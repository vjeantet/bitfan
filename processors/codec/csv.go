package codec

import (
	"encoding/csv"
	"fmt"
	"io"
)

type csvDecoder struct {
	more        bool
	r           *csv.Reader
	columnnames []string
}

func NewCsvDecoder(r io.Reader) Decoder {
	d := &csvDecoder{
		r:    csv.NewReader(r),
		more: true,
	}
	d.r.Comma = ','
	return d
}

func (c *csvDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	record, err := csv.NewReader(r).Read()
	if err == io.EOF {
		return data, err
	}
	if c.columnnames == nil {
		c.columnnames = record
		return nil, nil
	}

	for i, v := range c.columnnames {
		data[v] = record[i]
		// data[fmt.Sprintf("col_%d", i)] = v
	}

	return data, nil
}

func (c *csvDecoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	record, err := c.r.Read()
	if err == io.EOF {
		c.more = false
		return data, err
	}
	for i, v := range record {
		data[fmt.Sprintf("col_%d", i)] = v
	}
	return data, nil
}

func (c *csvDecoder) More() bool {
	return c.more
}
