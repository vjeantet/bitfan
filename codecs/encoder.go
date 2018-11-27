package codecs

import "bitfan/commons"

type Encoder interface {
	Encode(map[string]interface{}) error
	SetOptions(map[string]interface{}, commons.Logger, string) error
}
