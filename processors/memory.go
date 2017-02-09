package processors

type Memory interface {
	Set(name string, value interface{})
	Get(name string) (interface{}, bool)
	Delete(name string)
	Items() map[string]interface{}
}
