//go:generate bitfanDoc
// Generate a blank event on interval
package inputeventprocessor

import (
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// string value to put in event
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
	p.opt.ProcessCommonOptions(e.Fields())
	p.Send(e)
	return nil
}
