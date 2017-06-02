package codec

import (
	"encoding/json"
	"io"

	"golang.org/x/net/html/charset"
)

type jsonDecoder struct {
	d       *json.Decoder
	options map[string]interface{}
}

func NewJsonDecoder(r io.Reader, opt map[string]interface{}) Decoder {
	return &jsonDecoder{
		d:       json.NewDecoder(r),
		options: opt,
	}
}
func (p *jsonDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
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

	err := json.NewDecoder(cr).Decode(&data)
	return data, err
}

func (p *jsonDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	err := p.d.Decode(&data)

	return data, err
}
func (p *jsonDecoder) More() bool {
	return p.d.More()
}
