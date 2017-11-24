package core

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

func (c *Codec) GetName() string {
	return c.Name
}

func (c *Codec) GetRole() string {
	return c.Role
}

func (c *Codec) GetOptions() map[string]interface{} {
	return c.Options
}
