package runtime

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/veino/config"
	"github.com/veino/runtime/memory"
	"github.com/veino/veino"
)

const (
	schedulerAgentName = "veino.scheduler"
)

type agent struct {
	ID            int
	ProcessorType string
	Name          string
	Schedule      string
	packetChan    chan veino.IPacket
	Done          chan bool
	subProcessor  veino.Processor
	recipients    map[int][]string
}

var memoryMap = map[string]*memory.Memory{}
var webHookMap = map[string]*WebHook{}

func NewAgent(agentConf config.Agent, id int) (agent, error) {
	// Check that the agent's processor type is supported
	if _, ok := availableProcessorsFactory[agentConf.Type]; !ok {
		return agent{}, fmt.Errorf("Processor %s not found", agentConf.Type)
	}

	// Create a new Processor processor
	baseProcessor := availableProcessorsFactory[agentConf.Type](logger)
	if baseProcessor == nil {
		return agent{}, fmt.Errorf("Can not start processo %s", agentConf.Type)
	}

	// Create an agent with the processor
	a := agent{
		subProcessor: baseProcessor,
		ID:           id,
		Done:         make(chan bool),
	}

	// Configure the agent (and its processor)
	if err := a.configure(agentConf); err != nil {
		return agent{}, fmt.Errorf("Can not configure agent %s : %s", agentConf.Type, err)
	}

	// Set agent recipients
	a.recipients = map[int][]string{}
	for _, r := range agentConf.XRecipients {
		a.recipients[r.PortNumber] = append(a.recipients[r.PortNumber], r.AgentName)
	}

	return a, nil
}

func (p *agent) configure(conf config.Agent) error {

	p.ProcessorType = conf.Type
	p.Schedule = conf.Schedule
	p.Name = conf.Name

	// Magic attribute - when a processor has an attribute "Memory" typed as *veino.Memory THEN set it up !
	if s, err := reflect.TypeOf(p.subProcessor).Elem().FieldByName("Memory"); err == true {
		if s.Type.String() == "*memory.Memory" {
			if mem, ok := memoryMap[p.Name]; ok {
				mv := reflect.ValueOf(mem)
				reflect.ValueOf(p.subProcessor).Elem().FieldByName("Memory").Set(mv)
			} else {
				//set New Memory Objet
				mem := memory.NewMemory(p.Name)
				mv := reflect.ValueOf(mem)
				reflect.ValueOf(p.subProcessor).Elem().FieldByName("Memory").Set(mv)
				memoryMap[p.Name] = mem
			}
		}
	}

	// Magic attribute - when a processor has an attribute "WebHook" typed as *veino.WebHook THEN set it up !
	if s, err := reflect.TypeOf(p.subProcessor).Elem().FieldByName("WebHook"); err == true {
		if s.Type.String() == "*veino.WebHook" {
			if wh, ok := webHookMap[p.Name]; ok {
				whv := reflect.ValueOf(wh)
				reflect.ValueOf(p.subProcessor).Elem().FieldByName("WebHook").Set(whv)
			} else {
				//set New Memory Objet
				wh := NewWebHook(p.Name)
				whv := reflect.ValueOf(wh)
				reflect.ValueOf(p.subProcessor).Elem().FieldByName("WebHook").Set(whv)
				webHookMap[p.Name] = wh
			}
		}
	}

	// Magic attribute - when a processor has an attribute "NewPacket" typed as *veino.PacketBuilder THEN set it up !
	if s, err := reflect.TypeOf(p.subProcessor).Elem().FieldByName("NewPacket"); err == true {
		if s.Type.String() == "veino.PacketBuilder" {
			whv := reflect.ValueOf(NewPacket)
			reflect.ValueOf(p.subProcessor).Elem().FieldByName("NewPacket").Set(whv)
		}
	}

	// Magic attribute - when a processor has an attribute "Out" typed as *veino.PacketBuilder THEN set it up !
	if s, err := reflect.TypeOf(p.subProcessor).Elem().FieldByName("Send"); err == true {
		if s.Type.String() == "veino.PacketSender" {
			whv := reflect.ValueOf(p.send)
			reflect.ValueOf(p.subProcessor).Elem().FieldByName("Send").Set(whv)
		}
	}

	return p.subProcessor.Configure(conf.Options)

}

