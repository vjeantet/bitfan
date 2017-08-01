//go:generate bitfanDoc
// A simple output which prints to the STDOUT of the shell running BitFan. This output can be quite convenient when debugging plugin configurations, by allowing instant access to the event data after it has passed through the inputs and filters.
//
// For example, the following output configuration, in conjunction with the BitFan -e command-line flag, will allow you to see the results of your event pipeline for quick iteration.
// ```
// output {
//   stdout {}
// }
// ```
// Useful codecs include:
//
// pp: outputs event data using the go "k0kubun/pp" package
// if codec is rubydebug, it will treated as "pp"
// ```
// output {
//   stdout { codec => pp }
// }
// ```
// json: outputs event data in structured JSON format
// ```
// output {
//   stdout { codec => json }
// }
// ```
package stdout

import (
	"os"

	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/processors"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

const (
	CODEC_PRETTYPRINT string = "pp"
	CODEC_LINE        string = "line"
	CODEC_RUBYDEBUG   string = "rubydebug"
	CODEC_JSON        string = "json"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// Codec can be one of  "json", "line", "pp" or "rubydebug"
	// @ExampleLS codec => "pp"
	// @Default "line"
	// @Enum "json","line","pp","rubydebug"
	// @Type Codec
	Codec codecs.Codec `mapstructure:"codec"`
}

// Prints events to the standard output
type processor struct {
	processors.Base

	opt *options

	enc codecs.Encoder
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {

	defaults := options{
		Codec: codecs.New("line", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
	}

	p.opt = &defaults
	if err := p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}

	var err error
	p.enc, err = p.opt.Codec.NewEncoder(os.Stdout)
	if err != nil {
		p.Logger.Errorln("codec error : ", err.Error())
		return err
	}

	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	p.enc.Encode(e.Fields().Old())
	p.Memory.Set("last", e.Fields().StringIndentNoTypeInfo(2))
	p.Send(e)
	return nil
}
