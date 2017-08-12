//go:generate bitfanDoc
package json

import (
	"encoding/json"

	"github.com/vjeantet/bitfan/processors"
)

// Parses JSON events
func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Parses JSON events
type processor struct {
	processors.Base

	opt *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The configuration for the JSON filter
	Source string

	// Define the target field for placing the parsed data. If this setting is omitted,
	// the JSON data will be stored at the root (top level) of the event
	Target string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {

	json_string, err := e.Fields().ValueForPathString(p.opt.Source)
	if err != nil {
		return err
	}

	byt := []byte(json_string)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		return err
	}

	if p.opt.Target != "" {
		e.Fields().SetValueForPath(dat, p.opt.Target)
	} else {
		for k, v := range dat {
			e.Fields().SetValueForPath(v, k)
		}
	}

	p.opt.ProcessCommonOptions(e.Fields())

	p.Send(e, 0)
	return nil
}

func (p *processor) Tick(e processors.IPacket) error { return nil }

func (p *processor) Start(e processors.IPacket) error { return nil }

func (p *processor) Stop(e processors.IPacket) error { return nil }
