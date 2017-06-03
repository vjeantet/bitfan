package codec

import (
	"fmt"
	"io"
)

var EOF = fmt.Errorf("End of file")
var NOMOREDATA = fmt.Errorf("No more data")

type Decoder interface {
	Decode() (map[string]interface{}, error)
	More() bool
	DecodeReader(io.Reader) (map[string]interface{}, error)
}
