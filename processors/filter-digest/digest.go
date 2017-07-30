//go:generate bitfanDoc
package digest

import (
	"fmt"
	"regexp"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

const (
	PORT_SUCCESS = 0
)

func (p *processor) MaxConcurent() int { return 1 }

// Digest events every x
type processor struct {
	processors.Base
	opt     *options
	scan_re *regexp.Regexp
	values  map[string]interface{}
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	AddField map[string]interface{} `mapstructure:"add_field"`

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	AddTag []string `mapstructure:"add_tag"`

	// If this filter is successful, remove arbitrary fields from this event. Example:
	// ` kv {
	// `   remove_field => [ "foo_%{somefield}" ]
	// ` }
	RemoveField []string `mapstructure:"remove_field"`

	// If this filter is successful, remove arbitrary tags from the event. Tags can be dynamic and include parts of the event using the %{field} syntax.
	// Example:
	// ` kv {
	// `   remove_tag => [ "foo_%{somefield}" ]
	// ` }
	// If the event has field "somefield" == "hello" this filter, on success, would remove the tag foo_hello if it is present. The second example would remove a sad, unwanted tag as well.
	RemoveTag []string `mapstructure:"remove_tag"`

	// Add received event fields to the digest field named with the key map_key
	// When this setting is empty, digest will merge fields from coming events
	// @ExampleLS key_map => "type"
	KeyMap string `mapstructure:"key_map"`

	// When should Digest send a digested event ?
	// Use CRON or BITFAN notation
	// @ExampleLS interval => "every_10s"
	Interval string `mapstructure:"interval"`

	// With min > 0, digest will not fire an event if less than min events were digested
	Count int `mapstructure:"count"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		Count: 0,
	}
	p.opt = &defaults
	if err = p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}
	p.values = map[string]interface{}{}

	if p.opt.Interval == "" && p.opt.Count == 0 {
		return fmt.Errorf("no interval and no Count settings set")
	}

	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	if p.opt.KeyMap == "" {
		for k, v := range *e.Fields() {
			p.values[k] = v
		}
	} else {
		k, err := e.Fields().ValueForPathString(p.opt.KeyMap)
		if err != nil {
			p.Logger.Errorf("can not find value for key", p.opt.KeyMap)
			return err
		}
		p.values[k] = e.Fields()
	}

	// When no interval, flush event when Count events where digested
	if p.opt.Interval == "" {
		if len(p.values) >= p.opt.Count {
			p.Logger.Debugf("Flush digester ! %d/%d events digested", len(p.values), p.opt.Count)
			return p.Tick(e)
		}
	}

	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	// When Interval is set, and total digested events < Count : ignore
	if p.opt.Interval != "" {
		if len(p.values) < p.opt.Count {
			p.Logger.Debugf("Ignore tick interval, %d/%s events digested", len(p.values), p.opt.Count)
			return nil
		}
	}

	ne := p.NewPacket("", p.values)
	processors.ProcessCommonFields2(ne.Fields(),
		p.opt.AddField,
		p.opt.AddTag,
		p.opt.RemoveField,
		p.opt.RemoveTag,
	)
	p.Send(ne, PORT_SUCCESS)
	p.values = map[string]interface{}{}
	return nil
}
