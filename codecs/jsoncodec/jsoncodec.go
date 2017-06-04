//go:generate bitfanDoc -codec codec
package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
)

type codec struct {
	d       *json.Decoder
	options options
}

type options struct {
}

func New(opt map[string]interface{}) *codec {
	d := &codec{
		options: options{},
	}
	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}
	return d
}
func (c *codec) Decoder(r io.Reader) *codec {
	c.d = json.NewDecoder(r)
	return c
}

func (p *codec) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	err := p.d.Decode(&data)

	return data, err
}
func (p *codec) More() bool {
	return p.d.More()
}
