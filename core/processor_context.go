package core

import "github.com/awillis/bitfan/processors"

type processorContext struct {
	packetSender          processors.PacketSender
	packetBuilder         processors.PacketBuilder
	logger                processors.Logger
	memory                processors.Memory
	webHook               processors.WebHook
	store                 processors.IStore
	dataLocation          string
	configWorkingLocation string
}

func (p processorContext) Log() processors.Logger {
	return p.logger
}
func (p processorContext) Memory() processors.Memory {
	return p.memory
}

func (p processorContext) WebHook() processors.WebHook {
	return p.webHook
}
func (p processorContext) PacketSender() processors.PacketSender {
	return p.packetSender
}
func (p processorContext) PacketBuilder() processors.PacketBuilder {
	return p.packetBuilder
}
func (p processorContext) ConfigWorkingLocation() string {
	return p.configWorkingLocation
}

func (p processorContext) DataLocation() string {
	return p.dataLocation
}

func (p processorContext) Store() processors.IStore {
	return p.store
}
