//go:generate bitfanDoc -codec codec
package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
)

type codec struct {
	d       *json.Decoder
	e       *json.Encoder
	options options
}

type options struct {
	// Set indentation
	// @Default ""
	// @ExampleLS indent => "    "
	Indent string `mapstructure:"indent"`
}

func New(opt map[string]interface{}) *codec {
	d := &codec{
		options: options{
			Indent: "",
		},
	}
	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}
	return d
}

func (c *codec) Encoder(w io.Writer) *codec {
	c.e = json.NewEncoder(w)
	c.e.SetIndent("", c.options.Indent)
	return c
}

func (c *codec) Decoder(r io.Reader) *codec {
	c.d = json.NewDecoder(r)
	return c
}

func (p *codec) Encode(data map[string]interface{}) error {
	p.e.Encode(data)
	return nil
}

func (p *codec) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	err := p.d.Decode(&data)

	return data, err
}

func (p *codec) More() bool {
	return p.d.More()
}
