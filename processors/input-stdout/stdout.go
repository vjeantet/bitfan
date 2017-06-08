//go:generate bitfanDoc
// A simple output which prints to the STDOUT of the shell running BitFan.

package inputstdout

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

	// WebHook *core.WebHook
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
	return nil
}

func (p *processor) Start(e processors.IPacket) error {
	p.enc.Encode(e.Fields().Old())
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	return nil
}
