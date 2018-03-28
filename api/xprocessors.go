package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/api/models"
	"github.com/vjeantet/bitfan/core"
)

type XProcessorApiController struct {
	path string
}

func (x *XProcessorApiController) Find(c *gin.Context) {
	behavior := c.Query("behavior")
	xprocessors := core.Storage().FindXProcessors(behavior)
	c.JSON(200, xprocessors)
}

func (p *XProcessorApiController) Create(c *gin.Context) {
	var xprocessor models.XProcessor
	err := c.BindJSON(&xprocessor)
	if err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	if xprocessor.Uuid == "" {
		uid, _ := uuid.NewV4()
		xprocessor.Uuid = uid.String()
	}

	core.Storage().CreateXProcessor(&xprocessor)

	c.Redirect(302, fmt.Sprintf("/%s/xprocessors/%s", p.path, xprocessor.Uuid))
}

func (p *XProcessorApiController) FindOneByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	mXProcessor, err := core.Storage().FindOneXProcessorByUUID(uuid, false)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}
	c.JSON(200, mXProcessor)
}

func (p *XProcessorApiController) UpdateByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	mXProcessor, err := core.Storage().FindOneXProcessorByUUID(uuid, false)
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

	if err := mapstructure.WeakDecode(data, &mXProcessor); err != nil {
		c.JSON(500, models.Error{Message: err.Error()})
		return
	}

	core.Storage().SaveXProcessor(&mXProcessor)

	c.Redirect(302, fmt.Sprintf("/%s/xprocessors/%s", p.path, uuid))
}

func (p *XProcessorApiController) DeleteByUUID(c *gin.Context) {
	uuid := c.Param("uuid")

	mXProcessor, err := core.Storage().FindOneXProcessorByUUID(uuid, false)
	if err != nil {
		c.JSON(404, models.Error{Message: err.Error()})
		return
	}

	core.Storage().DeleteXProcessor(&mXProcessor)

	c.JSON(204, "")
}
