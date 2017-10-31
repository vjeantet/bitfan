package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/k0kubun/pp"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/core/models"
	"github.com/vjeantet/bitfan/lib"
)

type AssetApiController struct {
	database     *gorm.DB
	dataLocation string
	path         string
}

func (a *AssetApiController) Create(c *gin.Context) {
	var asset models.Asset
	err := c.BindJSON(&asset)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	uid, _ := uuid.NewV4()
	asset.Uuid = uid.String()
	asset.Size = len(asset.Value)
	asset.ContentType = http.DetectContentType(asset.Value[:512])

	a.database.Create(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}

func (a *AssetApiController) FindOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}

	c.JSON(200, asset)
}

func (a *AssetApiController) DownloadOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+asset.Name+"\"")
	c.Data(200, "application/octet-stream", asset.Value)
}

func (a *AssetApiController) UpdateByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}

	err := c.BindJSON(&asset)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	a.database.Save(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}

func (a *AssetApiController) DeleteByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}

	a.database.Delete(&asset)

	c.JSON(204, "")
}

func (a *AssetApiController) ReplaceByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	asset := models.Asset{Uuid: uuid}
	if a.database.Where(&asset).First(&asset).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Asset " + uuid + " not found"})
		return
	}
	tmpasset := models.Asset{}
	err := c.BindJSON(&tmpasset)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	asset.Name = tmpasset.Name
	asset.Value = tmpasset.Value
	asset.Size = len(asset.Value)
	asset.ContentType = http.DetectContentType(asset.Value[:512])

	a.database.Save(&asset)

	c.Redirect(302, fmt.Sprintf("/%s/assets/%s", a.path, asset.Uuid))
}

type PipelineApiController struct {
	database     *gorm.DB
	dataLocation string
	path         string
}

func (p *PipelineApiController) Create(c *gin.Context) {
	var pipeline models.Pipeline
	err := c.BindJSON(&pipeline)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	uid, _ := uuid.NewV4()
	pipeline.Uuid = uid.String()
	for i, _ := range pipeline.Assets {
		uid, _ := uuid.NewV4()
		pipeline.Assets[i].Uuid = uid.String()
	}

	p.database.Create(&pipeline)

	// Handle optinal Start
	if pipeline.Active == true {
		err = p.startPipeline(pipeline.Uuid)
		if err != nil {
			c.JSON(500, models.Error{Message: err.Error()})
			return
		}
	}

	c.Redirect(302, fmt.Sprintf("/%s/pipelines/%s", p.path, pipeline.Uuid))
}

func (p *PipelineApiController) Find(c *gin.Context) {
	var pipelines []models.Pipeline
	var err error

	pipelines = []models.Pipeline{}
	p.database.Find(&pipelines)

	runningPipelines := Pipelines() //core
	for i, p := range pipelines {
		if pup, ok := runningPipelines[p.Uuid]; ok {
			pipelines[i].Active = true
			pipelines[i].LocationPath = pup.ConfigLocation
			pipelines[i].StartedAt = pup.StartedAt
		}
	}

	if err == nil {
		c.JSON(200, pipelines)
	} else {
		c.JSON(404, models.Error{Message: "no pipelines(s) running"})
	}
}

func (p *PipelineApiController) FindOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	mPipeline := models.Pipeline{Uuid: uuid}
	if p.database.Preload("Assets").Where(&mPipeline).First(&mPipeline).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Pipeline " + uuid + " not found"})
		return
	}

	runningPipeline, found := pipelines.Load(uuid)
	if found == true {
		mPipeline.StartedAt = runningPipeline.(*Pipeline).StartedAt
		mPipeline.Active = true
		mPipeline.LocationPath = runningPipeline.(*Pipeline).ConfigLocation
	}

	c.JSON(200, mPipeline)
}

