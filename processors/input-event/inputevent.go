//go:generate bitfanDoc
// Generate a blank event on interval
package inputeventprocessor

import (
	"sync"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// string value to put in event
	Message string

	// How many events to generate
	// @Default 1
	// @ExampleLS count => 1000000
	Count int `mapstructure:"count"`

	// Use CRON or BITFAN notation
	// When omited, event will be generated on start
	// @ExampleLS interval => "@every 10s"
	Interval string `mapstructure:"interval"`
}

type processor struct {
	processors.Base
	opt *options
	wg  sync.WaitGroup
}

func (p *processor) MaxConcurent() int { return 0 }
func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Count: 1,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {
	if p.opt.Interval == "" {
		go p.Tick(e)
	}
	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	p.wg.Add(1)
	for i := 1; i <= p.opt.Count; i++ {
		e := p.NewPacket(
			p.opt.Message,
			map[string]interface{}{"number": i},
		)
		p.opt.ProcessCommonOptions(e.Fields())
		p.Send(e)
	}
	p.wg.Done()
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.Logger.Debugf("finishing event generation...")
	p.wg.Wait()
	return nil
}
