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
	"fmt"
	"time"

	"github.com/k0kubun/pp"
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
	Codec string
}

// Prints events to the standard output
type processor struct {
	processors.Base

	// WebHook *core.WebHook
	opt *options
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	if err := p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}

	if p.opt.Codec == CODEC_RUBYDEBUG {
		p.opt.Codec = CODEC_PRETTYPRINT
	}

	if p.opt.Codec == "" {
		p.opt.Codec = CODEC_LINE
	}

	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	switch p.opt.Codec {
	case CODEC_LINE:
		t, _ := e.Fields().ValueForPath("@timestamp")
		fmt.Printf("%s %s %s\n",
			t.(time.Time).Format(timeFormat),
			e.Fields().ValueOrEmptyForPathString("host"),
			e.Message(),
		)
	case CODEC_JSON:
		json, _ := e.Fields().Json()
		fmt.Printf("%s\n", json)
		break
	case CODEC_PRETTYPRINT:
		pp.Printf("%s\n", e.Fields())
		break
	default:
		p.Logger.Errorf("unknow codec %s", p.opt.Codec)
	}

	p.Memory.Set("last", e.Fields().StringIndentNoTypeInfo(2))
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
