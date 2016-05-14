package runtime

import (
	"fmt"
	"math"

	"github.com/veino/config"
	"github.com/veino/runtime/memory"
	"github.com/veino/runtime/metrics"
	"github.com/veino/runtime/scheduler"
	"github.com/veino/veino"

	"time"
)

var availableProcessorsFactory = map[string]veino.ProcessorFactory{}
var agentsRegistry = map[string][]*agent{}

var Scheduler veino.Scheduler

var stater metrics.IStats
var countEventOnboarded int64
var countEventConsumed int64
var countEventDropped int64
var countGotine int64

func init() {
	Scheduler = scheduler.NewScheduler()
	Scheduler.Start()
	stater = &metrics.StatsVoid{}
}

// RegisterProcessor is called by the processor loader when the program starts
// each processor give its name and factory func()
func RegisterProcessor(name string, procFact veino.ProcessorFactory) {
	Logger().Debugf("%s processor registered", name)
	availableProcessorsFactory[name] = procFact
}

//Start initialize the veino system
func Start() {
	listenAndServeHTTP("127.0.0.1:9090")
	memoryMap["veino"] = memory.NewMemory("veino")
	Logger().Infoln("veino started")
}

// StartAgent create the agent's pool with instances of agents, start them, and attach them to the bus
func StartAgent(name string, agentConf config.Agent) (int, error) {

	//Check if agent is already running
	if _, ok := agentsRegistry[name]; ok {
		Logger().Debugf("agent '%s' exists in agentsRegistry", name)
		return 0, fmt.Errorf("an agent nammed '%s' is already running", name)
	}

	Logger().Debugf("Starting agent %s with conf %v", name, agentConf)

	// Create new agent's processors
	nb := 0
	for i := 0; i < agentConf.PoolSize; i++ {
		config.Normalize(&agentConf)

		agent, err := NewAgent(agentConf, i)
		if err != nil {
			Logger().Errorf("can not build agent for %s : %s", name, err.Error())
			return nb, fmt.Errorf("can not build agent for %s : %s", name, err.Error())
		}
		agentsRegistry[name] = append(agentsRegistry[name], &agent)

		nb++
	}

	config.AddAgent(&agentConf)

	Logger().Debugf("%d agent processors created for agent %s", nb, name)

	//All created agents can now consume events.
	agentPacketChan := make(chan veino.IPacket, agentConf.Buffer)
	for _, agent := range agentsRegistry[name] {
		agent.packetChan = agentPacketChan
		agent.listen()
	}

	//Start agent, with the first agent's processor from the agent pool
	agent := agentsRegistry[name][0]

	agent.start(nil)
	Logger().Infof("agent %s started (%d processors)", name, nb)

	return nb, nil
}

func StartAgents(configAgents []config.Agent) error {
	orderedAgentConfList := config.Sort(configAgents, config.SortOutputsFirst)

	for _, agentConf := range orderedAgentConfList {
		_, err := StartAgent(agentConf.Name, agentConf)
		if err != nil {
			Logger().Fatalf("%s agent '%-s' can not start", agentConf.Type, agentConf.Name)
			return err
		}
	}

	return nil
}

// Stop stops all agents
func Stop() {
	Logger().Infoln("stopping...")

	//Stop scheduler
	Scheduler.Stop()

	//Get a list of agent's name sorted to stop them in a logical order
	agentNamesList, cycle := config.Agents().NamesSort(config.SortInputsFirst)
	if cycle != nil {
		// This case agent A wait for B events, and B waits for A events.. (can be more complex : A->B->C->A = cycle !)
		Logger().Warnf("Cyclic Agent circuit detected: %s", cycle)
		Logger().Warningln("I don't know how to handle this case in a clean way...")
	}

	//Stop Agents
	Logger().Infof("stopping agents %s", agentNamesList)

	for _, agentName := range agentNamesList {
		Logger().Infof("Stopping agent %s", agentName)
		StopAgent(agentName)
	}

	Logger().Infof("%d events received", countEventOnboarded)
	Logger().Infof("%d events consumed, %d lost", countEventConsumed, countEventDropped)
	Logger().Infoln("Good bye !")
}

// StopAgent removes any existing scheduling for the agent,
//  then stop each agent instance in the agent's pool, and close the pipeline
//  it returns the number of agent stopped
func StopAgent(agentName string) (int, error) {
	if _, ok := agentsRegistry[agentName]; !ok {
		Logger().Warnf("agent %s not found in agentsRegistry", agentName)
		return 0, fmt.Errorf("agent %s not found", agentName)
	}

	Scheduler.Remove(agentName)

	// Send a stop signal to each agent
	agents := agentsRegistry[agentName]
	close(agents[0].packetChan)
	for _, agent := range agents {
		Logger().Debugf("Waiting for agent %s to stop", agent.Name)
		agent.stop(nil)
	}

	config.RemoveAgent(agentName)

	// Remove agent from the registry
	delete(agentsRegistry, agentName)

	Logger().Debugf("agent %s stopped (%d processors)", agentName, len(agents))
	return len(agents), nil
}

func ShowMetrics() {
	go func() {
		var lastconsume, tmp int64
		for {
			ticker := time.Tick(time.Millisecond * 1000)
			<-ticker
			if countEventOnboarded == 0 {
				fmt.Printf("\r[waiting for events]")
				continue
			}
			div := math.Floor(float64(countEventConsumed) / float64(countEventOnboarded) * 100)
			tmp = countEventConsumed
			packetsPerSecond := (tmp - lastconsume)
			lastconsume = tmp

			fmt.Printf("\r[%3d] %6dm/s (waiting:%-3d) Received:%-5d Consumed:%-5d Lost:%-5d                             ",
				int(div),
				packetsPerSecond,
				countGotine,
				countEventOnboarded,
				countEventConsumed,
				countEventDropped)
		}
	}()

}

func SetIStat(s metrics.IStats) {
	stater = s
}
func CountEventDropped() int64 {
	return countEventDropped
}
func CountEventConsumed() int64 {
	return countEventConsumed
}
