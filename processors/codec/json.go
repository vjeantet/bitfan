package codec

import (
	"encoding/json"
	"io"
)

type jsonDecoder struct {
	d *json.Decoder
}

func NewJsonDecoder(r io.Reader) Decoder {
	return &jsonDecoder{
		d: json.NewDecoder(r),
	}
}
func (p *jsonDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json.NewDecoder(r).Decode(&data)
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
