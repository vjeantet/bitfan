package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vjeantet/bitfan/core"
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
