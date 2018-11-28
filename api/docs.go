package api

import (
	"bitfan/core"
	"github.com/gin-gonic/gin"
)

type DocsController struct {
}

func (d *DocsController) FindAllProcessors(c *gin.Context) {
	docs := core.ProcessorsDocs("")
	c.JSON(200, docs)
}

func (d *DocsController) FindOneProcessorByCode(c *gin.Context) {
	docs := core.ProcessorsDocs(c.Param("code"))
	c.JSON(200, docs)
}
