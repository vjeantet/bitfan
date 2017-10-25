//go:generate bitfanDoc
package templateprocessor

import (
	"bytes"
	"text/template"

	"github.com/vjeantet/bitfan/core/location"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// Go Template content
	//
	// set inline content, a path or an url to the template content
	//
	// Go template : https://golang.org/pkg/html/template/
	// @ExampleLS location => "test.tpl"
	// @Type Location
	Location string `mapstructure:"location" validate:"required"`

	// You can set variable to be used in template by using ${var}.
	// each reference will be replaced by the value of the variable found in Template's content
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`

	// Define the target field for placing the template execution result. If this setting is omitted,
	// the data will be stored in the "output" field
	// @ExampleLS target => "mydata"
	// @Default "output"
	Target string `mapstructure:"target"`
}

type processor struct {
	processors.Base
	opt *options
	q   chan bool
	Tpl *template.Template
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Target: "output",
	}

	p.opt = &defaults

	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	loc, err := location.NewLocation(p.opt.Location, p.ConfigWorkingLocation)
	if err != nil {
		return err
	}

	tpl, _, err := loc.TemplateWithOptions(p.opt.Var)
	if err != nil {
		return err
	}

	p.Tpl = tpl
	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func (p *processor) Receive(e processors.IPacket) error {

	buff := bytes.NewBufferString("")
	err := p.Tpl.Execute(buff, e.Fields())
	if err != nil {
		p.Logger.Errorf("template error : %v", err)
		return err
	}

	if len(p.opt.Var) > 0 {
		e.Fields().SetValueForPath(p.opt.Var, "var")
	}
	e.Fields().SetValueForPath(buff.String(), p.opt.Target)

	p.opt.ProcessCommonOptions(e.Fields())
	p.Send(e)

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	return nil
}