// listen plugs the agent processor to its event chan
func (p *agent) listen() {
	Logger().Debugf("Starting EventLoop on %s[%d]", p.Name, p.ID)

	go func() {
		for {
			select {
			case e := <-p.packetChan:
				// Receive a work request.
				if e == nil {
					Logger().Debugf("%s processor %d listen received null - quitting loop", p.Name, p.ID)
					close(p.Done)
					return
				}

				Logger().WithEvent(e).Debugf("Processor%d (%s) event received", p.ID, p.Name)
				switch e.Kind() {
				case tickEvent:
					if err := p.subProcessor.Tick(e); err != nil {
						Logger().WithEvent(e).Errorf("processor %s[%d]: %s", p.ProcessorType, p.ID, err.Error())
					}
					break
				default:
					p.receive(e)
					if err := p.subProcessor.Receive(e); err != nil {
						Logger().WithEvent(e).Errorf("agent %s[%d]: %s", p.ProcessorType, p.ID, err.Error())
					}
					break
				}
			}
		}
	}()
}

func (p *agent) start(e veino.IPacket) {
	// return
	Logger().Debugf("Start event received on processor %s for agent %s", p.Name, p.ProcessorType)

	//Init received event to 0
	memoryMap["veino"].Set("e.received."+p.Name, 0)

	if p.Schedule != "" {
		err := Scheduler.Add(p.Name, p.Schedule, func() {
			e := NewPacket("", nil)
			e.SetKind(tickEvent)
			routePacket(e, p.Name)
		})
		if err != nil {
			Logger().Errorf("schedule start failed - %s %d : %s", p.ProcessorType, p.ID, err.Error())
		}
	}

	if err := p.subProcessor.Start(e); err != nil {
		Logger().Errorf("%s %d : %s", p.ProcessorType, p.ID, err.Error())
	}
}

// Agent-wide receive func
func (a *agent) receive(e veino.IPacket) {
	memoryMap["veino"].IncrementInt("e.received."+a.Name, 1)
}

func (p *agent) send(packet veino.IPacket, portNumbers ...int) bool {
	if len(portNumbers) == 0 {
		portNumbers = []int{0}
	}
	for _, portNumber := range portNumbers {
		agentDest := p.recipients[portNumber]
		if len(agentDest) == 1 {
			routePacket(packet, agentDest[0])
		} else {
			for _, d := range agentDest {
				newPacket := packet.Clone()
				routePacket(newPacket, d)
			}
			packet = nil
		}

	}
	return true
}

// routePacket push event to the processor type dedicated channel
func routePacket(e veino.IPacket, dst string) {
	Logger().WithEvent(e).Debug("packet onboarded")
	atomic.AddInt64(&countGotine, 1)
	stater.Increment(dst+".wait", 1)
	defer stater.Decrement(dst+".wait", 1)
	defer atomic.AddInt64(&countGotine, -1)

	stater.Increment(dst+".sent", 1)
	atomic.AddInt64(&countEventOnboarded, 1)
	stater.Increment(dst+".onboarded", 1)

	//Find the destination packet channel

	agents, ok := agentsRegistry[dst]
	if !ok {
		Logger().WithEvent(e).Warnf("packet lost (no agent %s started)", dst)
		atomic.AddInt64(&countEventDropped, 1)
		stater.Increment(dst+".drop", 1)
		return
	}

	Logger().WithEvent(e).Debug("packet wait")

	agents[0].packetChan <- e

	atomic.AddInt64(&countEventConsumed, 1)
	stater.Increment(dst+".consumed", 1)
	Logger().WithEvent(e).Debug("packet consumed")
}

func (p *agent) stop(e veino.IPacket) {
	Logger().Debugf("Processor %s stop() [%d]", p.Name, p.ID)
	if err := p.subProcessor.Stop(e); err != nil {
		Logger().WithEvent(e).Errorf("%s %d : %s", p.ProcessorType, p.ID, err.Error())
	}
	<-p.Done
	Logger().Debugf("Processor %s stopped [%d]", p.Name, p.ID)
}
