//go:generate bitfanDoc
//HTTPPoller allows you to intermittently poll remote HTTP URL, decode the output into an event
package httppoller

import (
	"bytes"
	"html/template"
	"io"

	"github.com/parnurzeal/gorequest"
	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/commons"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @Default "plain"
	// @Type codec
	Codec codecs.CodecCollection `mapstructure:"codec"`

	// Use CRON or BITFAN notation
	// @ExampleLS interval => "every_10s"
	Interval string `mapstructure:"interval"`

	// Http Method
	// @Default "GET"
	Method string `mapstructure:"method"`

	// Define headers for the request.
	// @ExampleLS headers => {"User-Agent":"Bitfan","Accept":"application/json"}
	Headers map[string]string `mapstructure:"headers"`

	// The request body (e.g. for an HTTP POST request). No default body is specified
	// @Type Location
	Body string `mapstructure:"body"`

	// URL
	// @ExampleLS url=> "http://google.fr"
	Url string `mapstructure:"url" validate:"required"`

	// When data is an array it stores the resulting data into the given target field.
	// When target is "" or "." it try to store retreived values at the root level of produced event
	// (usefull with json content -> codec)
	// @Default "output"
	Target string `mapstructure:"target"`

	// When true, unsuccessful HTTP requests, like unreachable connections, will
	// not raise an event, but a log message.
	// When false an event is generated with a tag _http_request_failure
	// @Default true
	IgnoreFailure bool `mapstructure:"ignore_failure"`

	// You can set variable to be used in Body by using ${var}.
	// each reference will be replaced by the value of the variable found in Body's content
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`
}

type processor struct {
	processors.Base
	q       chan bool
	opt     *options
	request *gorequest.SuperAgent
	BodyTpl *template.Template
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec: codecs.CodecCollection{
			Dec: codecs.New("plain", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
		Method:        "GET",
		Target:        "output",
		IgnoreFailure: true,
		Headers:       map[string]string{},
		Body:          "",
	}
	p.opt = &defaults

	err := p.ConfigureAndValidate(ctx, conf, p.opt)

	if p.opt.Body != "" {
		loc, err := commons.NewLocation(p.opt.Body, p.ConfigWorkingLocation)
		if err != nil {
			return err
		}
		content, _, err := loc.ContentWithOptions(p.opt.Var)
		if err != nil {
			return err
		}
		p.opt.Body = string(content)
		p.BodyTpl, err = template.New("body").Parse(p.opt.Body)
		if err != nil {
			p.Logger.Errorf("Body tpl error : %v", err)
			return err
		}
	}

	return err
}

func (p *processor) Start(e processors.IPacket) error {
	p.q = make(chan bool)
	p.request = gorequest.New()
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	close(p.q)
	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func (p *processor) Receive(e processors.IPacket) error {
	var (
		errs []error
		resp gorequest.Response
	)

	switch p.opt.Method {
	case "GET":
		sagent := p.request.Get(p.opt.Url)
		for k, v := range p.opt.Headers {
			sagent.Set(k, v)
		}

		if p.opt.Body != "" {
			buff := bytes.NewBufferString("")
			p.BodyTpl.Execute(buff, e.Fields())
			sagent.Send(buff.String())
		}

		resp, _, errs = sagent.End()
		e.Fields().SetValueForPath(p.request.Url, "httpRequestURL")
	default:
		p.Logger.Warnf("Method %s not implemented", p.opt.Method)
		return nil
	}

	if errs != nil {
		if p.opt.IgnoreFailure {
			for _, err := range errs {
				p.Logger.Warnf("while http requesting %s : %v", p.opt.Url, err)
			}
		} else {
			processors.AddTags([]string{"_http_request_failure"}, e.Fields())
			p.Send(e)
		}
		return nil
	}

	// Create a reader
	var dec codecs.Decoder
	var err error
	if dec, err = p.opt.Codec.NewDecoder(resp.Body); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		return nil
	}

	res := map[string]interface{}{}
	for i, h := range resp.Header {
		if len(h) > 0 {
			res[i] = h[0]
		}
	}
	res["status"] = resp.Status
	res["statusCode"] = resp.StatusCode
	res["proto"] = resp.Proto
	res["ContentLength"] = resp.ContentLength

	for dec.More() {
		var record interface{}
		if err = dec.Decode(&record); err != nil {
			if err == io.EOF {
				p.Logger.Debugln("error while http read docoding : ", err)
			} else {
				p.Logger.Errorln("error while http read docoding : ", err)
			}
			return nil
		}

		e2 := e.Clone()
		e2.Fields().SetValueForPath(res, "response")

		if p.opt.Target == "" || p.opt.Target == "." {

			switch v := record.(type) {
			case nil:
				break
			case string:
				e2.Fields().SetValueForPath(record, "output")
			case []interface{}:
				e2.Fields().SetValueForPath(record, "output")
			case map[string]interface{}:
				for k, v := range record.(map[string]interface{}) {
					e2.Fields().SetValueForPath(v, k)
				}
			default:
				p.Logger.Errorf("Unknow structure %#v", v)
			}

		} else {
			e2.Fields().SetValueForPath(record, p.opt.Target)
		}

		p.opt.ProcessCommonOptions(e2.Fields())
		p.Send(e2)
		select {
		case <-p.q:
			return nil
		default:
		}
	}

	return nil
}
