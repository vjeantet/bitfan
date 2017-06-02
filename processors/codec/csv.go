package codec

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"

	"golang.org/x/net/html/charset"
)

type csvDecoder struct {
	more        bool
	r           *csv.Reader
	columnnames []string
	options     csvDecoderOptions
	comma       rune
}

type csvDecoderOptions struct {
	Charset   string
	Separator string
}

func NewCsvDecoder(r io.Reader, opt map[string]interface{}) Decoder {
	d := &csvDecoder{
		r:    csv.NewReader(r),
		more: true,
		options: csvDecoderOptions{
			Charset:   "utf-8",
			Separator: "-",
		},
		comma: '-',
	}

	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}

	d.r.Comma = []rune(d.options.Separator)[0]
	d.comma = d.r.Comma
	return d
}

func (c *csvDecoder) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	var cr io.Reader

	var err error
	cr, err = charset.NewReaderLabel(c.options.Charset, r)
	if err != nil {
		return nil, err
	}
	csvr := csv.NewReader(cr)
	csvr.Comma = c.comma

	record, err := csvr.Read()

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
