package processors

import "github.com/vjeantet/bitfan/processors/doc"

type Processor interface {
	Configure(ProcessorContext, map[string]interface{}) error
	Start(IPacket) error
	Tick(IPacket) error
	Receive(IPacket) error
	Stop(IPacket) error
	Doc() *doc.Processor
	MaxConcurent() int
}

type ProcessorContext interface {
	Log() Logger
	PacketSender() PacketSender
	PacketBuilder() PacketBuilder
	Memory() Memory
	ConfigWorkingLocation() string
	DataLocation() string
}
