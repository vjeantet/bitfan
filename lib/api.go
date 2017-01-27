package lib

import (
	"log"

	"github.com/rsms/gotalk"
	config "github.com/veino/veino/config"
	"github.com/veino/veino/runtime"
)

type ApiStarter struct {
	Pipeline config.Pipeline
	Agents   []config.Agent
}

func ApiServe(veinoHost string) {
	// gotalk.Handle("findPipelines", func() (config.PipelineList, error) {
	// 	return config.Pipelines(), nil
	// })

	gotalk.Handle("stopPipeline", func(ID int) (bool, error) {
		err := runtime.StopPipeline(ID)
		return true, err
	})

	gotalk.Handle("startPipeline", func(starter ApiStarter) (int, error) {
		ID, err := runtime.StartPipeline(&starter.Pipeline, starter.Agents)
		return ID, err
	})

	s, err := gotalk.Listen("tcp", veinoHost)
	if err != nil {
		log.Fatalln(err)
	}
	go s.Accept()
}

func ApiClient(veinoHost string) *gotalk.Sock {
	s, err := gotalk.Connect("tcp", veinoHost)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}
