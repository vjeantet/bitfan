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

	"bitfan/processors"
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
	processors.CommonOptions `mapstructure:",squash"`

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

	p.opt.ProcessCommonOptions(e.Fields())
	p.Send(e, 0)
	return nil
}
