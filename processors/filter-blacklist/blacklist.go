//go:generate bitfanDoc
// The blacklist rule will check a certain field against a blacklist, and match if it is in the blacklist.
package blacklist

import (
	"fmt"

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

// drop event when term not in a given list
type processor struct {
	processors.Base
	opt *options
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

	// The name of the field to use to compare to the blacklist.
	// If the field is null, those events will be ignored.
	// @ExampleLS compare_field => "message"
	CompareField string `mapstructure:"compare_field" validate:"required"`

	// A list of blacklisted values.
	// The compare_field term must be equal to one of these values for it to match.
	// @ExampleLS list => ["val1","val2","val3"]
	List []string `mapstructure:"list" validate:"required"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{}
	p.opt = &defaults
	err = p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if len(p.opt.List) == 0 {
		return fmt.Errorf("blacklist option should have at least one value")
	}

	return err
}

func (p *processor) Receive(e processors.IPacket) error {
	for _, v := range p.opt.List {
		if v == e.Fields().ValueOrEmptyForPathString(p.opt.CompareField) {
			p.Logger.Debugf("blacklisted word %s found in %s", v, p.opt.CompareField)

			processors.ProcessCommonFields2(e.Fields(),
				p.opt.AddField,
				p.opt.AddTag,
				p.opt.RemoveField,
				p.opt.RemoveTag,
			)
			p.Send(e, 0)
			return nil
		}
	}
	p.Logger.Debugf("content of [%s] is not in the blacklist", p.opt.CompareField)

	return nil
}
