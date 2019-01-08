//go:generate bitfanDoc
// The blacklist rule will check a certain field against a blacklist, and match if it is in the blacklist.
package blacklist

import (
	"fmt"

	"github.com/awillis/bitfan/processors"
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
	processors.CommonOptions `mapstructure:",squash"`

	// The name of the field to use to compare to the blacklist.
	// If the field is null, those events will be ignored.
	// @ExampleLS compare_field => "message"
	CompareField string `mapstructure:"compare_field" validate:"required"`

	// List of blacklisted terms.
	// The compare_field term must be equal to one of these values for it to match.
	// @ExampleLS terms => ["val1","val2","val3"]
	Terms []string `mapstructure:"terms"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{}
	p.opt = &defaults
	err = p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if len(p.opt.Terms) == 0 {
		return fmt.Errorf("blacklist option should have at least one value")
	}

	return err
}

func (p *processor) Receive(e processors.IPacket) error {
	for _, v := range p.opt.Terms {
		if v == e.Fields().ValueOrEmptyForPathString(p.opt.CompareField) {
			p.Logger.Debugf("blacklisted word %s found in %s", v, p.opt.CompareField)

			p.opt.ProcessCommonOptions(e.Fields())

			p.Send(e, 0)
			return nil
		}
	}
	p.Logger.Debugf("content of [%s] is not in the blacklist", p.opt.CompareField)

	return nil
}
