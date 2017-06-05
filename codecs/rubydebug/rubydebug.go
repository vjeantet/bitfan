//go:generate bitfanDoc -codec encoder,decoder
package rubydebugcodec

import (
	"io"

	"github.com/k0kubun/pp"
	"github.com/mitchellh/mapstructure"
)

type encoder struct {
	w       io.Writer
	options options
}

type options struct {
}

func NewEncoder(w io.Writer) *encoder {
	return &encoder{
		w:       w,
		options: options{},
	}
}

func (e *encoder) SetOptions(conf map[string]interface{}) error {
	if err := mapstructure.Decode(conf, &e.options); err != nil {
		return err
	}
	return nil
}

func (e *encoder) Encode(data map[string]interface{}) error {
	pp.Fprintln(e.w, data)
	return nil
}
