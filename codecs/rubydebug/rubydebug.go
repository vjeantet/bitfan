//go:generate bitfanDoc -codec codec
package rubydebugcodec

import (
	"io"

	"github.com/k0kubun/pp"
	"github.com/mitchellh/mapstructure"
)

type codec struct {
	w       io.Writer
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

func (r *codec) Encoder(w io.Writer) *codec {
	r.w = w
	return r
}

func (p *codec) Encode(data map[string]interface{}) error {
	pp.Fprintln(p.w, data)
	return nil
}
