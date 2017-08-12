//go:generate bitfanDoc
// This filter helps automatically parse messages (or specific event fields)
// which are of the foo=bar variety.
package kv

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vjeantet/bitfan/processors"
)

const (
	PORT_SUCCESS = 0
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Parses key-value pairs
type processor struct {
	processors.Base
	opt     *options
	scan_re *regexp.Regexp
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// A bool option for removing duplicate key/value pairs.
	// When set to false, only one unique key/value pair will be preserved.
	// For example, consider a source like from=me from=me.
	// [from] will map to an Array with two elements: ["me", "me"].
	// to only keep unique key/value pairs, you could use this configuration
	// ```
	// kv {
	//   allow_duplicate_values => false
	// }
	// ```
	AllowDuplicateValues bool `mapstructure:"allow_duplicate_values"`

	// A hash specifying the default keys and their values which should be added
	// to the event in case these keys do not exist in the source field being parsed.
	//
	// Example
	// ```
	// kv {
	//   default_keys => { "from"=> "logstash@example.com",
	//                    "to"=> "default@dev.null" }
	// }
	// ```
	DefaultKeys map[string]interface{} `mapstructure:"default_keys"`

	// An array specifying the parsed keys which should not be added to the event.
	//
	// By default no keys will be excluded.
	//
	// For example, consider a source like Hey, from=<abc>, to=def foo=bar.
	//
	// To exclude from and to, but retain the foo key, you could use this configuration:
	// ```
	// kv {
	//   exclude_keys => [ "from", "to" ]
	// }
	// ```
	ExcludeKeys []string `mapstructure:"exclude_keys"`

	// A string of characters to use as delimiters for parsing out key-value pairs.
	//
	// These characters form a regex character class and thus you must escape special regex characters like [ or ] using \.
	// #### Example with URL Query Strings
	// For example, to split out the args from a url query string such as ?pin=12345~0&d=123&e=foo@bar.com&oq=bobo&ss=12345:
	// ```
	//  kv {
	//    field_split => "&?"
	//  }
	// ```
	// The above splits on both & and ? characters, giving you the following fields:
	//
	// * pin: 12345~0
	// * d: 123
	// * e: foo@bar.com
	// * oq: bobo
	// * ss: 12345
	FieldSplit string `mapstructure:"field_split"`

	// A boolean specifying whether to include brackets as value wrappers (the default is true)
	// ```
	// kv {
	//   include_brackets => true
	// }
	// ```
	// For example, the result of this line: bracketsone=(hello world) bracketstwo=[hello world]
	// will be:
	//
	// * bracketsone: hello world
	// * bracketstwo: hello world
	//
	// instead of:
	//
	// * bracketsone: (hello
	// * bracketstwo: [hello
	IncludeBrackets bool `mapstructure:"include_brackets"`

	// An array specifying the parsed keys which should be added to the event. By default all keys will be added.
	//
	// For example, consider a source like Hey, from=<abc>, to=def foo=bar. To include from and to, but exclude the foo key, you could use this configuration:
	// ```
	// kv {
	//   include_keys => [ "from", "to" ]
	// }
	// ```
	IncludeKeys []string `mapstructure:"include_keys"`

	// A string to prepend to all of the extracted keys.
	//
	// For example, to prepend arg_ to all keys:
	// ```
	// kv {
	//   prefix => "arg_" }
	// }
	// ```
	Prefix string

	// A boolean specifying whether to drill down into values and recursively get more key-value pairs from it. The extra key-value pairs will be stored as subkeys of the root key.
	//
	// Default is not to recursive values.
	// ```
	// kv {
	//  recursive => "true"
	// }
	// ```
	Recursive bool

	// The field to perform key=value searching on
	//
	// For example, to process the not_the_message field:
	// ```
	// kv { source => "not_the_message" }
	// ```
	Source string

	// The name of the container to put all of the key-value pairs into.
	//
	// If this setting is omitted, fields will be written to the root of the event, as individual fields.
	//
	// For example, to place all keys into the event field kv:
	// ```
	// kv { target => "kv" }
	// ```
	Target string

	// A string of characters to trim from the value. This is useful if your values are wrapped in brackets or are terminated with commas (like postfix logs).
	//
	// For example, to strip <, >, [, ] and , characters from values:
	// ```
	// kv {
	//   trim => "<>[],"
	// }
	// ```
	Trim string

	// A string of characters to trim from the key. This is useful if your keys are wrapped in brackets or start with space.
	//
	// For example, to strip < > [ ] and , characters from keys:
	// ```
	// kv {
	//   trimkey => "<>[],"
	// }
	// ```
	Trimkey string `mapstructure:"trimkey"`

	// A string of characters to use as delimiters for identifying key-value relations.
	//
	// These characters form a regex character class and thus you must escape special regex characters like [ or ] using \.
	//
	// For example, to identify key-values such as key1:value1 key2:value2:
	// ```
	// { kv { value_split => ":" }
	// ```
	ValueSplit string `mapstructure:"value_split"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		AllowDuplicateValues: true,
		FieldSplit:           " ",
		IncludeBrackets:      true,
		Recursive:            false,
		Source:               "message",
		ValueSplit:           "=",
	}
	p.opt = &defaults

	if err = p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}

	valueRxString := `(?:\"([^\"]+)\"|'([^']+)'`

	if p.opt.IncludeBrackets == true {
		valueRxString += `|\(([^\)]+)\)|\[([^\]]+)\]|<([^>]+)>`
	}

	valueRxString += `|((?:\\ |[^` + p.opt.FieldSplit + `])+))`
	p.scan_re, _ = regexp.Compile(`((?:\\ |[^` + p.opt.FieldSplit + p.opt.ValueSplit + `])+)\s*[` + p.opt.ValueSplit + `]\s*` + valueRxString)
	return nil
}

func (p *processor) parse(value string, e processors.IPacket, kv map[string]interface{}) (map[string]interface{}, error) {
	// Short circuit parsing if the text does not contain the p.opt.value_split
	if strings.Contains(value, p.opt.ValueSplit) == false {
		return kv, nil
	}

	for _, matches := range p.scan_re.FindAllStringSubmatch(value, -1) {
		key := matches[1]
		// trimkey
		if p.opt.Trimkey != "" {
			key = strings.Trim(key, p.opt.Trimkey)
		}

		// handle only IncludeKeys if set
		if len(p.opt.IncludeKeys) > 0 {
			if stringInSlice(key, p.opt.IncludeKeys) == false {
				continue
			}
		} else {
			// excludekeys is set
			if stringInSlice(key, p.opt.ExcludeKeys) {
				continue
			}
		}

		// prefix
		value := strings.Join(matches[2:], "")
		value = fmt.Sprintf("%s%s", p.opt.Prefix, value)

		// Trim Value
		if p.opt.Trim != "" {
			value = strings.Trim(value, p.opt.Trim)
		}

		// handle multiples values for key
		if _, ok := kv[key]; ok {
			// handle allow_duplicate_values
			if p.opt.AllowDuplicateValues == true {

				switch kv[key].(type) {
				case string:
					kv[key] = []string{kv[key].(string), value}
				case []string:
					kv[key] = append(kv[key].([]string), value)
				}
			}
		} else {
			kv[key] = value
		}

	}
	return kv, nil
}

func (p *processor) Receive(e processors.IPacket) error {

	value, _ := e.Fields().ValueForPath(p.opt.Source)

	kv := map[string]interface{}{}
	switch value.(type) {
	case string:
		kv, _ = p.parse(value.(string), e, kv)
	case []string:
		for _, pvalue := range value.([]string) {
			kv, _ = p.parse(pvalue, e, kv)
		}
	default:
		p.Logger.Warn("kv filter has no support for this type of data")
	}

	if len(kv) > 0 {
		// Set Values
		// # Use Target is mentionned
		if p.opt.Target != "" {
			e.Fields().SetValueForPath(kv, p.opt.Target)
		} else {
			for k, v := range kv {
				e.Fields().SetValueForPath(v, k)
			}
		}

		// # Add default key-values for missing keys
		for k, v := range p.opt.DefaultKeys {
			path := k
			if p.opt.Target != "" {
				path = p.opt.Target + "." + k
			}
			if e.Fields().Exists(path) == false {
				e.Fields().SetValueForPath(v, path)
			}
		}

		p.opt.ProcessCommonOptions(e.Fields())

	}

	p.Send(e, PORT_SUCCESS)
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
