package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/core/models"
)

type PipelineApiController struct {
	path string
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

	core.Storage().CreatePipeline(&pipeline)

	// Handle optinal Start
	if pipeline.Active == true {
		err = core.StartPipelineByUUID(pipeline.Uuid)
		if err != nil {
			c.JSON(500, models.Error{Message: err.Error()})
			return
		}
	}

	c.Redirect(302, fmt.Sprintf("/%s/pipelines/%s", p.path, pipeline.Uuid))
}

func (p *PipelineApiController) Find(c *gin.Context) {

	pipelines := core.Storage().FindPipelines()

	runningPipelines := core.Pipelines() //core
	for i, p := range pipelines {
		if pup, ok := runningPipelines[p.Uuid]; ok {
			pipelines[i].Active = true
			pipelines[i].LocationPath = pup.ConfigLocation
			pipelines[i].StartedAt = pup.StartedAt
		}
	}

	c.JSON(200, pipelines)

}

func (p *PipelineApiController) FindOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	mPipeline, err := core.Storage().FindOnePipelineByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	runningPipeline, found := core.GetPipeline(uuid)
	if found == true {
		mPipeline.StartedAt = runningPipeline.StartedAt
		mPipeline.Active = true
		mPipeline.LocationPath = runningPipeline.ConfigLocation
	}

	c.JSON(200, mPipeline)
}

func (p *PipelineApiController) UpdateByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	mPipeline, err := core.Storage().FindOnePipelineByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
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

	core.Storage().SavePipeline(&mPipeline)

	// handle Start / Stop / Restart
	_, active := core.GetPipeline(uuid)
	if nextActive, ok := data["active"]; ok {

		switch active {
		case true: // pipeline is on
			switch nextActive {
			case true: // restart
				core.Log().Debugf("restarting pipeline %s", uuid)
				err := core.StopPipeline(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
				err = core.StartPipelineByUUID(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
			case false: // stop
				core.Log().Debugf("stopping pipeline %s", uuid)
				err := core.StopPipeline(uuid)
				if err != nil {
					c.JSON(500, models.Error{Message: err.Error()})
					return
				}
			}
		case false: // pipeline is off
			switch nextActive {
			case true: // start pipeline
				core.Log().Debugf("starting pipeline %s", uuid)
				err := core.StartPipelineByUUID(uuid)
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

	mPipeline, err := core.Storage().FindOnePipelineByUUID(uuid)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	core.Storage().DeletePipeline(&mPipeline)

	c.JSON(204, "")
}
