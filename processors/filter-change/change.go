//go:generate bitfanDoc
// This rule will monitor a certain field and match if that field changes. The field must change with respect to the last event
package change

import (
	"sync"
	"time"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

const (
	PORT_SUCCESS = 0
)

// no concurency ! only one worker
func (p *processor) MaxConcurent() int { return 1 }

// drop event when field value is the same in the last event
type processor struct {
	processors.Base
	opt   *options
	first bool

	mu        sync.Mutex
	lastValue string
	hop       *time.Timer
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	AddField map[string]interface{} `mapstructure:"add_field"`

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	AddTag []string `mapstructure:"add_tag"`

	// If this filter is successful, remove arbitrary fields from this event. Example:
	// ```
	// kv {
	//   remove_field => [ "foo_%{somefield}" ]
	// }
	// ```
	RemoveField []string `mapstructure:"remove_field"`

	// If this filter is successful, remove arbitrary tags from the event. Tags can be dynamic and include parts of the event using the %{field} syntax.
	//
	// Example:
	// ```
	// kv {
	//   remove_tag => [ "foo_%{somefield}" ]
	// }
	// ```
	// If the event has field "somefield" == "hello" this filter, on success, would remove the tag foo_hello if it is present. The second example would remove a sad, unwanted tag as well.
	RemoveTag []string `mapstructure:"remove_tag"`

	// The name of the field to use to compare to the blacklist.
	//
	// If the field is null, those events will be ignored.
	// @ExampleLS compare_field => "message"
	CompareField string `mapstructure:"compare_field" validate:"required"`

	// If true, events without a compare_key field will not count as changed.
	// @Default true
	IgnoreMissing bool `mapstructure:"ignore_missing"`

	// The maximum time in seconds between changes. After this time period, Bitfan will forget the old value of the compare_field field.
	// @Default 0 (no timeframe)
	// @ExampleLS timeframe => 10
	Timeframe int `mapstructure:"timeframe"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		IgnoreMissing: true,
		Timeframe:     0,
	}
	p.opt = &defaults
	p.first = true
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	eValue, err := e.Fields().ValueForPathString(p.opt.CompareField)
	if err != nil { // path not found
		if p.opt.IgnoreMissing == true {
			return nil
		}
		p.Logger.Debugf("missing field [%s]", p.opt.CompareField)
	} else {
		p.mu.Lock()

		if p.lastValue == eValue {
			p.mu.Unlock()
			return nil
		}

		p.Logger.Debugf("[%s] value change from '%s' to '%s'", p.opt.CompareField, p.lastValue, eValue)
		p.lastValue = eValue

		if p.first == true {
			p.Logger.Debugf("ignore first change on field [%s]", p.opt.CompareField)
			p.first = false
			p.mu.Unlock()
			return nil
		}

		// Change occured !

		if p.opt.Timeframe > 0 {
			if p.hop == nil { // Initiate timer
				p.Logger.Debugf("Timer inited")
				p.hop = time.AfterFunc(time.Second*time.Duration(p.opt.Timeframe), func() { // when timeframe expires, reset old value
					p.Logger.Debugf("expired !")
					p.mu.Lock()
					p.lastValue = ""
					p.hop = nil
					p.mu.Unlock()
				})
			} else { // Change occured before timeout -> reset timeframe
				p.Logger.Debugf("change before timeout")
				if !p.hop.Stop() {
					<-p.hop.C
					p.Logger.Debugf("expired ! B")
					p.lastValue = ""
				}
				p.Logger.Debugf("reset")
				p.hop.Reset(time.Second * time.Duration(p.opt.Timeframe))
			}
		}

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

func (p *processor) Stop(e processors.IPacket) error {
	if p.hop != nil {
		p.hop.Stop()
	}

	return nil
}
