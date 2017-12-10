package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"github.com/vjeantet/bitfan/api/models"
)

type playgroundRequest struct {
	Event      string `json:"event"`
	EventType  string `json:"event_type"`
	FilterPart string `json:"filter"`
	UUID       string `json:"uuid"`
}

func playgroundsFilter(c *gin.Context) {
	c.HTML(200, "playgrounds/filters", withCommonValues(c, gin.H{}))
}
func playgroundsFilterExit(c *gin.Context) {
	pgReq := playgroundRequest{}
	_ = c.BindJSON(&pgReq)
	pp.Println("bye-->", pgReq.UUID)

	// Stop pipeline if running
	_, _ = apiClient.StopPipeline(pgReq.UUID)
	c.JSON(200, gin.H{"ok": "ok"})
}

func playgroundsFilterDo(c *gin.Context) {
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
	pgFullConfig := `input{
  httpserver{
	codec => ` + pgReq.EventType + ` 		
  }
  websocket wsin{
  	codec => ` + pgReq.EventType + ` 
  	uri => "wsin"
  }
}
filter{
` + pgReq.FilterPart + `
} 
output{
  websocket wsout{
  	codec => json {indent => "    "}
  	uri => "wsout"
  }
}`

	// Stop pipeline if running
	_, _ = apiClient.StopPipeline(pgReq.UUID)

	// start pipeline
	defaultValue := []byte(pgFullConfig)
	var p = models.Pipeline{
		Playground:  true,
		Uuid:        pgReq.UUID,
		Active:      true,
		Label:       "playground-filter " + pgReq.UUID,
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
		case "httpserver":
			httpin = wh.Url
		}
	}

	c.JSON(200, withCommonValues(c, gin.H{
		"wsin":   wsin,
		"wsout":  wsout,
		"httpin": httpin,
	}))
}
