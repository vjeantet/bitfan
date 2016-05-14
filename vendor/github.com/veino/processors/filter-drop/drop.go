// Drops everything received
package drop

import (
	"math/rand"

	"github.com/mitchellh/mapstructure"
	"github.com/veino/field"
	"github.com/veino/veino"
)

func New(l veino.Logger) veino.Processor {
	return &processor{}
}

type processor struct {
	Send veino.PacketSender
	opt  *options
}

type options struct {
	// If this event survice to drop, add any arbitrary fields to this event.
	// Field names can be dynamic and include parts of the event using the %{field}.
	Add_field map[string]interface{}

	// If this event survice to drop, add arbitrary tags to the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax.
	Add_tag []string

	// If this event survice to drop, remove arbitrary fields from this event.
	Remove_field []string

	// If this event survice to drop, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	Remove_Tag []string

	// Drop all the events within a pre-configured percentage.
	// This is useful if you just need a percentage but not the whole.
	Percentage int
}

func (p *processor) Configure(conf map[string]interface{}) error {
	cf := options{Percentage: 100}
	if mapstructure.Decode(conf, &cf) != nil {
		return nil
	}
	p.opt = &cf
	return nil
}

func (p *processor) Receive(e veino.IPacket) error {

	if p.opt.Percentage == 100 || rand.Intn(100) < p.opt.Percentage {
		return nil
	}

	field.ProcessCommonFields2(e.Fields(),
		p.opt.Add_field,
		p.opt.Add_tag,
		p.opt.Remove_field,
		p.opt.Remove_Tag,
	)
	p.Send(e, 0)
	return nil
}

func (p *processor) Tick(e veino.IPacket) error { return nil }

func (p *processor) Start(e veino.IPacket) error { return nil }

func (p *processor) Stop(e veino.IPacket) error { return nil }
