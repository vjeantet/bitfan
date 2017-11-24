//go:generate bitfanDoc -codec json
package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/commons"
)

type decoder struct {
	d       *json.Decoder
	options decoderOptions

	log commons.Logger
}

type decoderOptions struct {
	// Set indentation
	// @Default ""
	// @ExampleLS indent => "    "
	Indent string `mapstructure:"indent"`

	// Json is an array, decode each element as a distinct dataframe
	// @Default false
	StreamArray bool `mapstructure:"stream_array"`
}

func NewDecoder(r io.Reader) *decoder {
	return &decoder{
		d:       json.NewDecoder(r),
		options: decoderOptions{},
	}
}
func (d *decoder) SetOptions(conf map[string]interface{}, logger commons.Logger, cwl string) error {
	d.log = logger

	err := mapstructure.Decode(conf, &d.options)

	if d.options.StreamArray == true {
		// read open bracket
		d.d.Token()
	}

	return err
}

func (d *decoder) Decode(v *interface{}) error {
	return d.d.Decode(v)
}

func (d *decoder) More() bool {
	return d.d.More()
}
