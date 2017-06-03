package codec

import (
	"bytes"
	"io"

	"github.com/vjeantet/bitfan/core/config"
	"golang.org/x/net/html/charset"
)

type Codec struct {
	Name    string
	Options map[string]interface{}
}

func New(name string) Codec {
	return Codec{
		Name:    name,
		Options: map[string]interface{}{},
	}
}

func NewFromConfig(conf *config.Codec) Codec {
	name := conf.Name
	c := New(name)
	for i, k := range conf.Options {
		c.Options[i] = k
	}
	return c
}

func (c *Codec) String() string {
	return c.Name
}

func (c *Codec) Decoder(f io.Reader) (Decoder, error) {
	var dec Decoder
	var cr io.Reader

	if f == nil {
		f = bytes.NewReader(nil)
	}

	//todo get Charset from Codec settings

	var err error
	var charsetLabel = "utf-8"
	if char7, ok := c.Options["charset"]; ok {
		charsetLabel = char7.(string)
	}

	cr, err = charset.NewReaderLabel(charsetLabel, f)
	if err != nil {
		return nil, err
	}

	switch c.Name {
	case "json":
		dec = NewJsonDecoder(cr, c.Options)
	case "csv":
		dec = NewCsvDecoder(cr, c.Options)
	case "multiline":
		dec = NewMultilineDecoder(cr, c.Options)
	default:
		dec = NewPlainDecoder(cr, c.Options)
	}
	return dec, nil
}
