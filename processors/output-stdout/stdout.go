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

	// WebHook *core.WebHook
	opt *options
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {

	defaults := options{
		Codec: codecs.New("line"),
	}

	p.opt = &defaults
	if err := p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}

	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	var enc codecs.Encoder
	var err error

	enc, err = p.opt.Codec.Encoder(os.Stdout)
	if err != nil {
		p.Logger.Errorln("encoder error : ", err.Error())
		return err
	}

	enc.Encode(e.Fields().Old())

	// switch p.opt.Codec.Name {
	// case CODEC_LINE:
	// 	buff := bytes.NewBufferString("")
	// 	p.formatTmp.Execute(buff, e.Fields())
	// 	fmt.Printf(buff.String())
	// case CODEC_JSON:
	// 	json, _ := e.Fields().Json()
	// 	fmt.Printf("%s\n", json)
	// 	break
	// case CODEC_PRETTYPRINT:
	// 	fallthrough
	// case CODEC_RUBYDEBUG:
	// 	pp.Printf("%s\n", e.Fields())
	// 	break
	// default:
	// 	p.Logger.Errorf("unknow codec %s", p.opt.Codec)
	// }

	p.Memory.Set("last", e.Fields().StringIndentNoTypeInfo(2))
	p.Send(e)
	return nil
}

func (p *processor) Start(e processors.IPacket) error {
	// p.WebHook.Add("events", p.HttpHandler)
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	return nil
}

// Handle Request received by bitfan for this agent (url hook should be registered during p.Start)
// func (p *processor) HttpHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	o := p.Memory.Items()
// 	for i, v := range o {
// 		// log.Printf("debug %s = %s", i, v)
// 		w.Write([]byte("<h3>" + i + "</h3><pre>" + v.(string) + "</pre>"))
// 	}
// }
