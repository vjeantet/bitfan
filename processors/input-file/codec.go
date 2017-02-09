package file

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

var EOF = fmt.Errorf("End of file")
var NOMOREDATA = fmt.Errorf("No more data")

type Decoder interface {
	Decode() (map[string]interface{}, error)
	More() bool
}

type csvDecoder struct {
	more bool
	r    *csv.Reader
}

func NewCsvDecoder(r io.Reader) Decoder {
	d := &csvDecoder{
		r:    csv.NewReader(r),
		more: true,
	}
	d.r.Comma = ';'
	return d
}
func (c *csvDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	record, err := c.r.Read()
	if err == io.EOF {
		c.more = false
		return data, err
	}
	data["values"] = record
	return data, nil
}
func (c *csvDecoder) More() bool {
	return c.more
}

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

type jsonDecoder struct {
	d *json.Decoder
}

func NewJsonDecoder(r io.Reader) Decoder {
	return &jsonDecoder{
		d: json.NewDecoder(r),
	}
}
func (p *jsonDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	err := p.d.Decode(&data)

	return data, err
}
func (p *jsonDecoder) More() bool {
	return p.d.More()
}
