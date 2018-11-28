//go:generate bitfanDoc
// The uuid filter allows you to generate a UUID and add it as a field to each processed event.
//
// This is useful if you need to generate a string that’s unique for every event, even if the same input is processed multiple times. If you want to generate strings that are identical each time a event with a given content is processed (i.e. a hash) you should use the fingerprint filter instead.
//
// The generated UUIDs follow the version 4 definition in RFC 4122).
package uuid

import (
	"bitfan/processors"
	"github.com/nu7hatch/gouuid"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Adds a UUID to events
type processor struct {
	processors.Base

	opt *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// If the value in the field currently (if any) should be overridden by the generated UUID.
	// Defaults to false (i.e. if the field is present, with ANY value, it won’t be overridden)
	Overwrite bool

	// Add a UUID to a field
	Target string `mapstructure:"target" validate:"required"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt.Overwrite = false
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	id, err := uuid.NewV4()

	if err == nil {
		if !(p.opt.Overwrite == false && e.Fields().Exists(p.opt.Target) == true) {
			e.Fields().SetValueForPath(id.String(), p.opt.Target)
		}

		p.opt.ProcessCommonOptions(e.Fields())

	}

	p.Send(e, 0)
	return nil
}
