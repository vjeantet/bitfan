package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/vjeantet/bitfan/core/config"
	"github.com/vjeantet/bitfan/processors"
)

type agent struct {
	ID               int
	Label            string
	processor        processors.Processor
	packetChan       chan *event
	outputs          map[int][]chan *event
	Done             chan bool
	concurentProcess int
	conf             config.Agent
}

func NewAgent(conf config.Agent) (*agent, error) {
	return newAgent(conf)
}

// build an agent and return its input chan
func newAgent(conf config.Agent) (*agent, error) {
	// Check that the agent's processor type is supported
	if _, ok := availableProcessorsFactory[conf.Type]; !ok {
		return nil, fmt.Errorf("Processor %s not found", conf.Type)
	}

	// Create a new Processor processor
	proc := availableProcessorsFactory[conf.Type]()
	if proc == nil {
		return nil, fmt.Errorf("Can not start processor %s", conf.Type)
	}

	a := &agent{
		packetChan: make(chan *event, conf.Buffer),
		outputs:    map[int][]chan *event{},
		processor:  proc,
		Done:       make(chan bool),
		conf:       conf,
	}

	// Configure the agent (and its processor)
	if err := a.configure(&conf); err != nil {
		return nil, fmt.Errorf("Can not configure agent %s : %v", conf.Type, err)
	}

	return a, nil
}

func (a *agent) configure(conf *config.Agent) error {
	a.ID = conf.ID
	a.Label = conf.Label
	a.processor.SetPipelineUUID(a.conf.PipelineUUID)

	ctx := processorContext{}
	ctx.logger = NewLogger("pipeline",
		map[string]interface{}{
			"processor_type":  conf.Type,
			"pipeline_uuid":   conf.PipelineUUID,
			"processor_label": conf.Label,
		},
	)

	ctx.packetSender = a.send // 	data["processor_type"] = proc_type
	// 	data["pipeline_uuid"] = pipelineUUID
	// 	data["processor_label"] = proc_label
	ctx.packetBuilder = NewPacket
	ctx.dataLocation = filepath.Join(dataLocation, conf.Type)
	ctx.configWorkingLocation = conf.Wd
	ctx.memory = myMemory.Space(conf.Type)
	ctx.webHook = newWebHook(conf.PipelineName, conf.Label)
	ctx.store = Storage().NewProcessorStorage(conf.Type)
	Log().Debugf("data location : %s", ctx.dataLocation)
	if _, err := os.Stat(ctx.dataLocation); os.IsNotExist(err) {
		if err = os.MkdirAll(ctx.dataLocation, 0777); err != nil {
			Log().Errorf("data location creation error : %v", err)
		}
	}

	return a.processor.Configure(ctx, conf.Options)
}

func (a *agent) traceEvent(way string, packet processors.IPacket, portNumbers ...int) {
	verb := "received"
	if way == "OUT" {
		verb = "sent"
	}
	Log().e.WithFields(
		map[string]interface{}{
			"processor_type":  a.conf.Type,
			"pipeline_uuid":   a.conf.PipelineUUID,
			"processor_label": a.conf.Label,
			"event":           packet.Fields().Old(),
			"ports":           portNumbers,
			"trace":           way,
		},
	).Info(verb + " event by " + a.conf.Label + " on pipeline '" + a.conf.PipelineName + "'")
}

func (a *agent) send(packet processors.IPacket, portNumbers ...int) bool {
	if len(portNumbers) == 0 {
		portNumbers = []int{0}
	}

	if a.conf.Trace {
		a.traceEvent("OUT", packet, portNumbers...)
	}

	// for each portNumbes
	// send packet to each a.outputs[portNumber]
	for _, portNumber := range portNumbers {
		if len(a.outputs[portNumber]) == 1 {
			a.outputs[portNumber][0] <- packet.(*event)
			metrics.increment(METRIC_PROC_OUT, a.conf.PipelineName, a.Label)
		} else {
			// do not use go routine nor waitgroup as it slow down the processing
			for _, out := range a.outputs[portNumber] {
				// Clone() is a time killer
				// TODO : failback if out does not take out packet on x ms (share on a bitfanSlave)
				out <- packet.Clone().(*event)
				metrics.increment(METRIC_PROC_OUT, a.conf.PipelineName, a.Label)
			}
		}
	}
	return true
}

