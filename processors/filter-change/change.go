//go:generate bitfanDoc
// This rule will monitor a certain field and match if that field changes. The field must change with respect to the last event
package change

import "github.com/vjeantet/bitfan/processors"

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
	opt       *options
	lastValue string
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

	// If true, events without a compare_key field will not count as changed.
	// @Default true
	IgnoreNull bool `mapstructure:"ignore_null"`

	// The maximum time in seconds between changes. After this time period, ElastAlert will forget the old value of the compare_key field.
	// Timeframe string `mapstructure:"timeframe"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		IgnoreNull: true,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {

	if p.opt.IgnoreNull == true && e.Fields().Exists(p.opt.CompareField) == false {
		p.Logger.Debugf("event does not have a field [%s]", p.opt.CompareField)
	}

	if p.lastValue != e.Fields().ValueOrEmptyForPathString(p.opt.CompareField) {
		p.Logger.Debugf("[%s] value change from '%s' to '%s'", p.opt.CompareField, p.lastValue, e.Fields().ValueOrEmptyForPathString(p.opt.CompareField))

		p.lastValue = e.Fields().ValueOrEmptyForPathString(p.opt.CompareField)

		processors.ProcessCommonFields2(e.Fields(),
			p.opt.AddField,
			p.opt.AddTag,
			p.opt.RemoveField,
			p.opt.RemoveTag,
		)
		p.Send(e, 0)
	}

	return nil
}
