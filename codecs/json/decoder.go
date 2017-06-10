//go:generate bitfanDoc -codec json
package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/codecs/lib"
)

type decoder struct {
	d       *json.Decoder
	options decoderOptions

	log lib.Logger
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
func (d *decoder) SetOptions(conf map[string]interface{}, logger lib.Logger, cwl string) error {
	d.log = logger

	if err := mapstructure.Decode(conf, &d.options); err != nil {
		return err
	}
	return nil
}

func (d *decoder) Decode(v *interface{}) error {

	err := d.d.Decode(v)

	return err
}

func (d *decoder) More() bool {
	return d.d.More()
}
