//go:generate bitfanDoc
// Stop after emiting a blank event on start
// Allow you to put first event and then stop processors as soon as they finish their job.
//
// Permit to launch bitfan with a pipeline and quit when work is done.
package stopprocessor

import (
	"os"

	"github.com/vjeantet/bitfan/core"
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

	// Stop bitfan with the pipeline ending ?
	// @Default true
	ExitBitfan bool
}

type processor struct {
	processors.Base
	opt *options
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		ExitBitfan: true,
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func (p *processor) Receive(e processors.IPacket) error {
	processors.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
	p.Send(e)

	core.StopPipeline(p.PipelineID)

	if true == p.opt.ExitBitfan {
		core.Stop()
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(os.Interrupt)
	}

	return nil
}
func (p *processor) Start(e processors.IPacket) error {
	p.Logger.Debug("start with e=", e)
	go p.Receive(e)
	return nil
}
func (p *processor) Stop(e processors.IPacket) error {
	p.Logger.Debug("stopping pipeline ID=", p.PipelineID)

	return nil
}
