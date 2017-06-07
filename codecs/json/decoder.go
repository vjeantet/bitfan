//go:generate bitfanDoc -codec json
package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
)

type decoder struct {
	d       *json.Decoder
	options decoderOptions
}

type decoderOptions struct {
	// Set indentation
	// @Default ""
	// @ExampleLS indent => "    "
	Indent string `mapstructure:"indent"`
}

func NewDecoder(r io.Reader) *decoder {
	return &decoder{
		d:       json.NewDecoder(r),
		options: decoderOptions{},
	}
}
func (d *decoder) SetOptions(opt map[string]interface{}) error {
	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return err
	}
	return nil
}

func (d *decoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	err := d.d.Decode(&data)

	return data, err
}

func (d *decoder) More() bool {
	return d.d.More()
}
