package core

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/sync/syncmap"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/spf13/viper"

	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vjeantet/bitfan/store"
)

var (
	metrics     Metrics
	myScheduler *scheduler
	myMemory    *memory
	myStore     *store.Store

	availableProcessorsFactory map[string]ProcessorFactory = map[string]ProcessorFactory{}
	dataLocation               string                      = "data"

	pipelines syncmap.Map = syncmap.Map{}
)

type fnMux func(sm *http.ServeMux)

type Options struct {
	Host         string
	HttpHandlers []fnMux
	Debug        bool
	VerboseLog   bool
	LogFile      string
	DataLocation string
}

func init() {
	metrics = &MetricsVoid{}
	myScheduler = newScheduler()
	myScheduler.Start()
	//Init Store
	myMemory = newMemory(dataLocation)
}

// RegisterProcessor is called by the processor loader when the program starts
// each processor give its name and factory func()
func RegisterProcessor(name string, procFact ProcessorFactory) {
	Log().Debugf("%s processor registered", name)
	availableProcessorsFactory[name] = procFact
}

func setMetrics(s Metrics) {
	metrics = s
}

func setDataLocation(location string) error {
	dataLocation = location
	fileInfo, err := os.Stat(dataLocation)
	if err != nil {
		err = os.MkdirAll(dataLocation, os.ModePerm)
		if err != nil {
			Log().Errorf("%s - %v", dataLocation, err)
			return err
		}
		Log().Debugf("created folder %s", dataLocation)
	} else {
		if !fileInfo.IsDir() {
			Log().Errorf("data path %s is not a directory", dataLocation)
			return err
		}
	}
	Log().Debugf("data location : %s", location)

	// DB
	myStore, err = store.New(location, Log())
	if err != nil {
		Log().Errorf("failed to start store : %s", err.Error())
		return err
	}

	return err
}

func webHookServer() fnMux {
	whPrefixURL = "/"
	commonHandlers := alice.New(loggingHandler, recoverHandler)
	return HTTPHandler("/", commonHandlers.ThenFunc(routerHandler))
}

// TODO : should be unexported
func PrometheusServer(path string) fnMux {
	setMetrics(NewPrometheus())
	return HTTPHandler(path, prometheus.Handler())
}

// TODO : should be unexported
func Storage() *store.Store {
	return myStore
}

func HTTPHandler(path string, s http.Handler) fnMux {
	return func(sm *http.ServeMux) {
		sm.Handle(path, s)
	}
}

func listenAndServe(addr string, hs ...fnMux) {
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

func Start(opt Options) {
	if opt.VerboseLog {
		setLogVerboseMode()
	}

	if opt.Debug {
		setLogDebugMode()
	}

	if opt.LogFile != "" {
		setLogOutputFile(viper.GetString("log"))
	}

	if err := setDataLocation(opt.DataLocation); err != nil {
		Log().Errorf("error with data location - %v", err)
		panic(err.Error())
	}

	if len(opt.HttpHandlers) > 0 {
		opt.HttpHandlers = append(opt.HttpHandlers, webHookServer())

		listenAndServe(opt.Host, opt.HttpHandlers...)
	}

	Log().Debugln("bitfan started")
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
		p, ok := GetPipeline(Uuid)
		if !ok {
			Log().Error("Stop Pipeline - pipeline " + Uuid + " not found")
			continue
		}
		err := p.Stop()
		if err != nil {
			Log().Error(err)
		}
	}

	myMemory.close()
	myStore.Close()
	return nil
}

func GetPipeline(UUID string) (*Pipeline, bool) {
	if i, found := pipelines.Load(UUID); found {
		return i.(*Pipeline), found
	} else {
		return nil, found
	}
}

// Pipelines returns running core.Pipeline
func Pipelines() map[string]*Pipeline {
	pps := map[string]*Pipeline{}
	pipelines.Range(func(key, value interface{}) bool {
		pps[key.(string)] = value.(*Pipeline)
		return true
	})
	return pps
}
