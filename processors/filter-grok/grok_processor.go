//go:generate bitfanDoc
package grok

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/processors"
	"github.com/vjeantet/grok"
)

const (
	PORT_SUCCESS = 0
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt  *options
	grok *grok.Grok
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// Break on first match. The first successful match by grok will result in the filter being
	// finished. If you want grok to try all patterns (maybe you are parsing different things),
	// then set this to false
	// @default : true
	BreakOnMatch bool `mapstructure:"break_on_match"`

	// If true, keep empty captures as event fields
	KeepEmptyCaptures bool `mapstructure:"keep_empty_captures"`

	// A hash of matches of field ⇒ value
	// @nodefault
	//
	// For example:
	// ```
	//     filter {
	//       grok { match => { "message" => "Duration: %{NUMBER:duration}" } }
	//     }
	// ```
	// If you need to match multiple patterns against a single field, the value can be an array of patterns
	// ```
	//     filter {
	//       grok { match => { "message" => [ "Duration: %{NUMBER:duration}", "Speed: %{NUMBER:speed}" ] } }
	//     }
	// ```
	Match map[string][]string `mapstructure:"match" validate:"required"`

	// If true, only store named captures from grok.
	// @default : true
	NamedCapturesOnly bool `mapstructure:"named_capture_only"`

	// BitFan ships by default with a bunch of patterns, so you don’t necessarily need to
	// define this yourself unless you are adding additional patterns. You can point to
	// multiple pattern directories using this setting Note that Grok will read all files
	// in the directory and assume its a pattern file (including any tilde backup files)
	// @default : []
	PatternsDir []string `mapstructure:"patterns_dir"`

	// Append values to the tags field when there has been no successful match
	// @default : ["_grokparsefailure"]
	TagOnFailure []string `mapstructure:"tag_on_failure"`
}

func (p *processor) fixGrokMatch(conf interface{}) (map[string][]string, error) {
	fixMatch := make(map[string][]string)
	switch m := conf.(type) {
	case map[string]interface{}:
		for fieldName, patterns := range m {
			switch v := patterns.(type) {
			case []string:
				fixMatch[fieldName] = v
			case string:
				fixMatch[fieldName] = []string{v}
			default:
				return fixMatch, fmt.Errorf("unsupported match value format %v", reflect.TypeOf(v))
			}
		}
	default:
		return fixMatch, fmt.Errorf("unsupported match format : %s", reflect.TypeOf(m))
	}

	return fixMatch, nil
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		NamedCapturesOnly: true,
		BreakOnMatch:      true,
		TagOnFailure:      []string{"_grokparsefailure"},
	}
	p.opt = &defaults

	conf["match"], err = p.fixGrokMatch(conf["match"])
	if err != nil {
		return err
	}

	if err = p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}

	p.grok, err = grok.NewWithConfig(&grok.Config{
		NamedCapturesOnly: p.opt.NamedCapturesOnly,
		RemoveEmptyValues: !p.opt.KeepEmptyCaptures,
	})
	if err != nil {
		return err
	}

	for _, path := range p.opt.PatternsDir {
		err := p.grok.AddPatternsFromPath(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	groked := false
	var values map[string]string

	for fkey, patterns := range p.opt.Match {

		if p.opt.BreakOnMatch && len(values) > 0 {
			break
		}

		for i := 0; i < len(patterns); i++ {
			values, _ = p.grok.Parse(patterns[i], e.Fields().ValueOrEmptyForPathString(fkey))
			if len(values) > 0 {
				groked = true

				if err := mapstructure.Decode(values, e.Fields()); err != nil {
					return fmt.Errorf("Error while groking : %v", err)
				}
				if p.opt.BreakOnMatch {
					break
				}
			}
		}
	}

	if groked {
		p.opt.ProcessCommonOptions(e.Fields())

	}

	if !groked {
		processors.AddTags(p.opt.TagOnFailure, e.Fields())
	}

	p.Send(e, PORT_SUCCESS)
	return nil
}
