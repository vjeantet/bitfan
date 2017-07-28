package api

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/k0kubun/pp"
	"github.com/vjeantet/bitfan/processors/doc"
)

type RestClient struct {
	host string
}

func NewRestClient(bitfanHost string) *RestClient {
	cli := &RestClient{
		host: "http://" + bitfanHost + "/api/v1/",
	}
	return cli
}

func (r *RestClient) client() *sling.Sling {
	return sling.New().Base(r.host)
}

func (r *RestClient) ListPipelines() ([]*Pipeline, error) {
	pipelines := new([]*Pipeline)
	apierror := new(Error)

	resp, err := r.client().Get("pipelines").Receive(pipelines, apierror)
	if err != nil {
		return *pipelines, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return *pipelines, err
}

func (r *RestClient) StopPipeline(ID string) error {
	retour := new(map[string]interface{})
	apierror := new(Error)
	resp, err := r.client().Delete("pipelines/"+ID).Receive(retour, apierror)

	if err != nil {
		return err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return err
}

func (r *RestClient) AddPipeline(pipeline *Pipeline) (*Pipeline, error) {
	newPipeline := &Pipeline{}
	apierror := new(Error)
	resp, err := r.client().Post("pipelines").BodyJSON(pipeline).Receive(newPipeline, apierror)

	if err != nil {
		return newPipeline, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return newPipeline, err
}

func (r *RestClient) Pipeline(ID string, full bool) (*Pipeline, error) {
	pipeline := &Pipeline{}
	apierror := new(Error)

	resp, err := r.client().Get("pipelines/"+ID).Receive(pipeline, apierror)

	if err != nil {
		return pipeline, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
		return pipeline, err
	}

	if full {
		assets := new([]Asset)
		resp, err = r.client().Get("pipelines/"+ID+"/assets").Receive(assets, apierror)
		if err != nil {
			return pipeline, err
		} else if resp.StatusCode > 400 {
			err = fmt.Errorf(apierror.Message)
			return pipeline, err
		}

		pipeline.Assets = *assets
	}

	return pipeline, err
}

func (r *RestClient) ListDoc() error {
	docs := make(map[string]map[string]*doc.Processor)

	apierror := new(Error)
	resp, err := r.client().Get("docs").Receive(&docs, apierror)
	pp.Println("docs-->", docs)
	if err != nil {
		return err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return err
}

func (r *RestClient) Test() {

}