type processorContext struct {
	packetSender          processors.PacketSender
	packetBuilder         processors.PacketBuilder
	logger                processors.Logger
	memory                processors.Memory
	webHook               processors.WebHook
	store                 processors.IStore
	dataLocation          string
	configWorkingLocation string
}

func (p processorContext) Log() processors.Logger {
	return p.logger
}
func (p processorContext) Memory() processors.Memory {
	return p.memory
}

func (p processorContext) WebHook() processors.WebHook {
	return p.webHook
}
func (p processorContext) PacketSender() processors.PacketSender {
	return p.packetSender
}
func (p processorContext) PacketBuilder() processors.PacketBuilder {
	return p.packetBuilder
}
func (p processorContext) ConfigWorkingLocation() string {
	return p.configWorkingLocation
}

func (p processorContext) DataLocation() string {
	return p.dataLocation
}

func (p processorContext) Store() processors.IStore {
	return p.store
}

func (a *agent) addOutput(in chan *event, portNumber int) error {
	a.outputs[portNumber] = append(a.outputs[portNumber], in)
	return nil
}

// Start agent
func (a *agent) start() error {
	// Start processor
	a.processor.Start(NewPacket("start", map[string]interface{}{}))

	// Maximum number of concurent packet consumption ?
	var maxConcurentPackets = a.conf.PoolSize

	if a.processor.MaxConcurent() > 0 && maxConcurentPackets > a.processor.MaxConcurent() {
		maxConcurentPackets = a.processor.MaxConcurent()
		Log().Infof("agent %s : starting only %d worker(s) (processor's limit)", a.Label, a.processor.MaxConcurent())
	}

	// Start in chan loop and a.processor.Receive(e) !
	Log().Debugf("agent %s : %d workers", a.Label, maxConcurentPackets)
	go func(maxConcurentPackets int) {
		var wg = &sync.WaitGroup{}

		wg.Add(maxConcurentPackets)
		for i := 1; i <= maxConcurentPackets; i++ {
			go a.listen(wg)
		}
		wg.Wait()

		Log().Debugf("processor (%d) - stopping (no more packets)", a.ID)
		if err := a.processor.Stop(NewPacket("", nil)); err != nil {
			Log().Errorf("%s %d : %v", a.conf.Type, a.ID, err)
		}
		close(a.Done)
		Log().Debugf("processor (%d) - stopped", a.ID)
	}(maxConcurentPackets)

	// Register scheduler if needed
	if a.conf.Schedule != "" {
		Log().Debugf("agent %s : schedule=%s", a.Label, a.conf.Schedule)
		err := myScheduler.Add(a.Label, a.conf.Schedule, func() {
			go a.processor.Tick(NewPacket("", nil))
		})
		if err != nil {
			Log().Errorf("schedule start failed - %s : %v", a.Label, err)
		} else {
			Log().Debugf("agent %s(%s) scheduled with %s", a.Label, a.ID, a.conf.Schedule)
		}
	}

	return nil
}

// listen plugs the agent processor to its event chan
func (a *agent) listen(wg *sync.WaitGroup) {
	Log().Debugf("Starting EventLoop on %d-%s", a.ID, a.Label)
	for e := range a.packetChan {
		// Receive a work request.
		metrics.set(METRIC_CONNECTION_TRANSIT, a.conf.PipelineName, a.Label, len(a.packetChan))

		if a.conf.Trace {
			a.traceEvent("IN", e, 0)
		}

		if err := a.processor.Receive(e); err != nil {
			Log().Errorf("agent %s: %v", a.conf.Type, err)
		}
		metrics.increment(METRIC_PROC_IN, a.conf.PipelineName, a.Label)
	}
	wg.Done()
}

func (a *agent) stop() {
	myScheduler.Remove(a.Label)
	Log().Debugf("agent %d schedule job removed", a.ID)

	// unregister processor's webhooks URLs
	if wh := a.processor.B().WebHook; wh != nil {
		wh.Unregister()
	}

	Log().Debugf("agent %d webhook routes unregistered", a.ID)

	Log().Debugf("Processor '%s' stopping... - %d in pipe ", a.Label, len(a.packetChan))
	close(a.packetChan)
	<-a.Done
	Log().Debugf("Processor %s stopped", a.Label)
}

func (a *agent) pause() {

}

func (a *agent) resume() {

}
