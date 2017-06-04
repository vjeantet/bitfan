//go:generate bitfanDoc -codec codec
package plaincodec

import (
	"io"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
)

type codec struct {
	more    bool
	r       io.Reader
	options options
}

type options struct {
}

func New(opt map[string]interface{}) *codec {
	d := &codec{
		more:    true,
		options: options{},
	}
	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}
	return d
}
func (c *codec) Decoder(r io.Reader) *codec {
	c.r = r
	return c
}

func (p *codec) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	p.more = false
	bytes, err := ioutil.ReadAll(p.r)
	if err != nil {
		return data, err
	}
	data["message"] = string(bytes)
	return data, nil
}

func (p *codec) More() bool {
	return p.more
}
