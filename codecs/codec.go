package codecs

import (
	"fmt"
	"io"

	"bitfan/codecs/csv"
	"bitfan/codecs/json"
	"bitfan/codecs/jsonlines"
	"bitfan/codecs/line"
	"bitfan/codecs/multiline"
	"bitfan/codecs/plain"
	"bitfan/codecs/rubydebug"
	"bitfan/codecs/w3c"
	"bitfan/commons"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

type CodecCollection struct {
	Default *Codec
	Enc     *Codec
	Dec     *Codec
}

func (c *CodecCollection) NewEncoder(w io.Writer) (Encoder, error) {
	if c.Enc != nil {
		return c.Enc.NewEncoder(w)
	} else if c.Default != nil {
		return c.Default.NewEncoder(w)
	} else {
		return nil, fmt.Errorf("no decoder available")
	}
}

func (c *CodecCollection) NewDecoder(r io.Reader) (Decoder, error) {
	if c.Dec != nil {
		return c.Dec.NewDecoder(r)
	} else if c.Default != nil {
		return c.Default.NewDecoder(r)
	} else {
		return nil, fmt.Errorf("no decoder available")
	}

}

type Codec struct {
	Name                  string
	Role                  string
	Charset               string
	Options               map[string]interface{}
	logger                commons.Logger
	configWorkingLocation string
}

func New(name string, conf map[string]interface{}, logger commons.Logger, cwl string) *Codec {
	c := &Codec{
		Name:                  name,
		Charset:               "utf-8",
		Options:               map[string]interface{}{},
		logger:                logger,
		configWorkingLocation: cwl,
	}
	for i, k := range conf {
		if i == "charset" {
			c.Charset = k.(string)
			continue
		}
		if i == "role" {
			c.Role = k.(string)
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

	// Charset
	var cw io.Writer
	var err error
	var encoding encoding.Encoding
	encoding, err = ianaindex.IANA.Encoding(c.Charset)
	cw = encoding.NewEncoder().Writer(w)
	if err != nil {
		return enc, err
	}

	switch c.Name {
	case "pp":
		enc = rubydebugcodec.NewEncoder(cw)
	case "rubydebug":
		enc = rubydebugcodec.NewEncoder(cw)
	case "line":
		enc = linecodec.NewEncoder(cw)
	case "json":
		enc = jsoncodec.NewEncoder(cw)
	case "plain":
		enc = plaincodec.NewEncoder(cw)
	default:
		return enc, fmt.Errorf("no encoder defined")
	}

	err = enc.SetOptions(c.Options, c.logger, c.configWorkingLocation)

	return enc, err
}

func (c *Codec) NewDecoder(r io.Reader) (Decoder, error) {
	var dec Decoder

	if r == nil {
		return dec, fmt.Errorf("codecs.Codec.Decoder error : no reader !")
	}

	// Charset
	var cr io.Reader
	var err error
	var encoding encoding.Encoding
	encoding, err = ianaindex.IANA.Encoding(c.Charset)
	cr = encoding.NewDecoder().Reader(r)
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
	case "w3c":
		dec = w3ccodec.NewDecoder(cr)
	case "json":
		dec = jsoncodec.NewDecoder(cr)
	case "json_lines":
		dec = jsonlinescodec.NewDecoder(cr)
	default:
		dec = plaincodec.NewDecoder(cr)
	}
	dec.SetOptions(c.Options, c.logger, c.configWorkingLocation)

	return dec, nil
}
