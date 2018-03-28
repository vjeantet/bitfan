package webfan

import (
	"fmt"

	"github.com/vjeantet/bitfan/processors"
)

func NewPass() processors.Processor {
	return &passProcessor{}
}

type passOptions struct {
}

// Prints events to the standard output
type passProcessor struct {
	processors.Base

	events chan processors.IPacket
}

func (p *passProcessor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {

	if err := p.ConfigureAndValidate(ctx, conf, &passOptions{}); err != nil {
		return err
	}

	if w, ok := conf["chan"]; ok {
		p.events = w.(chan processors.IPacket)
	} else {
		return fmt.Errorf("no chan")
	}

	return nil
}

func (p *passProcessor) Receive(e processors.IPacket) error {
	p.events <- e
	return nil
}

func (p *passProcessor) Stop(e processors.IPacket) error {
	close(p.events)
	return nil
}
