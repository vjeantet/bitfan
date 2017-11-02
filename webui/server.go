package webui

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/vjeantet/bitfan/api"
	"github.com/vjeantet/bitfan/core/models"
	eztemplate "github.com/vjeantet/ez-gin-template"
)

var apiClient *api.RestClient

func Handler(assetsPath, URLpath string, dbpath string, apiBaseUrl string) http.Handler {
	apiClient = api.New(apiBaseUrl)

	r := gin.New()
	render := eztemplate.New()
	render.TemplatesDir = assetsPath + "/views/" // default
	render.Ext = ".html"                         // default
	render.Debug = true                          // default
	render.TemplateFuncMap = template.FuncMap{
		"dateFormat": (*templateFunctions)(nil).dateFormat,
		"ago":        (*templateFunctions)(nil).dateAgo,
		"string":     (*templateFunctions)(nil).toString,
		"b64":        (*templateFunctions)(nil).toBase64,
		"int":        (*templateFunctions)(nil).toInt,
		"time":       (*templateFunctions)(nil).asTime,
		"now":        (*templateFunctions)(nil).now,
		"isset":      (*templateFunctions)(nil).isSet,

		"numFmt": (*templateFunctions)(nil).numFmt,

		"safeHTML":     (*templateFunctions)(nil).safeHtml,
		"hTMLUnescape": (*templateFunctions)(nil).htmlUnescape,
		"hTMLEscape":   (*templateFunctions)(nil).htmlEscape,
		"lower":        (*templateFunctions)(nil).lower,
		"upper":        (*templateFunctions)(nil).upper,
		"trim":         (*templateFunctions)(nil).trim,
		"trimPrefix":   (*templateFunctions)(nil).trimPrefix,
		"hasPrefix":    (*templateFunctions)(nil).hasPrefix,
		"replace":      (*templateFunctions)(nil).replace,
		"markdown":     (*templateFunctions)(nil).toMarkdown,
	}

	r.HTMLRender = render.Init()
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store), gin.Recovery())
	// r.Use(gin.Recovery())
	g := r.Group(URLpath)
	{
		g.StaticFS("/public", http.Dir(assetsPath+"/public"))
		g.GET("/", getPipelines)

		// list pipelines
		g.GET("/logs", getLogs)

		// list pipelines
		g.GET("/pipelines", getPipelines)
		// New pipeline
		g.GET("/pipelines/:id/new", newPipeline)
		// Create pipeline
		g.POST("/pipelines", createPipeline)
		// Show pipeline
		g.GET("/pipelines/:id", editPipeline)
		// Save pipeline
		g.POST("/pipelines/:id", updatePipeline)

		// Start pipeline
		g.GET("/pipelines/:id/start", startPipeline)
		// Restart pipeline
		g.GET("/pipelines/:id/restart", startPipeline)
		// Stop pipeline
		g.GET("/pipelines/:id/stop", stopPipeline)

		// Delete asset
		g.GET("/pipelines/:id/delete", deletePipeline)
		// Show asset
		g.GET("/pipelines/:id/assets/:assetID", showAsset)
		// Create asset
		g.POST("/pipelines/:id/assets", createAsset)
		// Update asset
		g.POST("/pipelines/:id/assets/:assetID", updateAsset)
		// Replace asset
		g.PUT("/pipelines/:id/assets/:assetID", replaceAsset)
		// Download asset
		g.GET("/pipelines/:id/assets/:assetID/download", downloadAsset)
		// Delete asset
		g.GET("/pipelines/:id/assets/:assetID/delete", deleteAsset)
	}

	return r
}

func getLogs(c *gin.Context) {
	c.HTML(200, "logs/logs", gin.H{})
}

func getPipelines(c *gin.Context) {
	pipelines, _ := apiClient.Pipelines()

	c.HTML(200, "pipelines/index", gin.H{
		"pipelines": pipelines,
	})
}

func editPipeline(c *gin.Context) {
	id := c.Param("id")

	p, _ := apiClient.Pipeline(id)

	flashes := []string{}
	for _, m := range sessions.Default(c).Flashes() {
		flashes = append(flashes, m.(string))
	}
	sessions.Default(c).Save()

	c.HTML(200, "pipelines/edit", gin.H{
		"pipeline": p,
		"flashes":  flashes,
	})

}

func newPipeline(c *gin.Context) {
	c.HTML(200, "pipelines/new", gin.H{})
}

func createPipeline(c *gin.Context) {
	c.Request.ParseForm()

	defaultValue := []byte("input{ }\n\nfilter{ }\n\noutput{ }")
	var p = models.Pipeline{
		Label:       c.Request.PostFormValue("label"),
		Description: c.Request.PostFormValue("description"),
		Assets: []models.Asset{{
			Name:        "bitfan.conf",
			Type:        "entrypoint",
			ContentType: "text/plain",
			Value:       defaultValue,
			Size:        len(defaultValue),
		}},
	}

	pnew, _ := apiClient.NewPipeline(&p)

	session := sessions.Default(c)
	session.AddFlash(fmt.Sprintf("Pipeline %s created", p.Label))
	session.Save()

	c.Redirect(302, fmt.Sprintf("/ui/pipelines/%s", pnew.Uuid))
}

