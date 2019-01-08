//go:generate bitfanDoc -codec plain
package plaincodec

import (
	"io"
	"io/ioutil"

	"github.com/awillis/bitfan/commons"
	"github.com/mitchellh/mapstructure"
)

type decoder struct {
	more    bool
	r       io.Reader
	options decoderOptions

	log commons.Logger
}

type decoderOptions struct {
}

func NewDecoder(r io.Reader) *decoder {
	d := &decoder{
		r:       r,
		more:    true,
		options: decoderOptions{},
	}

	return d
}

func (d *decoder) SetOptions(conf map[string]interface{}, logger commons.Logger, cwl string) error {
	d.log = logger

	return mapstructure.Decode(conf, &d.options)
}

func (d *decoder) Decode(v *interface{}) error {
	*v = nil
	d.more = false
	bytes, err := ioutil.ReadAll(d.r)
	if err != nil {
		return err
	}
	*v = string(bytes)
	return nil
}

func (d *decoder) More() bool {
	return d.more
}

func (d *decoder) Buffer() []byte {
	return []byte{}
}
