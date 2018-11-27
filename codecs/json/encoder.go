package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
	"bitfan/commons"
)

type encoder struct {
	e       *json.Encoder
	options encoderOptions
}

type encoderOptions struct {
	// Set indentation
	// @Default ""
	// @ExampleLS indent => "    "
	Indent string `mapstructure:"indent"`
}

func NewEncoder(w io.Writer) *encoder {
	e := &encoder{
		e: json.NewEncoder(w),
		options: encoderOptions{
			Indent: "",
		},
	}

	return e
}

func (e *encoder) SetOptions(conf map[string]interface{}, logger commons.Logger, cwl string) error {
	if err := mapstructure.Decode(conf, &e.options); err != nil {
		return err
	}
	e.e.SetIndent("", e.options.Indent)
	return nil
}

func (e *encoder) Encode(data map[string]interface{}) error {
	e.e.Encode(data)
	return nil
}
