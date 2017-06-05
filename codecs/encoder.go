package codecs

type Encoder interface {
	Encode(map[string]interface{}) error
	SetOptions(map[string]interface{}) error
}
