//go:generate bitfanDoc -codec codec
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

type codec struct {
	more    bool
	r       *bufio.Scanner
	w       io.Writer
	options options

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

func New(opt map[string]interface{}) *codec {
	var err error
	d := &codec{
		more: true,
		options: options{
			Delimiter: "\n",
			Format:    "{{Timestamp .}} {{.host}} {{.message}}\n",
		},
	}
	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}

	if d.options.Format != "" {

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

		d.formatTmp, err = template.New("format").Funcs(funcMap).Parse(d.options.Format)
		if err != nil {
			fmt.Errorf("stdout Format tpl error : %s", err)
			return nil
		}
	}

	return d
}

func (c *codec) Encoder(w io.Writer) *codec {
	c.w = w
	return c
}

func (p *codec) Encode(data map[string]interface{}) error {
	buff := bytes.NewBufferString("")
	p.formatTmp.Execute(buff, data)
	p.w.Write(buff.Bytes())
	return nil
}

func (c *codec) Decoder(r io.Reader) *codec {
	c.r = bufio.NewScanner(r)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Return nothing if at end of file and no data passed
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Find the index of the input of a newline followed by a
		// pound sign.
		if i := strings.Index(string(data), c.options.Delimiter); i >= 0 {
			return i + 1, data[0:i], nil
		}

		// If at end of file with data return the data
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}

	c.r.Split(split)

	return c
}

func (c *codec) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	if true == c.r.Scan() {
		c.more = true
		data["message"] = c.r.Text()
	} else {
		c.more = false
		return data, io.EOF
	}

	return data, nil
}

func (c *codec) More() bool {
	return c.more
}
