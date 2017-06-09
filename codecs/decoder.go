package codecs

import "github.com/vjeantet/bitfan/codecs/lib"

type Decoder interface {
	Decode() (map[string]interface{}, error)
	SetOptions(map[string]interface{}, lib.Logger, string) error
	More() bool
}
