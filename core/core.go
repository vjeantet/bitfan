package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/syncmap"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/spf13/viper"

	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vjeantet/bitfan/core/config"
	"github.com/vjeantet/bitfan/core/models"
	"github.com/vjeantet/bitfan/core/store"
	"github.com/vjeantet/bitfan/lib"
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

	if opt.AutoStart {
		pipelinesToStart := myStore.FindPipelinesWithAutoStart()
		for _, p := range pipelinesToStart {
			StartPipelineByUUID(p.Uuid)
		}
	}

	if len(opt.HttpHandlers) > 0 {
		opt.HttpHandlers = append(opt.HttpHandlers, webHookServer())

		listenAndServe(opt.Host, opt.HttpHandlers...)
	}

	Log().Debugln("bitfan started")
}

func StartPipelineByUUID(UUID string) error {
	tPipeline, err := myStore.FindOnePipelineByUUID(UUID, true)
	if err != nil {
		return err
	}

	uidString := fmt.Sprintf("%s_%d", tPipeline.Uuid, time.Now().Unix())

	cwd := filepath.Join(dataLocation, "_pipelines", uidString)
	Log().Debugf("configuration %s stored to %s", uidString, cwd)
	os.MkdirAll(cwd, os.ModePerm)

	//Save assets to cwd
	for _, asset := range tPipeline.Assets {
		dest := filepath.Join(cwd, asset.Name)
		dir := filepath.Dir(dest)
		os.MkdirAll(dir, os.ModePerm)
		if err := ioutil.WriteFile(dest, asset.Value, 07770); err != nil {
			return err
		}

		if asset.Type == models.ASSET_TYPE_ENTRYPOINT {
			tPipeline.ConfigLocation = filepath.Join(cwd, asset.Name)
		}

		if tPipeline.ConfigLocation == "" {
			return fmt.Errorf("missing entrypoint for pipeline %s", tPipeline.Uuid)
		}

		Log().Debugf("configuration %s asset %s stored", uidString, asset.Name)
	}

	Log().Debugf("configuration %s pipeline %s ready to be loaded", uidString, tPipeline.ConfigLocation)

	//TODO : resolve lib.Location dans location.Location

	var loc *lib.Location
	loc, err = lib.NewLocation(tPipeline.ConfigLocation, cwd)
	if err != nil {
		return err
	}

	ppl := loc.ConfigPipeline()
	ppl.Name = tPipeline.Label
	ppl.Uuid = tPipeline.Uuid

	agt, err := loc.ConfigAgents()
	if err != nil {
		return err
	}

	nUUID, err := StartPipeline(&ppl, agt)
	if err != nil {
		return err
	}

	Log().Debugf("Pipeline %s started UUID=%s", tPipeline.Label, nUUID)
	return nil
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
