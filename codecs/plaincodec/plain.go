//go:generate bitfanDoc -codec plainDecoder
package plaincodec

import (
	"io"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
)

type plainDecoder struct {
	more    bool
	r       io.Reader
	options options
}

type options struct {
}

func New(r io.Reader, opt map[string]interface{}) *plainDecoder {
	d := &plainDecoder{
		r:       r,
		more:    true,
		options: options{},
	}

	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}

	return d
}

func (p *plainDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	p.more = false
	bytes, err := ioutil.ReadAll(p.r)
	if err != nil {
		return data, err
	}
	data["message"] = string(bytes)
	return data, nil
}

func (p *plainDecoder) More() bool {
	return p.more
}
