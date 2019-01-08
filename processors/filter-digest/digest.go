//go:generate bitfanDoc
package digest

import (
	"fmt"
	"regexp"

	"github.com/awillis/bitfan/processors"
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
	packets int
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

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
	if p.opt.Count < 0 {
		return fmt.Errorf("Negative count setting")
	}
	//
	//if p.opt.Interval != "" {
	//	c := cron.New()
	//	c.AddFunc().Add(a.Label, a.Schedule, func() {
	//		go a.processor.Tick(newPacket("", nil))
	//	})
	//}
	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	p.packets++
	if p.opt.KeyMap == "" {
		// No key map: merge the event fields with the current data
		for k, v := range *e.Fields() {
			p.values[k] = v
		}
	} else {
		k, err := e.Fields().ValueForPathString(p.opt.KeyMap)
		if err != nil {
			p.Logger.Errorf("can not find value for key %s", p.opt.KeyMap)
			return err
		}
		p.values[k] = e.Fields().Old()
	}

	// When no interval, flush event when 'Count' events are digested
	if p.opt.Interval == "" {
		if p.packets >= p.opt.Count {
			p.Logger.Debugf("Flush digester ! %d/%d events digested", p.packets, p.opt.Count)
			return p.Tick(e)
		}
	}

	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	// When Interval is set, and total digested events < Count : ignore
	if p.opt.Interval != "" {
		if p.packets < p.opt.Count {
			p.Logger.Errorf("Ignore tick interval, %d/%s events digested", p.packets, p.opt.Count)
			return nil
		}
	}

	ne := p.NewPacket(p.values)
	p.opt.ProcessCommonOptions(ne.Fields())
	p.Send(ne, PORT_SUCCESS)
	p.values = map[string]interface{}{}
	p.packets = 0
	return nil
}
