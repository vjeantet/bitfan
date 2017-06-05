package codecs

type Decoder interface {
	Decode() (map[string]interface{}, error)
	SetOptions(map[string]interface{}) error
	More() bool
}
