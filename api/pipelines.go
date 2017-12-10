package api

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/api/models"
	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/entrypoint"
	"github.com/vjeantet/jodaTime"
)

type PipelineApiController struct {
	path string
}

func (b *PipelineApiController) DownloadAll(c *gin.Context) {

	// Create a new zip archive.
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Find all Pipelines
	pipelines := core.Storage().FindPipelines(true)

	for _, p := range pipelines {
		folderName := slugify(p.Label) + "_" + p.Uuid
		for _, a := range p.Assets {
			zipFile, err := zipWriter.Create(folderName + "/" + a.Name)
			if err != nil {
				c.String(500, err.Error())
				return
			}
			_, err = zipFile.Write(a.Value)
			if err != nil {
				c.String(500, err.Error())
				return
			}
		}
	}

	// Make sure to check the error on Close.
	err := zipWriter.Close()
	if err != nil {
		c.String(500, err.Error())
		return
	}

	filename := jodaTime.Format("'bitfan_pipelines_'YYYYMMdd-HHmmss'.zip'", time.Now())
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", strconv.Itoa(buf.Len()))

	c.Data(200, "application/zip", buf.Bytes())
}

func (p *PipelineApiController) Create(c *gin.Context) {
	var pipeline models.Pipeline
	err := c.BindJSON(&pipeline)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	if pipeline.Uuid == "" {
		uid, _ := uuid.NewV4()
		pipeline.Uuid = uid.String()
	}

	for i, _ := range pipeline.Assets {
		uid, _ := uuid.NewV4()
		pipeline.Assets[i].Uuid = uid.String()
	}

	if pipeline.Playground == false {
		core.Storage().CreatePipeline(&pipeline)
	}

	// Handle optinal Start
	if pipeline.Active == true {
		err = p.startPipeline(&pipeline)
		if err != nil {
			c.JSON(500, models.Error{Message: err.Error()})
			return
		}
	}

	c.Redirect(302, fmt.Sprintf("/%s/pipelines/%s", p.path, pipeline.Uuid))
}

func (p *PipelineApiController) startPipelineByUUID(UUID string) error {
	tPipeline, err := core.Storage().FindOnePipelineByUUID(UUID, true)
	if err != nil {
		return err
	}

	return p.startPipeline(&tPipeline)
}

func (p *PipelineApiController) startPipeline(tPipeline *models.Pipeline) error {
	entryPointPath, err := core.Storage().PreparePipelineExecutionStage(tPipeline)
	if err != nil {
		return err
	}

	var loc *entrypoint.Entrypoint
	loc, err = entrypoint.New(entryPointPath, "", entrypoint.CONTENT_REF)
	if err != nil {
		return err
	}

	ppl, err := loc.Pipeline()
	if err != nil {
		return err
	}

	ppl.Label = tPipeline.Label
	ppl.Uuid = tPipeline.Uuid

	nUUID, err := ppl.Start()
	if err != nil {
		return err
	}

	apiLogger.Debugf("Pipeline %s started UUID=%s", tPipeline.Label, nUUID)
	return nil
}

func (p *PipelineApiController) Find(c *gin.Context) {

	pipelines := core.Storage().FindPipelines(false)

	runningPipelines := core.Pipelines() //core
	for i, p := range pipelines {
		if pup, ok := runningPipelines[p.Uuid]; ok {
			pipelines[i].Active = true
			pipelines[i].LocationPath = pup.ConfigLocation
			pipelines[i].StartedAt = pup.StartedAt
			pipelines[i].Webhooks = pup.Webhooks
		}
	}

	c.JSON(200, pipelines)

}

func (p *PipelineApiController) FindOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	mPipeline, err := core.Storage().FindOnePipelineByUUID(uuid, false)
	if err != nil {

		if _, active := core.GetPipeline(uuid); !active {
			c.JSON(404, models.Error{Message: err.Error()})
			return
		}

		mPipeline.Playground = true

	}

	runningPipeline, found := core.GetPipeline(uuid)
	if found == true {
		mPipeline.StartedAt = runningPipeline.StartedAt
		mPipeline.Active = true
		mPipeline.LocationPath = runningPipeline.ConfigLocation
		mPipeline.Webhooks = runningPipeline.Webhooks

	}

	c.JSON(200, mPipeline)
}

func (p *PipelineApiController) UpdateByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	mPipeline, err := core.Storage().FindOnePipelineByUUID(uuid, false)
	if err != nil {

		// Is pipline is not running -> error
		if _, active := core.GetPipeline(uuid); !active {
			c.JSON(404, models.Error{Message: err.Error()})
			return
		}
		mPipeline.Playground = true
	}

	data := map[string]interface{}{}
	err = c.BindJSON(&data)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	if err := mapstructure.WeakDecode(data, &mPipeline); err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	if !mPipeline.Playground { // Ignore playground pipelines
		core.Storage().SavePipeline(&mPipeline)
	}

	// handle Start / Stop / Restart
	_, active := core.GetPipeline(uuid)
	if nextActive, ok := data["active"]; ok {

		switch active {
		case true: // pipeline is on
			switch nextActive {
			case true: // restart
				apiLogger.Debugf("restarting pipeline %s", uuid)
				err := core.StopPipeline(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
				err = p.startPipelineByUUID(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
			case false: // stop
				apiLogger.Debugf("stopping pipeline %s", uuid)
				err := core.StopPipeline(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
			}
		case false: // pipeline is off
			switch nextActive {
			case true: // start pipeline
				apiLogger.Debugf("starting pipeline %s", uuid)
				err := p.startPipelineByUUID(uuid)
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

	mPipeline, err := core.Storage().FindOnePipelineByUUID(uuid, false)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	core.Storage().DeletePipeline(&mPipeline)

	c.JSON(204, "")
}
