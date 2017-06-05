//go:generate bitfanDoc -codec line
package linecodec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

type decoder struct {
	more    bool
	r       *bufio.Scanner
	options options
}

type encoder struct {
	w         io.Writer
	options   options
	formatTmp *template.Template
}

type options struct {
	// Change the delimiter that separates lines
	// @Default "\\n"
	Delimiter string

	// Format (Encoder) as a golang text/template
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
		options: options{
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
	return nil
}

func NewDecoder(r io.Reader) *decoder {
	d := &decoder{
		r:    bufio.NewScanner(r),
		more: true,
		options: options{
			Delimiter: "\n",
		},
	}

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Return nothing if at end of file and no data passed
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Find the index of the input of a newline followed by a
		// pound sign.
		if i := strings.Index(string(data), d.options.Delimiter); i >= 0 {
			return i + 1, data[0:i], nil
		}

		// If at end of file with data return the data
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	d.r.Split(split)

	return d
}

func (d *decoder) SetOptions(conf map[string]interface{}) error {
	if err := mapstructure.Decode(conf, &d.options); err != nil {
		return err
	}
	return nil
}

func (d *decoder) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	if true == d.r.Scan() {
		d.more = true
		data["message"] = d.r.Text()
	} else {
		d.more = false
		return data, io.EOF
	}

	return data, nil
}

func (d *decoder) More() bool {
	return d.more
}
