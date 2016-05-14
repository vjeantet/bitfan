package grok

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/veino/field"
	"github.com/veino/veino"
	"github.com/vjeantet/grok"
)

const (
	PORT_SUCCESS = 0
)

func New(l veino.Logger) veino.Processor {
	return &processor{}
}

type processor struct {
	Send veino.PacketSender
	grok *grok.Grok

	// If this filter is successful, add any arbitrary fields to this event. Field names can
	// be dynamic and include parts of the event using the %{field}.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	Add_tag []string

	// Break on first match. The first successful match by grok will result in the filter being
	// finished. If you want grok to try all patterns (maybe you are parsing different things),
	// then set this to false
	// @default : true
	Break_on_match bool

	// If true, keep empty captures as event fields
	Keep_empty_captures bool

	// A hash of matches of field ⇒ value
	// TODO : keep order
	Match map[string]string

	// If true, only store named captures from grok.
	// @default : true
	Named_captures_only bool

	// Veino ships by default with a bunch of patterns, so you don’t necessarily need to
	// define this yourself unless you are adding additional patterns. You can point to
	// multiple pattern directories using this setting Note that Grok will read all files
	// in the directory and assume its a pattern file (including any tilde backup files)
	Patterns_dir []string

	// If this filter is successful, remove arbitrary fields from this event
	Remove_field []string

	// If this filter is successful, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	Remove_tag []string

	// Append values to the tags field when there has been no successful match
	// @default : ["_grokparsefailure"]
	Tag_on_failure []string
}

func (p *processor) Configure(conf map[string]interface{}) error {
	p.Named_captures_only = true
	p.Break_on_match = true
	p.Tag_on_failure = []string{"_grokparsefailure"}

	var err error
	if err = mapstructure.Decode(conf, p); err != nil {
		return err
	}

	p.grok, err = grok.NewWithConfig(&grok.Config{
		NamedCapturesOnly: p.Named_captures_only,
		RemoveEmptyValues: !p.Keep_empty_captures,
	})
	if err != nil {
		return err
	}

	for _, path := range p.Patterns_dir {
		err := p.grok.AddPatternsFromPath(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *processor) Receive(e veino.IPacket) error {
	groked := false
	for fkey, pattern := range p.Match {
		values, _ := p.grok.Parse(pattern, e.Fields().ValueOrEmptyForPathString(fkey))
		if len(values) > 0 {
			groked = true
			// if f, err := mxj.Ma(values); err == nil {
			if err := mapstructure.Decode(values, e.Fields()); err != nil {
				return fmt.Errorf("Error while groking : %s", err.Error())
			}
			if p.Break_on_match == true {
				break
			}
		}
	}

	if groked {
		field.AddFields(p.Add_field, e.Fields())
		field.RemoveFields(p.Remove_field, e.Fields())
		if len(p.Add_tag) > 0 {
			field.AddTags(p.Add_tag, e.Fields())
		}
		field.RemoveTags(p.Remove_tag, e.Fields())

	}

	if !groked {
		tags, err := e.Fields().ValueForPath("tags")
		if err != nil {
			tags = []string{}
		}
		newtags := append(tags.([]string), p.Tag_on_failure...)
		e.Fields().SetValueForPath(newtags, "tags")
	}

	p.Send(e, PORT_SUCCESS)
	return nil
}

func (p *processor) Stop(e veino.IPacket) error { return nil }

func (p *processor) Tick(e veino.IPacket) error { return nil }

func (p *processor) Start(e veino.IPacket) error { return nil }
