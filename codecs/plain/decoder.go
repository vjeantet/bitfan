//go:generate bitfanDoc -codec plain
package plaincodec

import (
	"io"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/codecs/lib"
)

type decoder struct {
	more    bool
	r       io.Reader
	options decoderOptions

	log lib.Logger
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

func (d *decoder) SetOptions(conf map[string]interface{}, logger lib.Logger, cwl string) error {
	d.log = logger

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