func updatePipeline(c *gin.Context) {
	c.Request.ParseForm()
	pipelineUUID := c.Param("id")

	var data = map[string]interface{}{
		"label":       c.Request.PostFormValue("label"),
		"description": c.Request.PostFormValue("description"),
		"auto_start":  false,
	}

	if _, ok := c.Request.PostForm["auto_start"]; ok {
		data["auto_start"] = true
	}

	pnew, _ := apiClient.UpdatePipeline(pipelineUUID, &data)
	session := sessions.Default(c)
	session.AddFlash(fmt.Sprintf("Pipeline %s saved", pnew.Label))
	session.Save()

	c.Redirect(302, fmt.Sprintf("/ui/pipelines/%s", pnew.Uuid))
}

func startPipeline(c *gin.Context) {
	pipelineUUID := c.Param("id")

	pipeline, err := apiClient.StartPipeline(pipelineUUID)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
	} else {
		c.JSON(200, pipeline)
		log.Printf("Started (UUID:%s) - %s\n", pipeline.Uuid, pipeline.Label)
	}
}

func stopPipeline(c *gin.Context) {
	pipelineUUID := c.Param("id")

	pipeline, err := apiClient.StopPipeline(pipelineUUID)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
	} else {
		c.JSON(200, pipeline)
		log.Printf("Stopped (UUID:%s) - %s\n", pipeline.Uuid, pipeline.Label)
	}
}

func deletePipeline(c *gin.Context) {
	pipelineUUID := c.Param("id")

	err := apiClient.DeletePipeline(pipelineUUID)

	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.Redirect(302, "/ui/pipelines")
}

func createAsset(c *gin.Context) {
	pipelineUUID := c.Param("id")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	filename := header.Filename

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	nasset := &models.Asset{
		Name:         filename,
		Value:        buf.Bytes(),
		PipelineUUID: pipelineUUID,
	}

	asset, err := apiClient.NewAsset(nasset)

	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.JSON(201, asset)
}

func showAsset(c *gin.Context) {
	pipelineUUID := c.Param("id")
	assetUUID := c.Param("assetID")

	p, _ := apiClient.Pipeline(pipelineUUID)
	a, _ := apiClient.Asset(assetUUID)

	flashes := []string{}
	for _, m := range sessions.Default(c).Flashes() {
		flashes = append(flashes, m.(string))
	}
	sessions.Default(c).Save()

	c.HTML(200, "assets/edit", gin.H{
		"asset":    a,
		"pipeline": p,
		"flashes":  flashes,
	})
}

func deleteAsset(c *gin.Context) {
	pipelineUUID := c.Param("id")
	assetUUID := c.Param("assetID")

	err := apiClient.DeleteAsset(assetUUID)

	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.Redirect(302, fmt.Sprintf("/ui/pipelines/%s", pipelineUUID))
}

func downloadAsset(c *gin.Context) {
	assetUUID := c.Param("assetID")

	asset, err := apiClient.Asset(assetUUID)

	if err != nil {
		c.JSON(404, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\""+asset.Name+"\"")
	c.Data(200, asset.ContentType, asset.Value)
}

func updateAsset(c *gin.Context) {
	c.Request.ParseForm()
	pipelineUUID := c.Param("id")
	assetUUID := c.Param("assetID")

	var data = map[string]interface{}{
		"name": c.Request.PostFormValue("name"),
	}
	if _, ok := c.Request.PostForm["content"]; ok {
		data["value"] = []byte(c.Request.PostFormValue("content"))
	}

	asset, err := apiClient.UpdateAsset(assetUUID, &data)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	session := sessions.Default(c)
	session.AddFlash(fmt.Sprintf("Asset %s saved", asset.Name))
	session.Save()

	c.Redirect(302, fmt.Sprintf("/ui/pipelines/%s/assets/%s", pipelineUUID, assetUUID))
}

func replaceAsset(c *gin.Context) {
	pipelineUUID := c.Param("id")
	assetUUID := c.Param("assetID")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	filename := header.Filename

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	nasset := &models.Asset{
		Uuid:         assetUUID,
		Name:         filename,
		Value:        buf.Bytes(),
		PipelineUUID: pipelineUUID,
	}

	asset, err := apiClient.ReplaceAsset(assetUUID, nasset)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	session := sessions.Default(c)
	session.AddFlash(fmt.Sprintf("Asset %s saved", asset.Name))
	session.Save()

	c.JSON(200, asset)
}

/*


func replaceAsset(c *gin.Context) {
	assetID, err := strconv.Atoi(c.Param("assetID"))
	if err != nil {
		pp.Println("c-->", c.Param("assetID"))
		c.Abort()
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	filename := header.Filename
	size := header.Size

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		c.Abort()
		return
	}
	// Reset the read pointer if necessary.
	file.Seek(0, 0)
	contentType := http.DetectContentType(buffer[:n])

	pp.Println(header)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Fatal(err)
	}

	var a = Asset{ID: assetID}
	db.First(&a)

	a.Name = filename
	a.ContentType = contentType
	a.Value = buf.Bytes()
	a.Size = int(size)

	db.Save(&a)
	c.String(200, "ok")
}













*/
