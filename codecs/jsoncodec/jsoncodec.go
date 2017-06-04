//go:generate bitfanDoc -codec jsonDecoder
package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
)

type jsonDecoder struct {
	d       *json.Decoder
	options options
}

type options struct {
}

func New(r io.Reader, opt map[string]interface{}) *jsonDecoder {
	d := &jsonDecoder{
		d:       json.NewDecoder(r),
		options: options{},
	}

	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}

	return d
}

func (p *jsonDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	err := p.d.Decode(&data)

	return data, err
}
func (p *jsonDecoder) More() bool {
	return p.d.More()
}
