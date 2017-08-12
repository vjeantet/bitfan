//go:generate bitfanDoc
package gennumbers

import "github.com/vjeantet/bitfan/processors"

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// How many events to generate
	// @Default 1000000
	// @ExampleLS count => 1000000
	Count int `mapstructure:"count"`

	// @ExampleLS interval => "10"
	// @Type Interval
	Interval string `mapstructure:"interval" `
}

// generate a number of event
type processor struct {
	processors.Base

	opt *options
	q   chan bool
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Count: 1000000,
	}

	p.opt = &defaults

	var err error
	err = p.ConfigureAndValidate(ctx, conf, p.opt)

	return err
}

// 0 = no limit !
func (p *processor) MaxConcurent() int { return 0 }

func (p *processor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func (p *processor) Receive(e processors.IPacket) error {
	p.opt.ProcessCommonOptions(e.Fields())
	p.Send(e)
	return nil
}

func (p *processor) Start(e processors.IPacket) error {
	p.q = make(chan bool)

	go func() {
		for i := 1; i <= p.opt.Count; i++ {
			select {
			case <-p.q:
				close(p.q)
				return
			default:
				e := p.NewPacket(
					"",
					map[string]interface{}{"number": i},
				)
				p.opt.ProcessCommonOptions(e.Fields())
				p.Send(e)
			}
		}
		<-p.q
		close(p.q)
	}()

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.q <- true
	<-p.q
	return nil
}
