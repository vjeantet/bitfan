//go:generate bitfanDoc -codec json
package jsoncodec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/awillis/bitfan/commons"
	"github.com/mitchellh/mapstructure"
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
	*v = nil
	return d.d.Decode(v)
}

func (d *decoder) More() bool {
	return d.d.More()
}

func (d *decoder) Buffer() []byte {
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, d.d.Buffered())
	if err != nil {
		fmt.Println(err)
	}

	return buf.Bytes()
}
