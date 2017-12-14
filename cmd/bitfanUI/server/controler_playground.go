package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"github.com/vjeantet/bitfan/api/models"
)

type playgroundRequest struct {
	UUID        string `json:"uuid"`
	InputValue  string `json:"input_value"`
	InputMode   string `json:"input_mode"`
	InputCodec  string `json:"input_codec"`
	FilterValue string `json:"filter_value"`
	FilterMode  string `json:"filter_mode"`
	OutputValue string `json:"output_value"`
	OutputMode  string `json:"output_mode"`
}

func playgroundsPlay(c *gin.Context) {
	c.HTML(200, "playgrounds/play", withCommonValues(c, gin.H{}))
}
func playgroundsPlayExit(c *gin.Context) {
	pgReq := playgroundRequest{}
	_ = c.BindJSON(&pgReq)
	pp.Println("bye-->", pgReq.UUID)

	// Stop pipeline if running
	_, _ = apiClient.StopPipeline(pgReq.UUID)
	c.JSON(200, gin.H{"ok": "ok"})
}

func playgroundsPlayDo(c *gin.Context) {
	pgReq := playgroundRequest{}
	err := c.BindJSON(&pgReq)

	if pgReq.UUID == "" {
		c.JSON(400, err.Error())
		log.Printf("error : no uuid provided\n")
		return
	}
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	// Build a complete bitfan configuration
	// - with input as WS
	// - with output as WS
	pgFullConfig := "input{\n"

	if pgReq.InputMode == "raw" {
		pgFullConfig = pgFullConfig + "\n  websocket wsin{ codec => " + pgReq.InputCodec + " uri => wsin}\n"
	} else {
		pgFullConfig = pgFullConfig + pgReq.InputValue
	}

	pgFullConfig = pgFullConfig + "\n} filter{\n"

	if pgReq.FilterMode == "configuration" {
		pgFullConfig = pgFullConfig + pgReq.FilterValue
	}

	pgFullConfig = pgFullConfig + "\n} output{\n"

	if pgReq.OutputMode == "raw" {
		pgFullConfig = pgFullConfig + "  websocket wsout{ codec => json {indent => '    '} uri => wsout  }\n"
	} else {
		pgFullConfig = pgFullConfig + pgReq.OutputValue
	}

	pgFullConfig = pgFullConfig + "\n}"

	// Stop pipeline if running
	_, _ = apiClient.StopPipeline(pgReq.UUID)

	// start pipeline
	defaultValue := []byte(pgFullConfig)
	var p = models.Pipeline{
		Playground:  true,
		Uuid:        pgReq.UUID,
		Active:      true,
		Label:       "playground-" + pgReq.UUID,
		Description: "",
		Assets: []models.Asset{{
			Name:        "play.conf",
			Type:        "entrypoint",
			ContentType: "text/plain",
			Value:       defaultValue,
			Size:        len(defaultValue),
		}},
	}

	tp, err := apiClient.NewPipeline(&p)
	if err != nil {
		c.JSON(500, err.Error())
		log.Printf("error : %v\n", err)
		return
	}

	// get its UUID
	// build its WS IN and OUT
	// returns WS adresses to client
	var wsout, wsin, httpin string
	for _, wh := range tp.Webhooks {
		switch wh.Namespace {
		case "wsin":
			wsin = wh.Url
		case "wsout":
			wsout = wh.Url
		}
	}

	c.JSON(200, withCommonValues(c, gin.H{
		"wsin":  wsin,
		"wsout": wsout,
	}))
}
