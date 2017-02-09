//go:generate bitfanDoc
package use

import "github.com/vjeantet/bitfan/processors"

const (
	PORT_SUCCESS = 0
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// If this processor is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// If this processor is successful, add arbitrary tags to the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax.
	Add_tag []string

	// If this processor is successful, remove arbitrary fields from this event.
	Remove_field []string

	// If this processor is successful, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	Remove_tag []string

	// Add a type field to all events handled by this processor
	Type string

	// Path to configuration to import in this pipeline, it could be a local file or an url
	// can be relative path to the current configuration.
	// SPLIT and JOIN : in filter Section, set multiples path to make a split and join into your pipeline
	// @ExampleLS path=> ["meteo-input.conf"]
	Path []string `mapstructure:"path" validate:"required"`

	// You can set variable references in the used configuration by using ${var}.
	// each reference will be replaced by the value of the variable found in this option
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`
}

// Include a config file
type processor struct {
	processors.Base

	opt *options
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt = &options{}
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	processors.ProcessCommonFields2(e.Fields(),
		p.opt.Add_field,
		p.opt.Add_tag,
		p.opt.Remove_field,
		p.opt.Remove_tag,
	)

	p.Send(e, PORT_SUCCESS)
	return nil
}
