package api

// Package Pipelines API.
//
// the purpose of this api....
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v1
//     Version: 0.0.1
//     License: Apache 2.0 http://www.apache.org/licenses/LICENSE-2.0.html
//     Contact: Valere JEANTET<valere.jeantet@gmail.com> http://vjeantet.fr
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//
// swagger:meta

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vjeantet/bitfan/core"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func ServeREST(hostport string) {

	r := gin.Default()
	r.Use(cors())
	v1 := r.Group("api/v1")
	{
		v1.GET("/pipelines", getPipelines)
		v1.DELETE("/pipelines/:id", deletePipeline)
		v1.POST("/pipelines", addPipeline)
		v1.GET("/pipelines/:id", getPipeline)

		// v1.PUT("/pipelines/:id", UpdateUser)
		// v1.OPTIONS("/pipelines", OptionsUser)     // POST
		// v1.OPTIONS("/pipelines/:id", OptionsUser) // PUT, DELETE
	}
	if hostport == "" {
		hostport = "127.0.0.1:8080"
	}
	go r.Run(hostport)
}

func getPipelines(c *gin.Context) {
	// swagger:route GET /pipelines listPipelines
	//
	// Lists pipelines.
	//
	// This will show all running pipelines.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//
	//     Responses:
	//       default: genericError
	//       200: []Pipeline
	var pipelines []Pipeline
	var err error

	pipelines = []Pipeline{}
	ppls := core.Pipelines()
	for _, p := range ppls {
		pipelines = append(pipelines, Pipeline{
			ID:                 p.ID,
			Label:              p.Label,
			ConfigLocation:     p.ConfigLocation,
			ConfigHostLocation: p.ConfigHostLocation,
		})
	}

	if err == nil {
		c.JSON(200, pipelines)
	} else {
		c.JSON(404, gin.H{"error": "no pipelines(s) running"})
	}
	// curl -i http://localhost:8080/api/v1/pipelines
}

func getPipeline(c *gin.Context) {
	var pipeline Pipeline
	var err error
	var id int

	id, err = strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pipeline = Pipeline{}
	ppls := core.Pipelines()
	for _, p := range ppls {
		if p.ID == id {
			pipeline = Pipeline{
				ID:                 p.ID,
				Label:              p.Label,
				ConfigLocation:     p.ConfigLocation,
				ConfigHostLocation: p.ConfigHostLocation,
			}
			c.JSON(200, pipeline)
			return
		}
	}

	c.JSON(404, gin.H{"error": "no pipelines(s) running"})

}

func addPipeline(c *gin.Context) {
	// swagger:route POST /pipelines addPipeline
	//
	// Start a pipeline.
	//
	// This will start pipeline.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//
	//     Responses:
	//       default: genericError
	//       200: Pipeline

	// ID, err := core.StartPipeline(&starter.Pipeline, starter.Agents)
	ID, err := 1, error(nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("/api/v1/pipelines/%d", ID))

		return
	}

	c.JSON(200, gin.H{"statut": "started"})

	// curl -i -X DELETE http://localhost:8080/api/v1/pipelines/1
}

func deletePipeline(c *gin.Context) {
	// swagger:route DELETE /pipelines/:id stopPipeline
	//
	// Stop a running pipeline.
	//
	// This will stop pipeline by ID.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//
	//     Responses:
	//       default: genericError
	//       200: []Pipeline
	var err error
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = core.StopPipeline(id)
	if err == nil {
		c.JSON(200, gin.H{c.Params.ByName("id"): "deleted"})
	} else {
		c.JSON(404, gin.H{"error": err.Error()})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/pipelines/1
}

// swagger:model Pipeline
// Pipeline represents a pipeline
//
// A Pipeline is ....
//
// A Pipeline can have.....
//
type Pipeline struct {
	// the id for this pipeline
	ID int `json:"id"`
	// the Label
	// min length: 3
	Label string `json:"label"`
	// the location
	ConfigLocation string `json:"config_location"`
	// the location's host
	ConfigHostLocation string `json:"config_host_location"`
}
