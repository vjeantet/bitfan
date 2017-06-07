package codecs

import "github.com/vjeantet/bitfan/codecs/lib"

type Encoder interface {
	Encode(map[string]interface{}) error
	SetOptions(map[string]interface{}, lib.Logger, string) error
}
