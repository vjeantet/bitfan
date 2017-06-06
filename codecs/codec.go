package codecs

import (
	"fmt"
	"io"

	"golang.org/x/net/html/charset"

	"github.com/vjeantet/bitfan/codecs/csv"
	"github.com/vjeantet/bitfan/codecs/json"
	"github.com/vjeantet/bitfan/codecs/jsonlines"
	"github.com/vjeantet/bitfan/codecs/line"
	"github.com/vjeantet/bitfan/codecs/multiline"
	"github.com/vjeantet/bitfan/codecs/plain"
	"github.com/vjeantet/bitfan/codecs/rubydebug"
)

type Codec struct {
	Name    string
	Charset string
	Options map[string]interface{}
}

func New(name string, conf map[string]interface{}) Codec {
	c := Codec{
		Name:    name,
		Charset: "utf-8",
		Options: map[string]interface{}{},
	}
	for i, k := range conf {
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

func (c *Codec) NewEncoder(w io.Writer) (Encoder, error) {
	var enc Encoder

	if w == nil {
		return enc, fmt.Errorf("codecs.Codec.Encoder error : no writer !")
	}

	// charset ?
	switch c.Name {
	case "pp":
		enc = rubydebugcodec.NewEncoder(w)
	case "rubydebug":
		enc = rubydebugcodec.NewEncoder(w)
	case "line":
		enc = linecodec.NewEncoder(w)
	case "json":
		enc = jsoncodec.NewEncoder(w)
	default:
		return enc, fmt.Errorf("no encoder defined")
	}
	enc.SetOptions(c.Options)

	return enc, nil
}

func (c *Codec) NewDecoder(r io.Reader) (Decoder, error) {
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
		dec = linecodec.NewDecoder(cr)
	case "multiline":
		dec = multilinecodec.NewDecoder(cr)
	case "csv":
		dec = csvcodec.NewDecoder(cr)
	case "json":
		dec = jsoncodec.NewDecoder(cr)
	case "json_lines":
		dec = jsonlinescodec.NewDecoder(cr)
	default:
		dec = plaincodec.NewDecoder(cr)
	}
	dec.SetOptions(c.Options)

	return dec, nil
}
