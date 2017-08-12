//go:generate bitfanDoc
package html

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vjeantet/bitfan/processors"
)

const (
	// all events
	PORT_SUCCESS = 0
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	Codec string

	// Which field contains the html document
	SourceField string `mapstructure:"source_field"`

	// Add fields with the text of elements found with css selector
	Text map[string]string `mapstructure:"text"`

	// Add fields with the number of elements found with css selector
	Size map[string]string `mapstructure:"size"`

	// Append values to the tags field when the html document can not be parsed
	// @default : ["_htmlparsefailure"]
	TagOnFailure []string `mapstructure:"tag_on_failure"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		SourceField:  "message",
		TagOnFailure: []string{"_htmlparsefailure"},
	}
	p.opt = &defaults

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(e.Fields().ValueOrEmptyForPathString(p.opt.SourceField)))

	if err != nil {
		tags, err := e.Fields().ValueForPath("tags")
		if err != nil {
			tags = []string{}
		}
		newtags := append(tags.([]string), p.opt.TagOnFailure...)
		e.Fields().SetValueForPath(newtags, "tags")
		p.Send(e, PORT_SUCCESS)
		return nil
	} else {
		// Text
		for fkey, query := range p.opt.Text {
			sel := doc.Find(query)
			if sel.Length() == 1 {
				e.Fields().SetValueForPath(sel.Text(), fkey)
			} else if sel.Length() > 1 {
				values := []string{}
				for i := range sel.Nodes {
					s := sel.Eq(i)
					values = append(values, s.Text())
				}
				e.Fields().SetValueForPath(values, fkey)
			} else {
				e.Fields().SetValueForPath("", fkey)
			}
		}

		// Size
		for fkey, query := range p.opt.Size {
			sel := doc.Find(query)
			e.Fields().SetValueForPath(sel.Length(), fkey)
		}
	}

	p.opt.ProcessCommonOptions(e.Fields())
	p.Send(e, PORT_SUCCESS)
	return nil
}
