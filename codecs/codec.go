// TODO : in each encoder, initialize encoder on configure instead of receive
// TODO : implement logger and cwl in decoders
package codecs

import (
	"fmt"
	"io"

	"github.com/vjeantet/bitfan/codecs/csv"
	"github.com/vjeantet/bitfan/codecs/json"
	"github.com/vjeantet/bitfan/codecs/jsonlines"
	"github.com/vjeantet/bitfan/codecs/lib"
	"github.com/vjeantet/bitfan/codecs/line"
	"github.com/vjeantet/bitfan/codecs/multiline"
	"github.com/vjeantet/bitfan/codecs/plain"
	"github.com/vjeantet/bitfan/codecs/rubydebug"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

type Codec struct {
	Name                  string
	Charset               string
	Options               map[string]interface{}
	logger                lib.Logger
	configWorkingLocation string
}

func New(name string, conf map[string]interface{}, logger lib.Logger, cwl string) Codec {
	c := Codec{
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
	default:
		return enc, fmt.Errorf("no encoder defined")
	}

	if err := enc.SetOptions(c.Options, c.logger, c.configWorkingLocation); err != nil {
		return enc, err
	}

	return enc, nil
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
