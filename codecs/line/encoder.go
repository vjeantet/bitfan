//go:generate bitfanDoc -codec encoder
// doc codec
package linecodec

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/codecs/lib"
	"github.com/vjeantet/bitfan/core/location"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

// doc encoder
type encoder struct {
	w         io.Writer
	options   encoderOptions
	formatTpl *template.Template

	log lib.Logger
}

// doc encoderOptions
type encoderOptions struct {
	// Change the delimiter that separates lines
	// @Default "\\n"
	Delimiter string

	// Format as a golang text/template
	// @Default "{{Timestamp .}} {{.host}} {{.message}}"
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
			Delimiter: "\n",
			Format:    "{{Timestamp .}} {{.host}} {{.message}}",
		},
	}

	return e
}

func (e *encoder) SetOptions(conf map[string]interface{}, logger lib.Logger, cwl string) error {
	e.log = logger

	if err := mapstructure.Decode(conf, &e.options); err != nil {
		return err
	}

	if e.options.Format != "" {

		//TODO : add a location.TemplateWithOptions to return golang text/template

		loc, err := location.NewLocation(e.options.Format, cwl)
		if err != nil {
			return err
		}

		content, _, err := loc.ContentWithOptions(e.options.Var)
		if err != nil {
			return err
		}
		e.options.Format = string(content)

		funcMap := template.FuncMap{
			"Timestamp": func(m map[string]interface{}) string {
				return m["@timestamp"].(time.Time).Format(timeFormat)
			},
		}

		e.formatTpl, err = template.New("format").Funcs(funcMap).Parse(e.options.Format)
		if err != nil {
			fmt.Errorf("stdout Format tpl error : %s", err)
			return err
		}
	}

	return nil
}

func (e *encoder) Encode(data map[string]interface{}) error {
	buff := bytes.NewBufferString("")
	e.formatTpl.Execute(buff, data)
	e.w.Write(buff.Bytes())
	e.w.Write([]byte(e.options.Delimiter))
	return nil
}
