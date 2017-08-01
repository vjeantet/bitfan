//go:generate bitfanDoc
// Sleep a given amount of time.
//
// This will cause bitfan to stall for the given amount of time.
//
// This is useful for rate limiting, etc.
package sleepprocessor

import (
	"time"

	"github.com/vjeantet/bitfan/processors"
)

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
	processors.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
	p.Send(e)

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	return nil
}
