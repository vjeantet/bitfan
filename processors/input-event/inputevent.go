//go:generate bitfanDoc
// Generate a blank event on interval
package inputeventprocessor

import "github.com/vjeantet/bitfan/processors"

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	Tags []string

	// Add a type field to all events handled by this input
	Type string

	Message string

	// Use CRON or BITFAN notation
	// @ExampleLS interval => "@every 10s"
	Interval string `mapstructure:"interval" validate:"required"`
}

type processor struct {
	processors.Base
	opt *options
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Tick(e processors.IPacket) error {
	processors.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
	p.Send(e)
	return nil
}
