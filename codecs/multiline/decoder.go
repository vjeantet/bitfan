//go:generate bitfanDoc -codec multiline
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
	"github.com/vjeantet/bitfan/codecs/lib"
)

// Merges multiline messages into a single event
type decoder struct {
	more    bool
	r       *bufio.Scanner
	options decoderOptions
	memory  string

	log lib.Logger
}

//
type decoderOptions struct {
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

func NewDecoder(r io.Reader) *decoder {
	d := &decoder{
		r:    bufio.NewScanner(r),
		more: true,
		options: decoderOptions{
			Delimiter: "\n",
			Negate:    false,
			Pattern:   "\\s",
			What:      "dd",
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
func (d *decoder) SetOptions(conf map[string]interface{}, logger lib.Logger, cwl string) error {
	d.log = logger

	if err := mapstructure.Decode(conf, &d.options); err != nil {
		return err
	}

	return nil
}

func (d *decoder) Decode(v *interface{}) error {
	for d.r.Scan() {
		d.more = true
		match, _ := regexp.MatchString(d.options.Pattern, d.r.Text())

		if d.options.Negate {
			match = !match
		}

		if d.options.What == "previous" {
			if match == true { // stick to previous
				d.memory += d.options.Delimiter + d.r.Text()
				continue
			} else if d.memory == "" {
				d.memory = d.r.Text()
				continue
			} else {
				*v = d.memory
				d.memory = d.r.Text()
				return nil
			}
		}
		if d.options.What == "next" {
			if match == true { // stick to au previous
				if d.memory != "" {
					d.memory += "\n"
				}
				d.memory += d.r.Text()
				continue
			} else if d.memory == "" {
				*v = d.r.Text()
				d.memory = ""
				return nil
			} else {
				d.memory += "\n" + d.r.Text()
				*v = d.memory
				d.memory = ""
				return nil
			}
		}
	}
	d.more = false
	return io.EOF
}

func (d *decoder) More() bool {
	return d.more
}
