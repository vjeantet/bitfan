package codecs

type Encoder interface {
	Encode(map[string]interface{}) error
}
