//go:generate bitfanDoc -codec codec
// The multiline codec will collapse multiline messages and merge them into a single event.
//
// The original goal of this codec was to allow joining of multiline messages from files into a single event. For example, joining Java exception and stacktrace messages into a single event.
//
// The config looks like this:
// ```
// input {
//   stdin {
//     codec => multiline {
//       pattern => "pattern, a regexp"
//       negate => true or false
//       what => "previous" or "next"
//     }
//   }
// }
// ```
// The pattern should match what you believe to be an indicator that the field is part of a multi-line event.
//
// The what must be previous or next and indicates the relation to the multi-line event.
//
// The negate can be true or false (defaults to false). If true, a message not matching the pattern will constitute a match of the multiline filter and the what will be applied. (vice-versa is also true)
//
// For example, Java stack traces are multiline and usually have the message starting at the far-left, with each subsequent line indented. Do this:
//
// ```
// input {
//   stdin {
//     codec => multiline {
//       pattern => "^\\s"
//       what => "previous"
//     }
//   }
// }
// ```
// This says that any line starting with whitespace belongs to the previous line.
//
// Another example is to merge lines not starting with a date up to the previous line..
//
// ```
// input {
//   file {
//     path => "/var/log/someapp.log"
//     codec => multiline {
//       # Grok pattern names are valid! :)
//       pattern => "^%{TIMESTAMP_ISO8601} "
//       negate => true
//       what => "previous"
//     }
//   }
// }
// ```
// This says that any line not starting with a timestamp should be merged with the previous line.
//
// One more common example is C line continuations (backslash). Here’s how to do that:
//
// ```
// filter {
//   multiline {
//     pattern => "\\$"
//     what => "next"
//   }
// }
// ```
// This says that any line ending with a backslash should be combined with the following line.
package multilinecodec

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Merges multiline messages into a single event
type codec struct {
	more    bool
	r       *bufio.Scanner
	options options
	memory  string
}

//
type options struct {
	// Change the delimiter that separates lines
	// @Default "\n"
	Delimiter string `mapstructure:"delimiter"`

	// Negate the regexp pattern (if not matched).
	// @Default false
	Negate bool `mapstructure:"negate"`

	// The regular expression to match
	// @ExampleLS pattern => "^\\s"
	Pattern string `mapstructure:"pattern" validate:"required"`

	// If the pattern matched, does event belong to the next or previous event?
	// @Enum previous,next
	// @Default "previous"
	What string `mapstructure:"what"`
}

func New(opt map[string]interface{}) *codec {
	d := &codec{
		more: true,
		options: options{
			Delimiter: "\n",
			Negate:    false,
			Pattern:   "\\s",
			What:      "dd",
		},
	}

	if err := mapstructure.Decode(opt, &d.options); err != nil {
		return nil
	}

	return d
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

func (p *codec) Decode() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	for p.r.Scan() {
		p.more = true
		match, _ := regexp.MatchString(p.options.Pattern, p.r.Text())

		if p.options.Negate {
			match = !match
		}

		if p.options.What == "previous" {
			if match == true { // stick to previous
				p.memory += p.options.Delimiter + p.r.Text()
				continue
			} else if p.memory == "" {
				p.memory = p.r.Text()
				continue
			} else {
				data["message"] = p.memory
				p.memory = p.r.Text()
				return data, nil
			}
		}
		if p.options.What == "next" {
			if match == true { // stick to au previous
				if p.memory != "" {
					p.memory += "\n"
				}
				p.memory += p.r.Text()
				continue
			} else if p.memory == "" {
				data["message"] = p.r.Text()
				p.memory = ""
				return data, nil
			} else {
				p.memory += "\n" + p.r.Text()
				data["message"] = p.memory
				p.memory = ""
				return data, nil
			}
		}
	}
	p.more = false
	return data, io.EOF

	// if p.options.What == "next" {
	// 	if match == true { // coller au previous
	// 		p.memory += "\n" + p.r.Text()
	// 		return nil, nil
	// 	} else if p.memory == "" {
	// 		data["message"] = p.r.Text()
	// 		p.memory = ""
	// 		return data, nil
	// 	} else {
	// 		p.memory += "\n" + p.r.Text()
	// 		data["message"] = p.memory
	// 		p.memory = ""
	// 		return data, nil
	// 	}
	// }

	// if p.options.What == "previous" {
	// 	if match == true { // coller au previous
	// 		p.memory += "\n" + p.r.Text()
	// 		return nil, nil
	// 	} else if p.memory == "" {
	// 		p.memory = p.r.Text()
	// 		return nil, nil
	// 	} else {
	// 		data["message"] = p.memory
	// 		p.memory = p.r.Text()
	// 		return data, nil
	// 	}
	// }

	return data, nil
}

func (p *codec) DecodeReader(r io.Reader) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	// var cr io.Reader

	// var err error
	// cr, err = charset.NewReaderLabel(p.options.Charset, r)
	// if err != nil {
	// 	return nil, err
	// }
	// bytes, err := ioutil.ReadAll(cr)
	// if err != nil {
	// 	return data, err
	// }

	// match, _ := regexp.MatchString(p.options.Pattern, string(bytes))

	// if p.options.What == "next" {
	// 	if match == true { // coller au previous
	// 		p.memory += "\n" + string(bytes)
	// 		return nil, nil
	// 	} else if p.memory == "" {
	// 		data["message"] = string(bytes)
	// 		p.memory = ""
	// 		return data, nil
	// 	} else {
	// 		p.memory += "\n" + string(bytes)
	// 		data["message"] = p.memory
	// 		p.memory = ""
	// 		return data, nil
	// 	}
	// }

	// if p.options.What == "previous" {
	// 	if match == true { // coller au previous
	// 		p.memory += "\n" + (string(bytes))
	// 		return nil, nil
	// 	} else if p.memory == "" {
	// 		p.memory = string(bytes)
	// 		return nil, nil
	// 	} else {
	// 		data["message"] = p.memory
	// 		p.memory = string(bytes)
	// 		return data, nil
	// 	}
	// }

	return data, nil
}

func (p *codec) More() bool {
	return p.more
}
