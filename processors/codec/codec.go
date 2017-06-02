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

func (c *Codec) String() string {
	return c.Name
}

func (c *Codec) Decoder() (Decoder, error) {
	// return NewDecoder(c.Name)
	var dec Decoder

	var cr io.Reader
	//todo get Charset from Codec settings
	if char7, ok := c.Options["charset"]; ok {
		var err error
		r1 := bytes.NewReader(nil)
		cr, err = charset.NewReaderLabel(char7.(string), r1)
		if err != nil {
			return nil, err
		}
	} else {
		cr = bytes.NewReader(nil)
	}

	switch c.Name {
	case "json":
		dec = NewJsonDecoder(cr, c.Options)
	case "csv":
		dec = NewCsvDecoder(cr, c.Options)
	default:
		dec = NewPlainDecoder(cr, c.Options)
	}
	return dec, nil
}

func NewFromConfig(conf *config.Codec) Codec {
	name := conf.Name
	c := New(name)
	for i, k := range conf.Options {
		c.Options[i] = k
	}
	return c
}
