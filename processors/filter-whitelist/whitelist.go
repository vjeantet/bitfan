//go:generate bitfanDoc
// Similar to blacklist, this processor will compare a certain field to a whitelist, and match
// if the list does not contain the term
package whitelist

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

// drop event when term is in a given list
type processor struct {
	processors.Base
	opt *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The name of the field to use to compare to the whitelist.
	// If the field is null, those events will be ignored.
	// @ExampleLS compare_field => "message"
	CompareField string `mapstructure:"compare_field" validate:"required"`

	// If true, events without a compare_key field will not match.
	// @Default true
	IgnoreMissing bool `mapstructure:"ignore_missing"`

	// A list of whitelisted terms.
	// The compare_field term must be in this list or else it will match.
	// @ExampleLS terms => ["val1","val2","val3"]
	Terms []string `mapstructure:"terms" validate:"required"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		IgnoreMissing: true,
	}
	p.opt = &defaults
	err = p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if len(p.opt.Terms) == 0 {
		return fmt.Errorf("whitelist option should have at least one value")
	}

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
		for _, v := range p.opt.Terms {
			if v == eValue {
				p.Logger.Debugf("white word %s found in %s", v, p.opt.CompareField)
				return nil
			}
		}
		p.Logger.Debugf("content of [%s] is not in the whitelist", p.opt.CompareField)
	}

	p.opt.ProcessCommonOptions(e.Fields())
	p.Send(e, 0)

	return nil
}
