package codecs

import "io"

type Decoder interface {
	Decode() (map[string]interface{}, error)
	More() bool
	DecodeReader(io.Reader) (map[string]interface{}, error)
}
