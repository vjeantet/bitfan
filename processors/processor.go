package processors

import "github.com/awillis/bitfan/processors/doc"

type Processor interface {
	B() *Base
	Configure(ProcessorContext, map[string]interface{}) error
	Start(IPacket) error
	Tick(IPacket) error
	Receive(IPacket) error
	Stop(IPacket) error
	Doc() *doc.Processor
	MaxConcurent() int
	SetPipelineUUID(string)
	SetProcessorIdentifiers(string, string)
}

type ProcessorContext interface {
	Log() Logger
	PacketSender() PacketSender
	PacketBuilder() PacketBuilder
	Memory() Memory
	WebHook() WebHook
	ConfigWorkingLocation() string
	DataLocation() string
	Store() IStore
}
