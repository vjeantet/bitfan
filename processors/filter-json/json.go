//go:generate bitfanDoc
package json

import (
	"encoding/json"

	"github.com/awillis/bitfan/processors"
)

// Parses JSON events
func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Parses JSON events
type processor struct {
	processors.Base

	opt *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// Allow to skip filter on invalid json
	// @Default false
	SkipOnInvalidJson bool `mapstructure:"skip_on_invalid_json"`

	// The configuration for the JSON filter
	Source string `mapstructure:"source" validate:"required"`

	// Define the target field for placing the parsed data. If this setting is omitted,
	// the JSON data will be stored at the root (top level) of the event
	Target string `mapstructure:"target"`

	// Append values to the tags field when there has been no successful match
	// @Default ["_jsonparsefailure"]
	TagOnFailure []string `mapstructure:"tag_on_failure"`

	// Remove source field if json was parsed successfully
	// @Default true
	RemoveSourceUponSuccess bool `mapstructure:"remove_source_upon_success"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt.TagOnFailure = []string{"_jsonparsefailure"}
	p.opt.SkipOnInvalidJson = false
	p.opt.RemoveSourceUponSuccess = true

	if err := p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}
	return nil
}

func (p *processor) Receive(e processors.IPacket) error {

	jsonString, err := e.Fields().ValueForPathString(p.opt.Source)
	if err != nil {
		p.Logger.Warnf("error while looking for `%s` field : %s", p.opt.Source, err.Error())
		return nil
	}

	byt := []byte(jsonString)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		if p.opt.SkipOnInvalidJson == false {
			p.Logger.Warnf("error while unmarshalling data : %s", err.Error())
			processors.AddTags(p.opt.TagOnFailure, e.Fields())
			p.Send(e, 0)
			return nil
		}
		return nil
	}

	if p.opt.Target != "" {
		e.Fields().SetValueForPath(dat, p.opt.Target)
	} else {
		for k, v := range dat {
			e.Fields().SetValueForPath(v, k)
		}
	}

	if p.opt.RemoveSourceUponSuccess {
		if err := e.Fields().Remove(p.opt.Source); err != nil {
			p.Logger.Warn("error while removing source field: %s", err.Error())
		}
	}

	p.opt.ProcessCommonOptions(e.Fields())

	p.Send(e, 0)
	return nil
}
