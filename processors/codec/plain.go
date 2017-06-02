package codec

import (
	"io"
	"io/ioutil"
)

type plainDecoder struct {
	more bool
	r    io.Reader
}

func NewPlainDecoder(r io.Reader) Decoder {
	return &plainDecoder{
		r:    r,
		more: true,
	}
}
func (p *plainDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return data, err
	}
	data["message"] = string(bytes)
	return data, nil
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
