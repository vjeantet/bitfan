//go:generate bitfanDoc -codec json
package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
)

type encoder struct {
	e       *json.Encoder
	options options
}

type decoder struct {
	d       *json.Decoder
	options options
}

type options struct {
	// Set indentation
	// @Default ""
	// @ExampleLS indent => "    "
	Indent string `mapstructure:"indent"`
}

func NewEncoder(w io.Writer) *encoder {
	e := &encoder{
		e: json.NewEncoder(w),
		options: options{
			Indent: "",
		},
	}

	return e
}

func (e *encoder) SetOptions(opt map[string]interface{}) error {
	if err := mapstructure.Decode(opt, &e.options); err != nil {
		return err
	}
	e.e.SetIndent("", e.options.Indent)
	return nil
}

func (e *encoder) Encode(data map[string]interface{}) error {
	e.e.Encode(data)
	return nil
}

func NewDecoder(r io.Reader) *decoder {
	return &decoder{
		d:       json.NewDecoder(r),
		options: options{},
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
