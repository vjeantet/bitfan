package client

import (
	"fmt"

	"bitfan/api/models"
	"github.com/dghubble/sling"
)

type RestClient struct {
	host string
}

func New(bitfanHost string) *RestClient {
	cli := &RestClient{
		host: "http://" + bitfanHost + "/api/v2/",
	}
	return cli
}

func (r *RestClient) client() *sling.Sling {
	return sling.New().Base(r.host)
}

func (r *RestClient) Envs() ([]models.Env, error) {
	varenvs := make([]models.Env, 0)
	apierror := new(models.Error)

	resp, err := r.client().Get("env").Receive(&varenvs, apierror)

	if err != nil {
		return varenvs, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return varenvs, nil
}

func (r *RestClient) XProcessors(behavior string) ([]models.XProcessor, error) {
	xprocessors := make([]models.XProcessor, 0)
	apierror := new(models.Error)

	resp, err := r.client().Get("xprocessors?behavior="+behavior).Receive(&xprocessors, apierror)

	if err != nil {
		return xprocessors, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return xprocessors, nil
}

func (r *RestClient) XProcessor(ID string) (*models.XProcessor, error) {
	xprocessor := &models.XProcessor{}
	apierror := new(models.Error)

	resp, err := r.client().Get("xprocessors/"+ID).Receive(xprocessor, apierror)

	if err != nil {
		return xprocessor, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
		return xprocessor, err
	}

	return xprocessor, err
}

func (r *RestClient) Pipelines() ([]models.Pipeline, error) {
	pipelines := make([]models.Pipeline, 0)
	apierror := new(models.Error)

	resp, err := r.client().Get("pipelines").Receive(&pipelines, apierror)

	if err != nil {
		return pipelines, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return pipelines, err
}
func (r *RestClient) Pipeline(ID string) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	apierror := new(models.Error)

	resp, err := r.client().Get("pipelines/"+ID).Receive(pipeline, apierror)

	if err != nil {
		return pipeline, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
		return pipeline, err
	}

	return pipeline, err
}

func (r *RestClient) NewPipeline(pipeline *models.Pipeline) (*models.Pipeline, error) {
	newPipeline := new(models.Pipeline)
	apierror := new(models.Error)

	resp, err := r.client().Post("pipelines").BodyJSON(pipeline).Receive(newPipeline, apierror)

	if err != nil {
		return newPipeline, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return newPipeline, err
}

func (r *RestClient) NewXProcessor(xprocessor *models.XProcessor) (*models.XProcessor, error) {
	newXProcessor := new(models.XProcessor)
	apierror := new(models.Error)
	resp, err := r.client().Post("xprocessors").BodyJSON(xprocessor).Receive(newXProcessor, apierror)

	if err != nil {
		return newXProcessor, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return newXProcessor, err
}

func (r *RestClient) UpdateXProcessor(UUID string, data *map[string]interface{}) (*models.XProcessor, error) {
	newXProcessor := new(models.XProcessor)
	apierror := new(models.Error)

	resp, err := r.client().Patch("xprocessors/"+UUID).BodyJSON(data).Receive(newXProcessor, apierror)

	if err != nil {
		return newXProcessor, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return newXProcessor, err
}

func (r *RestClient) UpdatePipeline(UUID string, data *map[string]interface{}) (*models.Pipeline, error) {
	newPipeline := new(models.Pipeline)
	apierror := new(models.Error)

	resp, err := r.client().Patch("pipelines/"+UUID).BodyJSON(data).Receive(newPipeline, apierror)

	if err != nil {
		return newPipeline, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return newPipeline, err
}

func (r *RestClient) StartPipeline(UUID string) (*models.Pipeline, error) {
	newPipeline := new(models.Pipeline)
	apierror := new(models.Error)

	var data = map[string]interface{}{
		"active": true,
	}

	resp, err := r.client().Patch("pipelines/"+UUID).BodyJSON(data).Receive(newPipeline, apierror)
	if err != nil {
		return newPipeline, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return newPipeline, err
}

func (r *RestClient) StopPipeline(UUID string) (*models.Pipeline, error) {
	newPipeline := new(models.Pipeline)
	apierror := new(models.Error)

	var data = map[string]interface{}{
		"active": false,
	}

	resp, err := r.client().Patch("pipelines/"+UUID).BodyJSON(data).Receive(newPipeline, apierror)
	if err != nil {
		return newPipeline, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return newPipeline, err
}

func (r *RestClient) DeleteXProcessor(UUID string) error {
	apierror := new(models.Error)

	resp, err := r.client().Delete("xprocessors/"+UUID).Receive(nil, apierror)
	if err != nil {
		return err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return nil
}
func (r *RestClient) DeletePipeline(UUID string) error {
	apierror := new(models.Error)

	resp, err := r.client().Delete("pipelines/"+UUID).Receive(nil, apierror)
	if err != nil {
		return err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return nil
}

func (r *RestClient) NewAsset(initAsset *models.Asset) (*models.Asset, error) {
	asset := new(models.Asset)
	apierror := new(models.Error)

	resp, err := r.client().Post("assets").BodyJSON(initAsset).Receive(asset, apierror)
	if err != nil {
		return asset, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return asset, nil
}

func (r *RestClient) Asset(UUID string) (*models.Asset, error) {
	asset := new(models.Asset)
	apierror := new(models.Error)

	resp, err := r.client().Get("assets/"+UUID).Receive(asset, apierror)

	if err != nil {
		return asset, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
		return asset, err
	}

	return asset, err
}

func (r *RestClient) DeleteAsset(UUID string) error {
	apierror := new(models.Error)

	resp, err := r.client().Delete("assets/"+UUID).Receive(nil, apierror)
	if err != nil {
		return err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return nil
}

func (r *RestClient) UpdateAsset(UUID string, data *map[string]interface{}) (*models.Asset, error) {
	newAsset := new(models.Asset)
	apierror := new(models.Error)
	resp, err := r.client().Patch("assets/"+UUID).BodyJSON(data).Receive(newAsset, apierror)

	if err != nil {
		return newAsset, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return newAsset, err
}

func (r *RestClient) ReplaceAsset(UUID string, initAsset *models.Asset) (*models.Asset, error) {
	newAsset := new(models.Asset)
	apierror := new(models.Error)

	resp, err := r.client().Put("assets/"+UUID).BodyJSON(initAsset).Receive(newAsset, apierror)

	if err != nil {
		return newAsset, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}
	return newAsset, err
}

func (r *RestClient) CheckSyntax(asset *models.Asset) (map[string]interface{}, error) {
	apierror := new(models.Error)

	syntaxCheckResult := new(map[string]interface{})
	resp, err := r.client().Post("assets/0/syntax-check").BodyJSON(asset).Receive(syntaxCheckResult, apierror)
	if err != nil {
		return *syntaxCheckResult, err
	} else if resp.StatusCode > 400 {
		err = fmt.Errorf(apierror.Message)
	}

	return *syntaxCheckResult, nil
}

// func debug(r io.ReadCloser) string {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r)
// 	return buf.String()
// }
