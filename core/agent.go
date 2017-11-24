package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/vjeantet/bitfan/core/metrics"
	"github.com/vjeantet/bitfan/processors"
)

type ProcessorFactory func() processors.Processor

type Agent struct {
	ID               int
	Label            string
	processor        processors.Processor
	packetChan       chan *event
	outputs          map[int][]chan *event
	Done             chan bool
	concurentProcess int
	// conf             config.Agent

	Sources         []string `json:"sources"`
	AgentSources    PortList
	AgentRecipients PortList
	Type            string `json:"type"`
	Schedule        string `json:"schedule"`
	Trace           bool   `json:"trace"`
	PoolSize        int    `json:"pool_size"`
	PipelineName    string
	PipelineUUID    string
	Buffer          int `json:"buffer_size"`
	Options         map[string]interface{}
	Wd              string
}

var agentIndex int = 0

func NewAgent() Agent {
	agentIndex++
	return Agent{
		ID: agentIndex,
	}
}

// build an agent and return its input chan
func buildAgent(conf *Agent) error {
	// Check that the agent's processor type is supported
	if _, ok := availableProcessorsFactory[conf.Type]; !ok {
		return fmt.Errorf("Processor %s not found", conf.Type)
	}

	// Create a new Processor processor
	proc := availableProcessorsFactory[conf.Type]()
	if proc == nil {
		return fmt.Errorf("Can not start processor %s", conf.Type)
	}

	conf.packetChan = make(chan *event, conf.Buffer)
	conf.outputs = map[int][]chan *event{}
	conf.processor = proc
	conf.Done = make(chan bool)
	conf.Options = conf.Options

	// Configure the agent (and its processor)
	if err := conf.configure(); err != nil {
		return fmt.Errorf("Can not configure agent %s : %v", conf.Type, err)
	}

	return nil
}

func (a *Agent) configure() error {

	a.processor.SetPipelineUUID(a.PipelineUUID)

	ctx := processorContext{}
	ctx.logger = NewLogger("pipeline",
		map[string]interface{}{
			"processor_type":  a.Type,
			"pipeline_uuid":   a.PipelineUUID,
			"processor_label": a.Label,
		},
	)

	ctx.packetSender = a.send // 	data["processor_type"] = proc_type
	// 	data["pipeline_uuid"] = pipelineUUID
	// 	data["processor_label"] = proc_label
	ctx.packetBuilder = newPacket
	ctx.dataLocation = filepath.Join(dataLocation, a.Type)
	ctx.configWorkingLocation = a.Wd
	ctx.memory = myMemory.Space(a.Type)
	ctx.webHook = newWebHook(a.PipelineName, a.Label)

	var err error
	ctx.store, err = Storage().NewProcessorStorage(a.Type)
	if err != nil {
		Log().Errorf("Storage error : %s", err.Error())
	}

	Log().Debugf("data location : %s", ctx.dataLocation)
	if _, err := os.Stat(ctx.dataLocation); os.IsNotExist(err) {
		if err = os.MkdirAll(ctx.dataLocation, 0777); err != nil {
			Log().Errorf("data location creation error : %v", err)
		}
	}

	return a.processor.Configure(ctx, a.Options)
}

func (a *Agent) traceEvent(way string, packet processors.IPacket, portNumbers ...int) {
	verb := "received"
	if way == "OUT" {
		verb = "sent"
	}
	Log().e.WithFields(
		map[string]interface{}{
			"processor_type":  a.Type,
			"pipeline_uuid":   a.PipelineUUID,
			"processor_label": a.Label,
			"event":           packet.Fields().Old(),
			"ports":           portNumbers,
			"trace":           way,
		},
	).Info(verb + " event by " + a.Label + " on pipeline '" + a.PipelineName + "'")
}

func (a *Agent) send(packet processors.IPacket, portNumbers ...int) bool {
	if len(portNumbers) == 0 {
		portNumbers = []int{0}
	}

	if a.Trace {
		a.traceEvent("OUT", packet, portNumbers...)
	}

	// for each portNumbes
	// send packet to each a.outputs[portNumber]
	for _, portNumber := range portNumbers {
		if len(a.outputs[portNumber]) == 1 {
			a.outputs[portNumber][0] <- packet.(*event)
			myMetrics.Increment(metrics.PROC_OUT, a.PipelineName, a.Label)
		} else {
			// do not use go routine nor waitgroup as it slow down the processing
			for _, out := range a.outputs[portNumber] {
				// Clone() is a time killer
				// TODO : failback if out does not take out packet on x ms (share on a bitfanSlave)
				out <- packet.Clone().(*event)
				myMetrics.Increment(metrics.PROC_OUT, a.PipelineName, a.Label)
			}
		}
	}
	return true
}

func (a *Agent) addOutput(in chan *event, portNumber int) error {
	a.outputs[portNumber] = append(a.outputs[portNumber], in)
	return nil
}

// Start agent
func (a *Agent) start() error {
	// Start processor
	a.processor.Start(newPacket("start", map[string]interface{}{}))

	// Maximum number of concurent packet consumption ?
	var maxConcurentPackets = a.PoolSize

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
		if err := a.processor.Stop(newPacket("", nil)); err != nil {
			Log().Errorf("%s %d : %v", a.Type, a.ID, err)
		}
		close(a.Done)
		Log().Debugf("processor (%d) - stopped", a.ID)
	}(maxConcurentPackets)

	// Register scheduler if needed
	if a.Schedule != "" {
		Log().Debugf("agent %s : schedule=%s", a.Label, a.Schedule)
		err := myScheduler.Add(a.Label, a.Schedule, func() {
			go a.processor.Tick(newPacket("", nil))
		})
		if err != nil {
			Log().Errorf("schedule start failed - %s : %v", a.Label, err)
		} else {
			Log().Debugf("agent %s(%s) scheduled with %s", a.Label, a.ID, a.Schedule)
		}
	}

	return nil
}

// listen plugs the agent processor to its event chan
func (a *Agent) listen(wg *sync.WaitGroup) {
	Log().Debugf("Starting EventLoop on %d-%s", a.ID, a.Label)
	for e := range a.packetChan {
		// Receive a work request.
		myMetrics.Set(metrics.CONNECTION_TRANSIT, a.PipelineName, a.Label, len(a.packetChan))

		if a.Trace {
			a.traceEvent("IN", e, 0)
		}

		if err := a.processor.Receive(e); err != nil {
			Log().Errorf("agent %s: %v", a.Type, err)
		}
		myMetrics.Increment(metrics.PROC_IN, a.PipelineName, a.Label)
	}
	wg.Done()
}

func (a *Agent) stop() {
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

func (a *Agent) pause() {

}

func (a *Agent) resume() {

}
