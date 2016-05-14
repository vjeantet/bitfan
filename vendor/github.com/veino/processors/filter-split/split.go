// splitt
package split

import (
	"github.com/mitchellh/mapstructure"
	"github.com/veino/field"
	"github.com/veino/veino"
)

const (
	PORT_SUCCESS = 0
	PORT_ERROR   = 1
)

func New(l veino.Logger) veino.Processor {
	return &processor{}
}

type processor struct {
	Send      veino.PacketSender
	NewPacket veino.PacketBuilder

	// The field which value is split by the terminator
	Field string
	// The field within the new event which the value is split into. If not set, target field defaults to split field name
	Target string
	// The string to split on. This is usually a line terminator, but can be any string
	// Default value is "\n"
	Terminator string

	// If this filter is successful, add any arbitrary fields to this event.
	// Field names can be dynamic and include parts of the event using the %{field}.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax.
	Add_tag []string

	// If this filter is successful, remove arbitrary fields from this event.
	Remove_field []string

	// If this filter is successful, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	Remove_Tag []string
}

func (p *processor) Configure(conf map[string]interface{}) error {
	if err := mapstructure.Decode(conf, p); err != nil {
		return err
	}
	return nil
}

func (p *processor) Receive(e veino.IPacket) error {
	// recupere les splits
	splits, _ := e.Fields().ValuesForPath(p.Field)
	// veino.Logger().Infof("err = %#v\nvalue=%#v\n\n", err, splits)

	if len(splits) == 0 {
		p.Send(e, PORT_ERROR)
		return nil
	}

	// iterate over found splits
	for _, split := range splits {
		// create a new event
		// set target value with split
		cp, _ := e.Fields().Copy()
		cp.SetValueForPath(split, p.Target)

		field.ProcessCommonFields2(&cp,
			p.Add_field,
			p.Add_tag,
			p.Remove_field,
			p.Remove_Tag,
		)

		// e := veino.NewEvent(e.ToAgentName(), e.Message(), cp)
		e := p.NewPacket(e.Message(), cp)
		p.Send(e, 0)
	}

	return nil
}

func (p *processor) Tick(e veino.IPacket) error { return nil }

func (p *processor) Start(e veino.IPacket) error { return nil }

func (p *processor) Stop(e veino.IPacket) error { return nil }