func (p *PipelineApiController) UpdateByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	mPipeline := models.Pipeline{Uuid: uuid}
	if p.database.Where(&mPipeline).First(&mPipeline).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Pipeline " + uuid + " not found"})
		return
	}

	data := map[string]interface{}{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	if err := mapstructure.WeakDecode(data, &mPipeline); err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}
	p.database.Save(&mPipeline)

	// handle Start / Stop / Restart
	_, active := pipelines.Load(uuid)
	if nextActive, ok := data["active"]; ok {

		switch active {
		case true: // pipeline is on
			switch nextActive {
			case true: // restart
				Log().Debugf("restarting pipeline %s", uuid)
				err := p.stopPipeline(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
				err = p.startPipeline(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
			case false: // stop
				Log().Debugf("stopping pipeline %s", uuid)
				err := p.stopPipeline(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
			}
		case false: // pipeline is off
			switch nextActive {
			case true: // start pipeline
				Log().Debugf("starting pipeline %s", uuid)
				err := p.startPipeline(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}

			case false: // can not stop a non running pipeline
				c.JSON(428, models.Error{Message: "pipeline " + uuid + " is not running"})
				return
			}
		}
	}

	c.Redirect(302, fmt.Sprintf("/%s/pipelines/%s", p.path, uuid))
}

func (p *PipelineApiController) DeleteByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	mPipeline := models.Pipeline{Uuid: uuid}
	if p.database.Where(&mPipeline).First(&mPipeline).RecordNotFound() {
		c.JSON(404, models.Error{Message: "Pipeline " + uuid + " not found"})
		return
	}

	p.database.Delete(&mPipeline)
	p.database.Delete(models.Asset{}, "pipeline_uuid = ?", uuid)
	// TODO : Delete related Assets

	c.JSON(204, "")
}

func (p *PipelineApiController) stopPipeline(uuid string) error {
	return StopPipeline(uuid)
}
func (p *PipelineApiController) startPipeline(uuid string) error {

	pipeline := models.Pipeline{Uuid: uuid}
	if p.database.Preload("Assets").Where(&pipeline).First(&pipeline).RecordNotFound() {
		return fmt.Errorf("Pipeline %s not found", uuid)
	}

	var cwd string

	// save Assets
	// directory = $data / remote / UUID /

	uidString := fmt.Sprintf("%s_%d", pipeline.Uuid, time.Now().Unix())

	cwd = filepath.Join(p.dataLocation, "_pipelines", uidString)
	Log().Debugf("configuration %s stored to %s", uidString, cwd)
	os.MkdirAll(cwd, os.ModePerm)

	pp.Println("pipeline-->", pipeline)

	//Save assets to cwd
	for _, asset := range pipeline.Assets {
		dest := filepath.Join(cwd, asset.Name)
		dir := filepath.Dir(dest)
		os.MkdirAll(dir, os.ModePerm)
		if err := ioutil.WriteFile(dest, asset.Value, 07770); err != nil {
			return err
		}

		if asset.Type == models.ASSET_TYPE_ENTRYPOINT {
			pipeline.ConfigLocation = filepath.Join(cwd, asset.Name)
		}

		if pipeline.ConfigLocation == "" {
			return fmt.Errorf("missing entrypont for pipeline %s", pipeline.Uuid)
		}

		Log().Debugf("configuration %s asset %s stored", uidString, asset.Name)
	}

	Log().Debugf("configuration %s pipeline %s ready to be loaded", uidString, pipeline.ConfigLocation)

	var loc *lib.Location
	var err error
	loc, err = lib.NewLocation(pipeline.ConfigLocation, cwd)

	if err != nil {
		return err
	}

	ppl := loc.ConfigPipeline()
	ppl.Name = pipeline.Label
	ppl.Uuid = pipeline.Uuid

	agt, err := loc.ConfigAgents()
	if err != nil {
		return err
	}

	UUID, err := StartPipeline(&ppl, agt)
	if err != nil {
		return err
	}
	Log().Debugf("Pipeline %s started UUID=%s, uuid=%s", pipeline.Label, UUID, pipeline.Uuid)

	return nil
}
