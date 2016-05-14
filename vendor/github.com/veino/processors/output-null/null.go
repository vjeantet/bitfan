// Drops everything received
package null

import "github.com/veino/veino"

func New(l veino.Logger) veino.Processor {
	return &processor{}
}

type processor struct {
}

func (p *processor) Configure(conf map[string]interface{}) error { return nil }

func (p *processor) Receive(e veino.IPacket) error { return nil }

func (p *processor) Tick(e veino.IPacket) error { return nil }

func (p *processor) Start(e veino.IPacket) error { return nil }

func (p *processor) Stop(e veino.IPacket) error { return nil }
