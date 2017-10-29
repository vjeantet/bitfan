package core

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/sync/syncmap"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vjeantet/bitfan/core/config"
)

var (
	metrics     Metrics
	myScheduler *scheduler
	myStore     *memory

	availableProcessorsFactory map[string]ProcessorFactory = map[string]ProcessorFactory{}
	dataLocation               string                      = "data"

	pipelines syncmap.Map = syncmap.Map{}
)

type FnMux func(sm *http.ServeMux)

func init() {
	metrics = &MetricsVoid{}
	myScheduler = newScheduler()
	myScheduler.Start()
	//Init Store
	myStore = NewMemory(dataLocation)
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

func SetDataLocation(location string) error {
	dataLocation = location
	fileInfo, err := os.Stat(dataLocation)
	if err != nil {
		err = os.MkdirAll(dataLocation, os.ModePerm)
		if err != nil {
			Log().Errorf("%s - %v", dataLocation, err)
		} else {
			Log().Debugf("created folder %s", dataLocation)
		}
	} else {
		if !fileInfo.IsDir() {
			Log().Errorf("data path %s is not a directory", dataLocation)
		}
	}
	Log().Debugf("data location : %s", location)
	return err
}

// DataLocation returns the bitfan's data filepath
func DataLocation() string {
	return dataLocation
}

func WebHookServer() FnMux {
	whPrefixURL = "/"
	commonHandlers := alice.New(loggingHandler, recoverHandler)
	return HTTPHandler("/", commonHandlers.ThenFunc(routerHandler))
}

func PrometheusServer(path string) FnMux {
	SetMetrics(NewPrometheus())
	return HTTPHandler(path, prometheus.Handler())
}

func HTTPHandler(path string, s http.Handler) FnMux {
	return func(sm *http.ServeMux) {
		sm.Handle(path, s)
	}
}

func ListenAndServe(addr string, hs ...FnMux) {
	httpServerMux := http.NewServeMux()
	for _, h := range hs {
		h(httpServerMux)
	}
	go http.ListenAndServe(addr, httpServerMux)

	addrSpit := strings.Split(addr, ":")
	if addrSpit[0] == "0.0.0.0" {
		addrSpit[0] = fqdn.Get()
	}

	baseURL = fmt.Sprintf("http://%s:%s", addrSpit[0], addrSpit[1])

	Log().Infof("Ready to serve on %s", baseURL)
}

// StartPipeline load all agents form a configPipeline and returns pipeline's ID
func StartPipeline(configPipeline *config.Pipeline, configAgents []config.Agent) (string, error) {
	p, err := newPipeline(configPipeline, configAgents)
	if err != nil {
		return "", err
	}
	if _, ok := pipelines.Load(p.Uuid); ok {
		// a pipeline with same uuid is already running
		return "", fmt.Errorf("a pipeline with uuid %s is already running", p.Uuid)
	}

	pipelines.Store(p.Uuid, p)

	err = p.start()

	return p.Uuid, err
}

func StopPipeline(Uuid string) error {
	var err error
	if p, ok := pipelines.Load(Uuid); ok {
		err = p.(*Pipeline).stop()
	} else {
		err = fmt.Errorf("Pipeline %s not found", Uuid)
	}

	if err != nil {
		return err
	}

	pipelines.Delete(Uuid)
	return nil
}

// Stop each pipeline
func Stop() error {
	var Uuids = []string{}
	pipelines.Range(func(key, value interface{}) bool {
		Uuids = append(Uuids, key.(string))
		return true
	})

	for _, Uuid := range Uuids {
		err := StopPipeline(Uuid)
		if err != nil {
			Log().Error(err)
		}
	}

	myStore.close()

	return nil
}

func Pipelines() map[string]*Pipeline {
	pps := map[string]*Pipeline{}
	pipelines.Range(func(key, value interface{}) bool {
		pps[key.(string)] = value.(*Pipeline)
		return true
	})
	return pps
}
