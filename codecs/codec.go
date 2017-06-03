package codecs

import (
	"bytes"
	"io"

	"golang.org/x/net/html/charset"

	"github.com/vjeantet/bitfan/codecs/csvcodec"
	"github.com/vjeantet/bitfan/codecs/jsoncodec"
	"github.com/vjeantet/bitfan/codecs/jsonlinescodec"
	"github.com/vjeantet/bitfan/codecs/linecodec"
	"github.com/vjeantet/bitfan/codecs/multilinecodec"
	"github.com/vjeantet/bitfan/codecs/plaincodec"
	"github.com/vjeantet/bitfan/core/config"
)

type Codec struct {
	Name    string
	Charset string
	Options map[string]interface{}
}

func New(name string) Codec {
	return Codec{
		Name:    name,
		Charset: "utf-8",
		Options: map[string]interface{}{},
	}
}

func NewFromConfig(conf *config.Codec) Codec {
	name := conf.Name
	c := New(name)
	for i, k := range conf.Options {
		if i == "charset" {
			c.Charset = k.(string)
			continue
		}
		c.Options[i] = k
	}
	return c
}

func (c *Codec) String() string {
	return c.Name
}

func (c *Codec) Decoder(r io.Reader) (Decoder, error) {
	var dec Decoder

	if r == nil {
		r = bytes.NewReader(nil)
	}

	var cr io.Reader
	var err error
	cr, err = charset.NewReaderLabel(c.Charset, r)
	if err != nil {
		return dec, err
	}

	switch c.Name {
	case "line": // OK
		dec = linecodec.New(cr, c.Options)
	case "multiline": // OK
		dec = multilinecodec.New(cr, c.Options)
	case "csv": // OK
		dec = csvcodec.New(cr, c.Options)
	case "json": // OK
		dec = jsoncodec.New(cr, c.Options)
	case "json_lines": // OK
		dec = jsonlinescodec.New(cr, c.Options)
	default:
		dec = plaincodec.New(cr, c.Options)
	}
	return dec, nil
}
