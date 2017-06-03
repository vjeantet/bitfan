package plaincodec

import (
	"io"
	"io/ioutil"

	"golang.org/x/net/html/charset"
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

func (p *plainDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	var cr io.Reader

	if char7, ok := p.options["charset"]; ok {
		var err error
		cr, err = charset.NewReaderLabel(char7.(string), r)
		if err != nil {
			return nil, err
		}
	} else {
		cr = r
	}

	bytes, err := ioutil.ReadAll(cr)
	if err != nil {
		return data, err
	}
	data["message"] = string(bytes)
	return data, nil
}
