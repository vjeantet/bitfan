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
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/lib"
	"github.com/vjeantet/bitfan/processors/doc"
)

var plugins map[string]map[string]core.ProcessorFactory

func ServeREST(hostport string, plugs map[string]map[string]core.ProcessorFactory) {
	plugins = plugs
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors())
	v1 := r.Group("api/v1")
	{
		v1.GET("/pipelines", getPipelines)
		v1.DELETE("/pipelines/:id", deletePipeline)
		v1.POST("/pipelines", addPipeline)
		v1.GET("/pipelines/:id", getPipeline)

		v1.GET("/docs", getDocs)
		v1.GET("/docs/inputs", getDocsInputs)
		v1.GET("/docs/inputs/:name", getDocsInputsByName)
		v1.GET("/docs/filters", getDocsFilters)
		v1.GET("/docs/filters/:name", getDocsFiltersByName)
		v1.GET("/docs/outputs", getDocsOutputs)
		v1.GET("/docs/outputs/:name", getDocsOutputsByName)
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

	// b, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pp.Println("c-->", string(b))

	//Bind request data
	var pipeline Pipeline
	err := c.BindJSON(&pipeline)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var loc *lib.Location
	cwd, _ := os.Getwd()

	if pipeline.Content != "" {
		loc, err = lib.NewLocationContent(pipeline.Content, cwd)
	} else {
		loc, err = lib.NewLocation(pipeline.ConfigLocation, cwd)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ppl := loc.ConfigPipeline()
	if pipeline.Label != "" {
		ppl.Name = pipeline.Label
	}

	agt, err := loc.ConfigAgents()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ID, err := core.StartPipeline(&ppl, agt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/api/v1/pipelines/%d", ID))
	return

	// c.JSON(200, gin.H{"statut": "started"})

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
		c.JSON(500, gin.H{"error": fmt.Sprintf("malformed id : %s", c.Params.ByName("id"))})
		return
	}

	err = core.StopPipeline(id)
	if err == nil {
		c.JSON(200, gin.H{"id": c.Params.ByName("id"), "status": "deleted"})
	} else {
		c.JSON(404, gin.H{"error": err.Error()})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/pipelines/1
}

func getDocs(c *gin.Context) {
	// swagger:route GET /docs listDocs
	//
	// Lists plugins.
	//
	// This will show all avaialable plugins.
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

	data := map[string]map[string]*doc.Processor{}

	data["input"] = map[string]*doc.Processor{}
	data["input"] = docsByKind("input")

	data["filter"] = map[string]*doc.Processor{}
	data["filter"] = docsByKind("filter")

	data["output"] = map[string]*doc.Processor{}
	data["output"] = docsByKind("output")

	if err == nil {
		c.JSON(200, data)
	} else {
		c.JSON(404, gin.H{"error": "not found"})
	}
	// curl -i http://localhost:8080/api/v1/docs
}

func getDocsInputs(c *gin.Context) {
	c.JSON(200, docsByKind("input"))
	// curl -i http://localhost:8080/api/v1/docs/inputs
}

func getDocsFilters(c *gin.Context) {
	c.JSON(200, docsByKind("filter"))
	// curl -i http://localhost:8080/api/v1/docs/filters
}

func getDocsOutputs(c *gin.Context) {
	c.JSON(200, docsByKind("output"))
	// curl -i http://localhost:8080/api/v1/docs/outputs
}

func getDocsInputsByName(c *gin.Context) {
	data, err := docsByKindByName("input", c.Params.ByName("name"))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, data)
	}
}

func getDocsFiltersByName(c *gin.Context) {
	data, err := docsByKindByName("filter", c.Params.ByName("name"))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, data)
	}
}

func getDocsOutputsByName(c *gin.Context) {
	data, err := docsByKindByName("output", c.Params.ByName("name"))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, data)
	}
}

// swagger:model Doc
// Doce represents a processor documentation
//
// A Doc is ....
//
// A Doc can have.....
//
type ProcessorDoc struct {
	Name     string `json:"name"`
	Doc      string `json:"doc"`
	DocShort string `json:"doc_short"`
	Options  *struct {
		Doc     string `json:"doc"`
		Options []*struct {
			Name         string      `json:"name"`
			Alias        string      `json:"alias"`
			Doc          string      `json:"doc"`
			Required     bool        `json:"requiered"`
			Type         string      `json:"type"`
			DefaultValue interface{} `json:"default_value"`
			//LogstashExample
			ExampleLS string `json:"example"`
		} `json:"options"`
	} `json:"options"`
	Ports []*struct {
		Default bool   `json:"default"`
		Name    string `json:"name"`
		Number  int    `json:"number"`
		Doc     string `json:"doc"`
	} `json:"ports"`
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

	Content string `json:"config_content"`
}

type Error struct {
	Message string `json:"error"`
}
