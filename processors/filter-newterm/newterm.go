//go:generate bitfanDoc
// This processor matches when a new value appears in a field that has never been seen before.
package newterm

import (
	"sync"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

const (
	PORT_SUCCESS = 0
)

// no concurency limit
func (p *processor) MaxConcurent() int { return 0 }

// drop event when term was already seen before
type processor struct {
	processors.Base
	opt *options

	mu    sync.RWMutex
	terms []string
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

	// The name of the field to use to compare to terms list.
	// If the field is null, those events will be ignored.
	// @ExampleLS compare_field => "message"
	CompareField string `mapstructure:"compare_field" validate:"required"`

	// If true, events without a compare_field field will be ignored.
	// @ExampleLS ignore_missing => true
	// @Default true
	IgnoreMissing bool `mapstructure:"ignore_missing"`

	// A list of initial terms to consider now new.
	// The compare_field term must be in this list or else it will match.
	// @ExampleLS terms => ["val1","val2","val3"]
	Terms []string `mapstructure:"terms"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		IgnoreMissing: true,
	}
	p.opt = &defaults
	err = p.ConfigureAndValidate(ctx, conf, p.opt)
	p.terms = p.opt.Terms
	return err
}

func (p *processor) Receive(e processors.IPacket) error {
	eValue, err := e.Fields().ValueForPathString(p.opt.CompareField)
	if err != nil { // path not found
		if p.opt.IgnoreMissing == true {
			return nil
		}
		p.Logger.Debugf("missing field [%s]", p.opt.CompareField)
	} else {
		p.mu.RLock()
		for _, v := range p.terms {
			if v == eValue {
				p.Logger.Debugf("ignore event, term '%s' already seen", eValue)
				p.mu.RUnlock()
				return nil
			}
		}
		p.mu.RUnlock()

		p.Logger.Debugf("new term '%s' found in [%s]", eValue, p.opt.CompareField)

		p.mu.Lock()
		p.terms = append(p.terms, eValue)
		p.mu.Unlock()
	}

	processors.ProcessCommonFields2(e.Fields(),
		p.opt.AddField,
		p.opt.AddTag,
		p.opt.RemoveField,
		p.opt.RemoveTag,
	)
	p.Send(e, 0)

	return nil
}
