//go:generate bitfanDoc
//HTTPPoller allows you to call an HTTP Endpoint, decode the output into an event
package httppoller

import (
	"io"

	"github.com/parnurzeal/gorequest"
	"github.com/vjeantet/bitfan/codecs"
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
	Codec codecs.Codec `mapstructure:"codec"`

	// Use CRON or BITFAN notation
	// @ExampleLS interval => "every_10s"
	Interval string `mapstructure:"interval"`

	// Http Method
	// @Default "GET"
	Method string `mapstructure:"method"`

	// URL
	// @ExampleLS url=> "http://google.fr"
	Url string `mapstructure:"url" validate:"required"`

	// When data is an array it stores the resulting data into the given target field.
	Target string `mapstructure:"target"`
}

type processor struct {
	processors.Base
	q       chan bool
	opt     *options
	request *gorequest.SuperAgent
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec:  codecs.New("plain", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		Method: "GET",
		Target: "output",
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
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
		resp, _, errs = p.request.Get(p.opt.Url).End()
	default:
		p.Logger.Warnf("Method %s not implemented", p.opt.Method)
		return nil
	}

	if errs != nil {
		p.Logger.Warnf("while http requesting %s : %#v", p.opt.Url, errs)
		return nil
	}
	if resp.StatusCode >= 400 {
		p.Logger.Warnf("http response code %s : %d (%s)", p.opt.Url, resp.StatusCode, resp.Status)
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
				p.Logger.Warnln("error while http read docoding : ", err)
			} else {
				p.Logger.Errorln("error while http read docoding : ", err)
				break
			}
		}

		e2 := e.Clone()
		e2.Fields().SetValueForPath(res, "response")
		e2.Fields().SetValueForPath(record, p.opt.Target)
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
