package date

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/veino/field"
	"github.com/veino/veino"
)

func New(l veino.Logger) veino.Processor {
	return &processor{}
}

type processor struct {
	Send             veino.PacketSender
	match_field_name string
	match_patterns   []string
	opt              *options
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// The date formats allowed are anything allowed by Golang time format.
	// You can see the docs for this format https://golang.org/src/time/format.go#L20
	// An array with field name first, and format patterns following, [ field, formats... ]
	Match []string

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	Tags []string

	// If this filter is successful, remove arbitrary fields from this event.
	Remove_field []string

	Remove_Tag []string

	// Append values to the tags field when there has been no successful match
	// @default : ["_grokparsefailure"]
	Tag_on_failure []string

	// Store the matching timestamp into the given target field. If not provided,
	// default to updating the @timestamp field of the event
	Target string

	// Specify a time zone canonical ID to be used for date parsing.
	// The valid IDs are listed on IANA Time Zone database, such as "America/New_York".
	// This is useful in case the time zone cannot be extracted from the value,
	// and is not the platform default. If this is not specified the platform default
	//  will be used. Canonical ID is good as it takes care of daylight saving time
	// for you For example, America/Los_Angeles or Europe/Paris are valid IDs.
	// This field can be dynamic and include parts of the event using the %{field} syntax
	Timezone string
}

func (p *processor) Configure(conf map[string]interface{}) error {
	cf := options{Target: "@timestamp", Tag_on_failure: []string{"_dateparsefailure"}}
	if mapstructure.Decode(conf, &cf) != nil {
		return nil
	}
	p.opt = &cf

	p.match_field_name = p.opt.Match[0]
	p.match_patterns = p.opt.Match[1:]

	return nil
}

func (p *processor) Receive(e veino.IPacket) error {
	dated := false
	var value string
	var err error
	value, err = e.Fields().ValueForPathString(p.match_field_name)
	if err == nil {
		for _, layout := range p.match_patterns {
			var t time.Time

			if p.opt.Timezone != "" {
				location, err := time.LoadLocation(p.opt.Timezone)
				if err == nil {
					t, err = time.ParseInLocation(layout, value, location)
				}
			} else {
				t, err = time.Parse(layout, value)
			}

			if err != nil {
				continue
			}

			dated = true
			e.Fields().SetValueForPath(t.Format(veino.VeinoTime), p.opt.Target)
			field.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, "")
			break
		}
	}

	if dated == false {
		field.AddTags(p.opt.Tag_on_failure, e.Fields())
	}

	p.Send(e, 0)
	return nil
}

func (p *processor) Tick(e veino.IPacket) error { return nil }

func (p *processor) Start(e veino.IPacket) error { return nil }

func (p *processor) Stop(e veino.IPacket) error { return nil }
