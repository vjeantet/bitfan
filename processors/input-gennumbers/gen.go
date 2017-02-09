//go:generate bitfanDoc
package gennumbers

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

	// How many events to generate
	// @Default 1000000
	// @ExampleLS count => 1000000
	Count int `mapstructure:"count"`
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
				processors.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
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
