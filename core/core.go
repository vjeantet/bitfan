package core

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/sync/syncmap"

	"github.com/vjeantet/bitfan/core/memory"
	"github.com/vjeantet/bitfan/core/metrics"
	"github.com/vjeantet/bitfan/core/webhook"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/store"
)

var (
	myMetrics   metrics.Metrics
	myScheduler *scheduler
	myMemory    *memory.Memory
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
	Prometheus   string
}

func init() {
	myMetrics = metrics.New()
	myScheduler = newScheduler()
	myScheduler.Start()
	//Init Store
	myMemory = memory.New()
}

// RegisterProcessor is called by the processor loader when the program starts
// each processor give its name and factory func()
func RegisterProcessor(name string, procFact ProcessorFactory) {
	Log().Debugf("%s processor registered", name)
	availableProcessorsFactory[name] = procFact
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
	Log().Infof("Ready to serve on %s", addr)
}

func Start(opt Options) {
	if opt.VerboseLog {
		setLogVerboseMode()
	}

	if opt.Debug {
		setLogDebugMode()
	}

	if opt.LogFile != "" {
		setLogOutputFile(opt.LogFile)
	}

	if err := setDataLocation(opt.DataLocation); err != nil {
		Log().Errorf("error with data location - %v", err)
		panic(err.Error())
	}

	if opt.Prometheus != "" {
		m := metrics.NewPrometheus(opt.Prometheus)
		opt.HttpHandlers = append(opt.HttpHandlers, HTTPHandler(m.Path, m.HTTPHandler()))
		myMetrics = m
	}

	// Load env
	envs := Storage().FindEnvs()
	for _, v := range envs {
		os.Setenv(v.Name, v.Value)
	}

	if len(opt.HttpHandlers) > 0 {
		webhook.Log = logger
		opt.HttpHandlers = append(opt.HttpHandlers, HTTPHandler("/h/", webhook.Handler(opt.Host)))
		opt.HttpHandlers = append(opt.HttpHandlers, HTTPHandler("/_/", webhook.Handler(opt.Host)))

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

	myMemory.Close()
	myStore.Close()
	return nil
}

func GetPipeline(UUID string) (*Pipeline, bool) {
	if i, found := pipelines.Load(UUID); found {
		i.(*Pipeline).Webhooks = webhook.WebHooks(UUID)
		if sjobs, ok := scheduleMap.Load(UUID); ok {
			i.(*Pipeline).Schedulers = sjobs.([]schedulerJob)
		}
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
		pps[key.(string)].Webhooks = webhook.WebHooks(key.(string))
		if sjobs, ok := scheduleMap.Load(key.(string)); ok {
			pps[key.(string)].Schedulers = sjobs.([]schedulerJob)
		}
		return true
	})
	return pps
}

// ProcessorsDocs returns available ProcessorDoc
func ProcessorsDocs(code string) map[string]*doc.Processor {
	docs := map[string]*doc.Processor{}
	if code != "" {
		if proc, ok := availableProcessorsFactory[code]; ok {
			docs[code] = proc().Doc()
		}

	} else {
		for code, proc := range availableProcessorsFactory {
			docs[code] = proc().Doc()
		}
	}

	return docs
}
