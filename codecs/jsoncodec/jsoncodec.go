package jsoncodec

import (
	"encoding/json"
	"io"
)

type jsonDecoder struct {
	d       *json.Decoder
	options map[string]interface{}
}

func New(r io.Reader, opt map[string]interface{}) *jsonDecoder {
	return &jsonDecoder{
		d:       json.NewDecoder(r),
		options: opt,
	}
}

func (p *jsonDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	err := p.d.Decode(&data)

	return data, err
}
func (p *jsonDecoder) More() bool {
	return p.d.More()
}
