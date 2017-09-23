//go:generate bitfanDoc
//HTTPPoller allows you to intermittently poll remote HTTP URL, decode the output into an event
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
	Codec codecs.CodecCollection `mapstructure:"codec"`

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

	// Level of failure
	//
	// 1 - noFailures
	// 2 - unsuccessful HTTP requests (unreachable connections)
	// 3 - unreachable connections and HTTP responses > 400 of successful HTTP requests
	// 4 - unreachable connections and non-2xx HTTP responses of successful HTTP requests
	// @default : 4
	FailureSeverity int `mapstructure:"failure_severity"`

	// When set, http failures will pass the received event and
	// append values to the tags field when there has been an failure
	// @ExampleLS tag_on_failure => ["_httprequestfailure"]
	// @default : []
	TagOnFailure []string `mapstructure:"tag_on_failure"`
}

const (
	failureSeverity_nothing          int = 1
	failureSeverity_unsuccessfulHTTP int = 2
	failureSeverity_HTTPover400      int = 3
	failureSeverity_HTTPnon2xx       int = 4
)

type processor struct {
	processors.Base
	q       chan bool
	opt     *options
	request *gorequest.SuperAgent
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec: codecs.CodecCollection{
			Dec: codecs.New("plain", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
		Method:          "GET",
		Target:          "output",
		FailureSeverity: failureSeverity_HTTPnon2xx,
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
		e.Fields().SetValueForPath(p.request.Url, "httpRequestURL")
	default:
		p.Logger.Warnf("Method %s not implemented", p.opt.Method)
		return nil
	}

	if errs != nil {
		p.Logger.Warnf("while http requesting %s : %#v", p.opt.Url, errs)

		if p.opt.FailureSeverity > failureSeverity_nothing {
			if len(p.opt.TagOnFailure) > 0 { // pass
				p.Logger.Debugf("network Failure pass event with tags %s", p.opt.TagOnFailure)
				processors.AddTags(p.opt.TagOnFailure, e.Fields())
				p.Send(e)
			}
		} else {
			p.Send(e)
		}
		return nil
	}

	if p.opt.FailureSeverity == failureSeverity_HTTPnon2xx {
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			if len(p.opt.TagOnFailure) > 0 { // pass
				processors.AddTags(p.opt.TagOnFailure, e.Fields())
			} else {
				p.Logger.Warnf("http response code %s : %d (%s)", p.opt.Url, resp.StatusCode, resp.Status)
				return nil
			}
		}
	} else if p.opt.FailureSeverity == failureSeverity_HTTPover400 {
		if resp.StatusCode >= 400 {
			if len(p.opt.TagOnFailure) > 0 { // pass
				processors.AddTags(p.opt.TagOnFailure, e.Fields())
			} else {
				p.Logger.Warnf("http response code %s : %d (%s)", p.opt.Url, resp.StatusCode, resp.Status)
				return nil
			}
		}
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
