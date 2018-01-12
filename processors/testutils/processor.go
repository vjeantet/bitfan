package testutils

import (
	"github.com/stretchr/testify/mock"
	"github.com/vjeantet/bitfan/processors"
)

type Processor struct {
	processors.Processor
	mock.Mock

	ctx *DummyProcessorContext
}

func StartNewProcessor(f func() processors.Processor, conf ...map[string]interface{}) (Processor, error) {
	p, err := NewProcessor(f, conf...)
	if err != nil {
		return p, err
	}
	err = p.Start(nil)
	if err != nil {
		return p, err
	}

	return p, nil

}
func NewProcessor(f func() processors.Processor, conf ...map[string]interface{}) (Processor, error) {
	var err error
	p := newMockedProcessor(f)
	p.ctx = NewProcessorContext()
	if len(conf) > 0 {
		err = p.Configure(p.ctx, conf[0])
	}
	return p, err
}

func newMockedProcessor(f func() processors.Processor) Processor {
	return Processor{Processor: f()}
}

func (p *Processor) SentPacketsCount(portNumber int) int {
	return p.ctx.SentPacketsCount(portNumber)
}
func (p *Processor) SentPackets(portNumber int) []processors.IPacket {
	return p.ctx.SentPackets(portNumber)
}
func (p *Processor) BuiltPackets() []processors.IPacket {
	return p.ctx.BuiltPackets()
}
func (p *Processor) BuiltPacketsCount() int {
	return p.ctx.BuiltPacketsCount()
}
