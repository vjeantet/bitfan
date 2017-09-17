package config

type Codec struct {
	Name    string
	Role    string //decoder/encoder
	Options map[string]interface{}
}

func NewCodec(name string) *Codec {
	return &Codec{
		Name:    name,
		Role:    "",
		Options: map[string]interface{}{},
	}
}

func (c *Codec) String() string {
	return c.Name
}
