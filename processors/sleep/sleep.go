//go:generate bitfanDoc
// Sleep a given amount of time.
//
// This will cause bitfan to stall for the given amount of time.
//
// This is useful for rate limiting, etc.
package sleepprocessor

import (
	"time"

	"bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The length of time to sleep, in Millisecond, for every event.
	Time int
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

func (p *processor) Receive(e processors.IPacket) error {
	time.Sleep(time.Millisecond * time.Duration(p.opt.Time))
	p.opt.ProcessCommonOptions(e.Fields())
	p.Send(e)

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	return nil
}
