package core

import "github.com/vjeantet/bitfan/core/config"

var (
	metrics     Metrics
	myScheduler *scheduler
	myStore     *memory

	availableProcessorsFactory map[string]ProcessorFactory = map[string]ProcessorFactory{}
	dataLocation               string                      = "data"

	pipelines map[int]*Pipeline = map[int]*Pipeline{}
)

func init() {
	metrics = &MetricsVoid{}
	myScheduler = newScheduler()
	myScheduler.Start()
}

// RegisterProcessor is called by the processor loader when the program starts
// each processor give its name and factory func()
func RegisterProcessor(name string, procFact ProcessorFactory) {
	Log().Debugf("%s processor registered", name)
	availableProcessorsFactory[name] = procFact
}

func SetMetrics(s Metrics) {
	metrics = s
}

func SetDataLocation(location string) {
	dataLocation = location
	Log().Debugf("data location : %s", location)
}

// Start runtime
func Start(addr string) {
	Log().Debugln("bitfan started")
	//Init Store
	myStore = NewMemory(dataLocation)
}

// StartPipeline load all agents form a configPipeline and returns pipeline's ID
func StartPipeline(configPipeline *config.Pipeline, configAgents []config.Agent) (int, error) {
	p, err := newPipeline(configPipeline, configAgents)
	if err != nil {
		return 0, err
	}
	pipelines[p.ID] = p

	err = p.start()

	return p.ID, err
}

func StopPipeline(ID int) error {
	err := pipelines[ID].stop()
	if err != nil {
		return err
	}
	delete(pipelines, ID)
	return nil
}

// Stop each pipeline
func Stop() error {
	var IDS = []int{}
	for ID, _ := range pipelines {
		IDS = append(IDS, ID)
	}
	for _, ID := range IDS {
		err := StopPipeline(ID)
		if err != nil {
			Log().Error(err)
		}
	}

	myStore.close()

	return nil
}

func Pipelines() map[int]*Pipeline {
	return pipelines
}
