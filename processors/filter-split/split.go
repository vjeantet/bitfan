//go:generate bitfanDoc
// The split filter clones an event by splitting one of its fields and placing each value resulting from the split into a clone of the original event. The field being split can either be a string or an array.
//
// An example use case of this filter is for taking output from the exec input plugin which emits one event for the whole output of a command and splitting that output by newline - making each line an event.
//
// The end result of each split is a complete copy of the event with only the current split section of the given field changed.
package split

import "github.com/vjeantet/bitfan/processors"

const (
	PORT_SUCCESS = 0
	PORT_ERROR   = 1
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Splits multi-line messages into distinct events
type processor struct {
	processors.Base
	opt *options
}

type options struct {
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

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {

	// recupere les splits
	splits, _ := e.Fields().ValuesForPath(p.opt.Field)
	// processors.Logger().Infof("err = %#v\nvalue=%#v\n\n", err, splits)

	if len(splits) == 0 {
		p.Send(e, PORT_ERROR)
		return nil
	}

	// iterate over found splits
	for _, split := range splits {
		// create a new event
		// set target value with split
		cp, _ := e.Fields().Copy()
		cp.SetValueForPath(split, p.opt.Target)
		processors.ProcessCommonFields2(&cp,
			p.opt.Add_field,
			p.opt.Add_tag,
			p.opt.Remove_field,
			p.opt.Remove_Tag,
		)

		// e := processors.NewEvent(e.ToAgentName(), e.Message(), cp)
		e2 := p.NewPacket(e.Message(), cp)
		p.Send(e2, 0)
	}

	return nil
}
