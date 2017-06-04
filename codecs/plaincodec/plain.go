package plaincodec

import (
	"io"
	"io/ioutil"
)

type plainDecoder struct {
	more    bool
	r       io.Reader
	options map[string]interface{}
}

func New(r io.Reader, opt map[string]interface{}) *plainDecoder {
	return &plainDecoder{
		r:       r,
		more:    true,
		options: opt,
	}
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
