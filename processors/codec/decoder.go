package codec

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/net/html/charset"
)

var EOF = fmt.Errorf("End of file")
var NOMOREDATA = fmt.Errorf("No more data")

type Decoder interface {
	Decode() (map[string]interface{}, error)
	More() bool
	DecodeReader(io.Reader) (map[string]interface{}, error)
}

func NewDecoder(name string) (Decoder, error) {
	var dec Decoder

	//todo get Charset from Codec settings
	f := bytes.NewReader(nil)
	cr, err := charset.NewReaderLabel("utf8", f)
	if err != nil {
		return nil, err
	}

	switch name {
	case "json":
		dec = NewJsonDecoder(cr)
	case "csv":
		dec = NewCsvDecoder(cr)
	default:
		dec = NewPlainDecoder(cr)
	}
	return dec, nil
}
