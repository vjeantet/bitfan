//go:generate bitfanDoc -codec encoder
package linecodec

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

type encoder struct {
	w         io.Writer
	options   encoderOptions
	formatTmp *template.Template
}

type encoderOptions struct {
	// Change the delimiter that separates lines
	// @Default "\\n"
	Delimiter string

	// Format as a golang text/template
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
			Format:    "{{Timestamp .}} {{.host}} {{.message}}\n",
		},
	}

	return e
}

func (e *encoder) SetOptions(conf map[string]interface{}) error {
	var err error
	if err := mapstructure.Decode(conf, &e.options); err != nil {
		return err
	}

	if e.options.Format != "" {

		// TODO !
		// loc, err := location.NewLocation(d.options.Format, d.ConfigWorkingLocation)
		// if err != nil {
		// 	return err
		// }

		// content, _, err := loc.ContentWithOptions(d.options.V)
		// if err != nil {
		// 	return err
		// }
		// d.options.Format = string(content)

		funcMap := template.FuncMap{
			"Timestamp": func(m map[string]interface{}) string {
				return m["@timestamp"].(time.Time).Format(timeFormat)
			},
		}

		e.formatTmp, err = template.New("format").Funcs(funcMap).Parse(e.options.Format)
		if err != nil {
			fmt.Errorf("stdout Format tpl error : %s", err)
			return err
		}
	}

	return nil
}

func (e *encoder) Encode(data map[string]interface{}) error {
	buff := bytes.NewBufferString("")
	e.formatTmp.Execute(buff, data)
	e.w.Write(buff.Bytes())
	e.w.Write([]byte(e.options.Delimiter))
	return nil
}
