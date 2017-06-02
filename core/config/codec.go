package config

type Codec struct {
	Name    string
	Options map[string]interface{}
}

func NewCodec(name string) *Codec {
	return &Codec{
		Name:    name,
		Options: map[string]interface{}{},
	}
}

func (c *Codec) String() string {
	return c.Name
}
