//go:generate bitfanDoc
// The date filter is used for parsing dates from fields, and then using that date or timestamp as the logstash timestamp for the event.
// For example, syslog events usually have timestamps like this:
// `"Apr 17 09:32:01"`
// You would use the date format MMM dd HH:mm:ss to parse this.
// The date filter is especially important for sorting events and for backfilling old data. If you don’t get the date correct in your event, then searching for them later will likely sort out of order.
// In the absence of this filter, logstash will choose a timestamp based on the first time it sees the event (at input time), if the timestamp is not already set in the event. For example, with file input, the timestamp is set to the time of each read.
package date

import (
	"strconv"
	"time"

	"github.com/vjeantet/bitfan/processors"
	"github.com/vjeantet/jodaTime"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Parses dates from fields to use as the BitFan timestamp for an event
type processor struct {
	processors.Base

	matchFieldName string
	matchPatterns  []string
	opt            *options
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	AddField map[string]interface{} `mapstructure:"add_field"`

	// The date formats allowed are anything allowed by Joda time format.
	// You can see the docs for this format http://www.joda.org/joda-time/key_format.html
	// An array with field name first, and format patterns following, [ field, formats... ]
	Match []string `mapstructure:"match"`

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	AddTag []string `mapstructure:"add_tag"`

	// If this filter is successful, remove arbitrary fields from this event.
	RemoveField []string `mapstructure:"remove_field"`

	RemoveTag []string `mapstructure:"remove_tag"`

	// Append values to the tags field when there has been no successful match
	// Default value is ["_dateparsefailure"]
	TagOnFailure []string `mapstructure:"tag_on_failure"`

	// Store the matching timestamp into the given target field. If not provided,
	// default to updating the @timestamp field of the event
	Target string `mapstructure:"target"`

	// Specify a time zone canonical ID to be used for date parsing.
	// The valid IDs are listed on IANA Time Zone database, such as "America/New_York".
	// This is useful in case the time zone cannot be extracted from the value,
	// and is not the platform default. If this is not specified the platform default
	//  will be used. Canonical ID is good as it takes care of daylight saving time
	// for you For example, America/Los_Angeles or Europe/Paris are valid IDs.
	// This field can be dynamic and include parts of the event using the %{field} syntax
	Timezone string `mapstructure:"timezone"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt.Target = "@timestamp"
	p.opt.TagOnFailure = []string{"_dateparsefailure"}

	if err := p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}

	p.matchFieldName = p.opt.Match[0]
	p.matchPatterns = p.opt.Match[1:]

	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	dated := false
	var value string
	var err error
	value, err = e.Fields().ValueForPathString(p.matchFieldName)
	if err == nil {
		for _, layout := range p.matchPatterns {
			var t time.Time

			if p.opt.Timezone != "" {
				t, err = jodaTime.ParseInLocation(layout, value, p.opt.Timezone)
			} else {
				if layout == "UNIX" {
					var i int64
					i, err = strconv.ParseInt(value, 10, 64)
					if err == nil {
						t = time.Unix(i, 0)
					}
				} else {
					t, err = jodaTime.Parse(layout, value)
				}

			}

			if err != nil {
				continue
			}

			dated = true
			e.Fields().SetValueForPath(t, p.opt.Target)
			processors.ProcessCommonFields2(e.Fields(),
				p.opt.AddField,
				p.opt.AddTag,
				p.opt.RemoveField,
				p.opt.RemoveTag,
			)
			break
		}
	}

	if dated == false {
		processors.AddTags(p.opt.TagOnFailure, e.Fields())
	}

	p.Send(e, 0)
	return nil
}
