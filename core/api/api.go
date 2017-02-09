package api

import (
	"log"

	"github.com/rsms/gotalk"
	"github.com/vjeantet/bitfan/core"
	config "github.com/vjeantet/bitfan/core/config"
)

type ApiStarter struct {
	Pipeline config.Pipeline
	Agents   []config.Agent
}

func ApiServe(bitfanHost string) {
	// GET /pipelines
	gotalk.Handle("findPipelines", func() (map[int]*core.Pipeline, error) {
		return core.Pipelines(), nil
	})

	// DELETE /pipelines/ID
	gotalk.Handle("stopPipeline", func(ID int) (bool, error) {
		err := core.StopPipeline(ID)
		return true, err
	})

	// POST /pipelines
	gotalk.Handle("startPipeline", func(starter ApiStarter) (int, error) {
		ID, err := core.StartPipeline(&starter.Pipeline, starter.Agents)
		return ID, err
	})

	s, err := gotalk.Listen("tcp", bitfanHost)
	if err != nil {
		log.Fatalln(err)
	}
	go s.Accept()
}

func ApiClient(bitfanHost string) *gotalk.Sock {
	s, err := gotalk.Connect("tcp", bitfanHost)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}
