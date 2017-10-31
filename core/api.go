package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// var plugins map[string]map[string]core.ProcessorFactory
func init() {
	gin.SetMode(gin.ReleaseMode)
}
func apiHandler(path string, db *gorm.DB, dataLocation string) http.Handler {

	r := gin.New()
	r.Use(
		gin.Recovery(),
		func(c *gin.Context) {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.Next()
		},
	)
	v2 := r.Group(path)
	{

		pipelineCtrl := &PipelineApiController{
			database:     db,
			dataLocation: dataLocation,
			path:         path,
		}
		assetCtrl := &AssetApiController{
			database:     db,
			dataLocation: dataLocation,
			path:         path,
		}

		// curl -i -X POST http://localhost:5123/api/v2/pipelines
		v2.POST("/pipelines", pipelineCtrl.Create) // créer pipeline

		// curl -i -X GET http://localhost:5123/api/v2/pipelines
		v2.GET("/pipelines", pipelineCtrl.Find) // list pipelines

		// curl -i -X GET http://localhost:5123/api/v2/pipelines/408b9a7b-933e-4d3d-6df1-65324a0a5315
		v2.GET("/pipelines/:uuid", pipelineCtrl.FindOneByUUID) // show pipeline

		// curl -i -X PATCH http://localhost:5123/api/v2/pipelines/408b9a7b-933e-4d3d-6df1-65324a0a5315
		v2.PATCH("/pipelines/:uuid", pipelineCtrl.UpdateByUUID) // update pipeline / stop / start / restart

		// curl -i -X DELETE http://localhost:5123/api/v2/pipelines/408b9a7b-933e-4d3d-6df1-65324a0a5315
		v2.DELETE("/pipelines/:uuid", pipelineCtrl.DeleteByUUID) // delete pipeline

		v2.POST("/assets", assetCtrl.Create)                         // create asset
		v2.GET("/assets/:uuid", assetCtrl.FindOneByUUID)             // show asset
		v2.GET("/assets/:uuid/content", assetCtrl.DownloadOneByUUID) // dl asset
		v2.PUT("/assets/:uuid", assetCtrl.ReplaceByUUID)             // replace asset
		v2.PATCH("/assets/:uuid", assetCtrl.UpdateByUUID)            // update asset
		v2.DELETE("/assets/:uuid", assetCtrl.DeleteByUUID)           // delete asset

		// v1.GET("/docs", getDocs)
		// v1.GET("/docs/inputs", getDocsInputs)
		// v1.GET("/docs/inputs/:name", getDocsInputsByName)
		// v1.GET("/docs/filters", getDocsFilters)
		// v1.GET("/docs/filters/:name", getDocsFiltersByName)
		// v1.GET("/docs/outputs", getDocsOutputs)
		// v1.GET("/docs/outputs/:name", getDocsOutputsByName)
	}

	Log().Debugf("Serving API on /%s/ ", path)

	return r
}
