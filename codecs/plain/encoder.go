//go:generate bitfanDoc -codec encoder
// doc codec
package plaincodec

import (
	"bytes"
	"io"
	"text/template"

	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/commons"
)

// doc encoder
type encoder struct {
	w         io.Writer
	options   encoderOptions
	formatTpl *template.Template

	log commons.Logger
}

// doc encoderOptions
type encoderOptions struct {

	// Format as a golang text/template
	// @Default "{{.message}}"
	// @Type Location
	Format string `mapstructure:"format"`

	// You can set variable to be used in Statements by using ${var}.
	// each reference will be replaced by the value of the variable found in Statement's content
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`
}

func NewEncoder(w io.Writer) *encoder {
	e := &encoder{
		w: w,
		options: encoderOptions{
			Format: `{{.message}}`,
		},
	}

	return e
}

func (e *encoder) SetOptions(conf map[string]interface{}, logger commons.Logger, cwl string) error {
	e.log = logger

	if err := mapstructure.Decode(conf, &e.options); err != nil {
		return err
	}

	if e.options.Format != "" {
		//TODO : add a location.TemplateWithOptions to return golang text/template

		loc, err := commons.NewLocation(e.options.Format, cwl)
		if err != nil {
			return err
		}

		e.formatTpl, _, err = loc.TemplateWithOptions(e.options.Var)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *encoder) Encode(data map[string]interface{}) error {
	buff := bytes.NewBufferString("")
	e.formatTpl.Execute(buff, data)
	e.w.Write(buff.Bytes())
	return nil
}
