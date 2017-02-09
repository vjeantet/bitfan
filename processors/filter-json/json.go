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
	// If this filter is successful, add any arbitrary fields to this event.
	// Field names can be dynamic and include parts of the event using the %{field}.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax.
	Add_tag []string

	// If this filter is successful, remove arbitrary fields from this event.
	Remove_field []string

	// If this filter is successful, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	Remove_tag []string

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

	processors.ProcessCommonFields2(e.Fields(),
		p.opt.Add_field,
		p.opt.Add_tag,
		p.opt.Remove_field,
		p.opt.Remove_tag,
	)

	p.Send(e, 0)
	return nil
}

func (p *processor) Tick(e processors.IPacket) error { return nil }

func (p *processor) Start(e processors.IPacket) error { return nil }

func (p *processor) Stop(e processors.IPacket) error { return nil }
