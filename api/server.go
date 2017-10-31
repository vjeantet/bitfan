// Bitfan API.
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
//
// Host: 127.0.0.1:5123
// BasePath: /api/v1
// Version: 0.0.1
// License: Apache 2.0 http://www.apache.org/licenses/LICENSE-2.0.html
// Contact: Valere JEANTET<valere.jeantet@gmail.com> http://vjeantet.fr
// Schemes: http
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
//
// swagger:meta
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vjeantet/bitfan/core"
)

var plugins map[string]map[string]core.ProcessorFactory

func init() {
	// gin.SetMode(gin.ReleaseMode)
}

func Handler(path string, plugs map[string]map[string]core.ProcessorFactory) http.Handler {
	plugins = plugs

	r := gin.New()
	r.Use(gin.Recovery(), cors())
	v1 := r.Group(path)
	{

		// swagger:operation GET /pipelines pipeline listPipelines
		//
		// Lists pipelines.
		//
		// This will show all running pipelines.
		//
		// ---
		//
		// produces:
		// - application/json
		//
		//
		// responses:
		//   200:
		//     description: pipelines response
		//     schema:
		//       type: array
		//       items:
		//         "$ref": "#/definitions/Pipeline"
		//   default:
		//     description: unexpected error
		//     schema:
		//       "$ref": "#/definitions/Error"
		v1.GET("/pipelines", getPipelines)

		// swagger:operation DELETE /pipelines/{id} pipeline stopPipeline
		//
		// Stop a running pipeline.
		//
		// This will stop pipeline by ID.
		//
		// ---
		//
		//
		// produces:
		// - application/json
		//
		// parameters:
		//   - name: "id"
		//     in: "path"
		//     description: "Pipeline ID"
		//     required: true
		//     type: integer
		//
		//
		// responses:
		//   200:
		//     description: pipelines response
		//   default:
		//     description: unexpected error
		//     schema:
		//       "$ref": "#/definitions/Error"
		v1.DELETE("/pipelines/:id", deletePipeline)

		// swagger:operation POST /pipelines pipeline addPipeline
		//
		// Start a pipeline.
		//
		// This will start pipeline.
		//
		// ---
		// consumes:
		// - application/json
		//
		// produces:
		// - application/json
		//
		// parameters:
		// - in: "body"
		//   name: "body"
		//   description: "Pipeline object that needs to be started"
		//   required: true
		//   schema:
		//     $ref: "#/definitions/Pipeline"
		//
		// responses:
		//   200:
		//     description: pipeline response
		//     schema:
		//       "$ref": "#/definitions/Pipeline"
		//   default:
		//     description: unexpected error
		//     schema:
		//       "$ref": "#/definitions/Error"
		v1.POST("/pipelines", addPipeline)

		// swagger:operation GET /pipelines/{id} pipeline getPipeline
		//
		// Get a pipeline.
		//
		// This will show a running pipeline.
		//
		// ---
		//
		// produces:
		// - application/json
		//
		//
		// parameters:
		//   - name: "id"
		//     in: "path"
		//     description: "Pipeline ID"
		//     required: true
		//     type: integer
		//
		// responses:
		//   200:
		//     description: pipeline response
		//     schema:
		//       "$ref": "#/definitions/Pipeline"
		//   default:
		//     description: unexpected error
		//     schema:
		//       "$ref": "#/definitions/Error"
		v1.GET("/pipelines/:id", getPipeline)

		// swagger:operation GET /pipelines/{id}/assets pipeline getPipeline
		//
		// Get pipeline's assets
		//
		// This will show configuration assets from a running pipeline.
		//
		// ---
		//
		// produces:
		// - application/json
		//
		//
		// parameters:
		//   - name: "id"
		//     in: "path"
		//     description: "Pipeline ID"
		//     required: true
		//     type: integer
		//
		// responses:
		//   200:
		//     description: assets response
		//     schema:
		//       type: array
		//       items:
		//         "$ref": "#/definitions/Asset"
		//   default:
		//     description: unexpected error
		//     schema:
		//       "$ref": "#/definitions/Error"
		v1.GET("/pipelines/:id/assets", getPipelineAssets)

		// swagger:operation GET /docs doc listDocs
		//
		// Lists plugins.
		//
		// This will show all avaialable plugins.
		//
		// ---
		// produces:
		// - application/json
		//
		// responses:
		//   200:
		//     description: processor doc response
		//     schema:
		//       type: array
		//       items:
		//         "$ref": "#/definitions/processorDoc"
		//   default:
		//     description: unexpected error
		//     schema:
		//       "$ref": "#/definitions/Error"
		v1.GET("/docs", getDocs)

		v1.GET("/docs/inputs", getDocsInputs)
		v1.GET("/docs/inputs/:name", getDocsInputsByName)
		v1.GET("/docs/filters", getDocsFilters)
		v1.GET("/docs/filters/:name", getDocsFiltersByName)
		v1.GET("/docs/outputs", getDocsOutputs)
		v1.GET("/docs/outputs/:name", getDocsOutputsByName)
	}

	core.Log().Debugf("Serving API on /%s/ ", path)

	return r
}
