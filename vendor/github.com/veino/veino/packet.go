package veino

import "github.com/clbanning/mxj"

type IPacket interface {
	Kind() int
	Message() string
	Fields() *mxj.Map

	SetKind(int)
	SetMessage(string)
	SetFields(mxj.Map)

	Clone() IPacket
}

type PacketBuilder func(string, map[string]interface{}) IPacket

type PacketSender func(IPacket, ...int) bool
