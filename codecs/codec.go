package codecs

import (
	"fmt"
	"io"

	"golang.org/x/net/html/charset"

	"github.com/vjeantet/bitfan/codecs/csvcodec"
	"github.com/vjeantet/bitfan/codecs/jsoncodec"
	"github.com/vjeantet/bitfan/codecs/jsonlinescodec"
	"github.com/vjeantet/bitfan/codecs/linecodec"
	"github.com/vjeantet/bitfan/codecs/multilinecodec"
	"github.com/vjeantet/bitfan/codecs/plaincodec"
	"github.com/vjeantet/bitfan/codecs/rubydebugcodec"
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

func (c *Codec) Encoder(w io.Writer) (Encoder, error) {
	var enc Encoder

	if w == nil {
		return enc, fmt.Errorf("codecs.Codec.Encoder error : no writer !")
	}

	// charset ?
	switch c.Name {
	case "pp":
		enc = rubydebugcodec.New(c.Options).Encoder(w)
	case "rubydebug":
		enc = rubydebugcodec.New(c.Options).Encoder(w)
	case "line":
		enc = linecodec.New(c.Options).Encoder(w)
		//TODO default
	}

	return enc, nil
}

func (c *Codec) Decoder(r io.Reader) (Decoder, error) {
	var dec Decoder

	if r == nil {
		return dec, fmt.Errorf("codecs.Codec.Decoder error : no reader !")
	}

	var cr io.Reader
	var err error
	cr, err = charset.NewReaderLabel(c.Charset, r)
	if err != nil {
		return dec, err
	}

	switch c.Name {
	case "line":
		dec = linecodec.New(c.Options).Decoder(cr)
	case "multiline":
		dec = multilinecodec.New(c.Options).Decoder(cr)
	case "csv":
		dec = csvcodec.New(c.Options).Decoder(cr)
	case "json":
		dec = jsoncodec.New(c.Options).Decoder(cr)
	case "json_lines":
		dec = jsonlinescodec.New(c.Options).Decoder(cr)
	default:
		dec = plaincodec.New(c.Options).Decoder(cr)
	}

	return dec, nil
}
