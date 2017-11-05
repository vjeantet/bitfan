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
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vjeantet/bitfan/core/config"
	"github.com/vjeantet/bitfan/core/models"
	"github.com/vjeantet/bitfan/lib"
)

var (
	metrics     Metrics
	myScheduler *scheduler
	myMemory    *memory

	availableProcessorsFactory map[string]ProcessorFactory = map[string]ProcessorFactory{}
	dataLocation               string                      = "data"

	pipelines syncmap.Map = syncmap.Map{}

	db *gorm.DB
)

type FnMux func(sm *http.ServeMux)

func init() {
	metrics = &MetricsVoid{}
	myScheduler = newScheduler()
	myScheduler.Start()
	//Init Store
	myMemory = NewMemory(dataLocation)
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

	db, err = gorm.Open("sqlite3", filepath.Join(location, "bitfan.db"))
	// db.LogMode(true)

	if err != nil {
		Log().Errorf("failed to connect database : %s", err.Error())
		return err
	}

	// Migrate the schema
	db.AutoMigrate(&models.Pipeline{}, &models.Asset{})

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

func Database() *gorm.DB {
	return db
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

func RunAutoStartPipelines() {
	pipelinesToStart := []models.Pipeline{}
	db.Where(&models.Pipeline{AutoStart: true}).Find(&pipelinesToStart).RecordNotFound()

	for _, p := range pipelinesToStart {
		StartPipelineByUUID(p.Uuid)
	}
}

func StartPipelineByUUID(UUID string) error {
	tPipeline := models.Pipeline{Uuid: UUID}
	if db.Preload("Assets").Where(&tPipeline).First(&tPipeline).RecordNotFound() {
		return fmt.Errorf("Pipeline %s not found", UUID)
	}

	uidString := fmt.Sprintf("%s_%d", tPipeline.Uuid, time.Now().Unix())

	cwd := filepath.Join(DataLocation(), "_pipelines", uidString)
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
			return fmt.Errorf("missing entrypont for pipeline %s", tPipeline.Uuid)
		}

		Log().Debugf("configuration %s asset %s stored", uidString, asset.Name)
	}

	Log().Debugf("configuration %s pipeline %s ready to be loaded", uidString, tPipeline.ConfigLocation)

	//TODO : resolve lib.Location dans location.Location

	var loc *lib.Location
	var err error
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
	db.Close()
	return nil
}

func GetPipeline(UUID string) (*Pipeline, bool) {
	if i, found := pipelines.Load(UUID); found {
		return i.(*Pipeline), found
	} else {
		return nil, found
	}
}

func Pipelines() map[string]*Pipeline {
	pps := map[string]*Pipeline{}
	pipelines.Range(func(key, value interface{}) bool {
		pps[key.(string)] = value.(*Pipeline)
		return true
	})
	return pps
}
