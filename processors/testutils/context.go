package testutils

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/vjeantet/bitfan/processors"
)

type DummyProcessorContext struct {
	logger        processors.Logger
	packetSender  processors.PacketSender
	packetBuilder processors.PacketBuilder
	sentPackets   map[int][]processors.IPacket
	builtPackets  []processors.IPacket
	memory        processors.Memory
	store         processors.IStore
	mock.Mock
}

func NewProcessorContext() *DummyProcessorContext {
	dp := &DummyProcessorContext{}
	// dp.logger = log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
	dp.logger = logrus.New()
	dp.sentPackets = map[int][]processors.IPacket{}
	dp.packetSender = newSender(dp)
	dp.packetBuilder = newPacket(dp)
	dp.memory = newMemory(dp)
	dp.store = newStore(dp)
	return dp
}

func newSender(p *DummyProcessorContext) processors.PacketSender {
	return func(packet processors.IPacket, portNumbers ...int) bool {
		if len(portNumbers) == 0 {
			portNumbers = []int{0}
		}
		for _, portNumber := range portNumbers {
			p.sentPackets[portNumber] = append(p.sentPackets[portNumber], packet)
		}
		return true
	}
}

func newPacket(p *DummyProcessorContext) processors.PacketBuilder {
	return func(fields map[string]interface{}) processors.IPacket {
		e := NewPacket(fields)
		p.builtPackets = append(p.builtPackets, e)
		return e
	}
}

func (p DummyProcessorContext) Log() processors.Logger {
	return p.logger
}
func (p DummyProcessorContext) PacketSender() processors.PacketSender {
	return p.packetSender
}
func (p DummyProcessorContext) PacketBuilder() processors.PacketBuilder {
	return p.packetBuilder
}
func (d *DummyProcessorContext) SentPackets(portNumber int) []processors.IPacket {
	return d.sentPackets[portNumber]
}
func (d *DummyProcessorContext) BuiltPackets() []processors.IPacket {
	return d.builtPackets
}
func (d *DummyProcessorContext) BuiltPacketsCount() int {
	return len(d.builtPackets)
}
func (d *DummyProcessorContext) SentPacketsCount(portNumber int) int {
	return len(d.SentPackets(portNumber))
}

func (d *DummyProcessorContext) ConfigWorkingLocation() string {
	return ""
}

func (d *DummyProcessorContext) DataLocation() string {
	return ""
}

func (d *DummyProcessorContext) Memory() processors.Memory {
	return d.memory
}

func (d *DummyProcessorContext) Store() processors.IStore {
	return d.store
}

var i = 0

func newMemory(p *DummyProcessorContext) processors.Memory {
	i += 1
	return NewMemory("").Space(fmt.Sprintf("test_%d", i))
}

func (d *DummyProcessorContext) WebHook() processors.WebHook {
	return nil
}
