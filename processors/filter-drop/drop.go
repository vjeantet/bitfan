//go:generate bitfanDoc
// Drops everything received
// Drops everything that gets to this filter.
//
// This is best used in combination with conditionals, for example:
// ```
// filter {
//   if [loglevel] == "debug" {
//     drop { }
//   }
// }
// ```
// The above will only pass events to the drop filter if the loglevel field is debug. This will cause all events matching to be dropped.
package drop

import (
	"math/rand"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Drops all events
type processor struct {
	processors.Base

	opt *options
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

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt.Percentage = 100
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {

	if p.opt.Percentage == 100 || rand.Intn(100) < p.opt.Percentage {
		return nil
	}

	processors.ProcessCommonFields2(e.Fields(),
		p.opt.Add_field,
		p.opt.Add_tag,
		p.opt.Remove_field,
		p.opt.Remove_Tag,
	)
	p.Send(e, 0)
	return nil
}
