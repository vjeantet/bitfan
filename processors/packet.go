package processors

import "github.com/clbanning/mxj"

type IPacket interface {
	Message() string
	Fields() *mxj.Map

	SetMessage(string)
	SetFields(map[string]interface{})

	Clone() IPacket
}

type PacketBuilder func(map[string]interface{}) IPacket

type PacketSender func(IPacket, ...int) bool
