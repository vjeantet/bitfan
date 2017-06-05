//go:generate bitfanDoc -codec codec
package plaincodec

import (
	"io"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
)

type decoder struct {
	more    bool
	r       io.Reader
	options options
}

type options struct {
}

func NewDecoder(r io.Reader) *decoder {
	d := &decoder{
		r:       r,
		more:    true,
		options: options{},
	}

	return d
}

func (d *decoder) SetOptions(conf map[string]interface{}) error {

	if err := mapstructure.Decode(conf, &d.options); err != nil {
		return err
	}

	return nil
}

func (d *decoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	d.more = false
	bytes, err := ioutil.ReadAll(d.r)
	if err != nil {
		return data, err
	}
	data["message"] = string(bytes)
	return data, nil
}

func (d *decoder) More() bool {
	return d.more
}
