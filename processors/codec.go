package processors

type ICodec interface {
	GetName() string
	GetOptions() map[string]interface{}
	GetRole() string
}
