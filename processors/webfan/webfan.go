//go:generate bitfanDoc
package webfan

import (
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/entrypoint"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	//
	// Default decode http request as json, response is json encoded.
	// Set multiple codec with role to customize
	// @Default "json"
	// @Type codec
	Codec codecs.CodecCollection

	// URI path /_/path
	Uri string `mapstructure:"uri" validate:"required"`

	// Path to pipeline's configuration to execute on request
	// This configuration should contains only a filter section an a output like ```output{pass{}}```
	Conf string `mapstructure:"conf" validate:"required"`

	// Headers to send back into each outgoing response
	// @ExampleLS {"X-Processor" => "bitfan"}
	Headers map[string]string `mapstructure:"headers"`
}

// Reads events from standard input
type processor struct {
	processors.Base

	opt  *options
	wg   *sync.WaitGroup
	host string
	ep   *entrypoint.Entrypoint
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.host, err = os.Hostname(); err != nil {
		p.Logger.Warnf("can not get hostname : %v", err)
	}

	if p.opt.Codec.Enc == nil {
		p.opt.Codec.Enc = codecs.New("json", nil, ctx.Log(), ctx.ConfigWorkingLocation())
	}

	if p.opt.Codec.Dec == nil && p.opt.Codec.Default == nil {
		p.opt.Codec.Dec = codecs.New("plain", nil, ctx.Log(), ctx.ConfigWorkingLocation())
	}

	return err
}
func (p *processor) Start(e processors.IPacket) error {
	p.wg = &sync.WaitGroup{}
	p.WebHook.AddNamed(p.opt.Uri, p.HttpHandler)

	var err error
	p.ep, err = entrypoint.New(p.opt.Conf, p.ConfigWorkingLocation, entrypoint.CONTENT_REF)
	if err != nil {
		p.Logger.Errorf("Error with entrypoint %s", p.opt.Conf)
	}

	return err
}

// Handle Request received by bitfan for this agent (url hook should be registered during p.Start)
func (p *processor) HttpHandler(w http.ResponseWriter, r *http.Request) {
	p.wg.Add(1)
	defer p.wg.Done()
	p.Logger.Debug("reading request")

	// Build Pipeline
	ppl, err := p.ep.Pipeline()
	if err != nil {
		p.Logger.Errorf("%s", err)
		return
	}
	// pp.Println("ppl-->", ppl)
	orderedAgentConfList := core.Sort(ppl.Agents(), core.SortInputsFirst)

	// Find Last Agent
	lastAgent := orderedAgentConfList[len(orderedAgentConfList)-1]
	// Connect its Recipient to this Receive

	// Find First agent
	firstAgent := orderedAgentConfList[0]
	// pp.Println("lastAgent-->", lastAgent)
	back := make(chan processors.IPacket)
	lastAgent.Options["chan"] = back

	done := make(chan bool)
	go func(back chan processors.IPacket) {
		defer close(done)

		// Encode content
		var enc codecs.Encoder
		enc, err = p.opt.Codec.NewEncoder(w)
		if err != nil {
			p.Logger.Errorln("codec error : ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(err.Error()))
			return
		}

		firstPass := true
		for e := range back {
			if firstPass {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				for hn, hv := range p.opt.Headers {
					w.Header().Set(hn, hv)
					p.Logger.Debugf("added header : %s -> %s", hn, hv)
				}
				w.WriteHeader(http.StatusAccepted)
				firstPass = false
			}
			enc.Encode(e.Fields().Old())
		}

	}(back)

	_, err = ppl.Start()
	if err != nil {
		p.Logger.Errorf("Can not start webfan request pipeline : %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(err.Error()))
		close(back)
		return
	}

	// Create a reader
	var dec codecs.Decoder
	if dec, err = p.opt.Codec.NewDecoder(r.Body); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		close(back)
		close(back)
		return
	}
	headersBytes, _ := httputil.DumpRequest(r, false)
	headers := string(headersBytes)

	req := map[string]interface{}{
		"remoteAddr":  r.RemoteAddr,
		"rawHeaders":  headers,
		"method":      r.Method,
		"requestURI":  r.RequestURI,
		"proto":       r.Proto,
		"host":        r.Host,
		"requestPath": r.URL.Path,
	}

	req["querystring"] = map[string]interface{}{}
	for i, v := range r.URL.Query() {
		if len(v) == 1 {
			req["querystring"].(map[string]interface{})[i] = v[0]
		} else {
			req["querystring"].(map[string]interface{})[i] = v
		}
	}

	req["headers"] = map[string]interface{}{}
	for i, v := range r.Header {
		if len(v) == 1 {
			req["headers"].(map[string]interface{})[i] = v[0]
		} else {
			req["headers"].(map[string]interface{})[i] = v
		}
	}

	if r.Method == "POST" {
		r.ParseForm()
		req["formvalues"] = map[string]interface{}{}
		for i, v := range r.PostForm {
			if len(v) == 1 {
				req["formvalues"].(map[string]interface{})[i] = v[0]
			} else {
				req["formvalues"].(map[string]interface{})[i] = v
			}
		}
	}

	p.Logger.Debug("request = ", req)
	p.Logger.Debug("start reading body content")
	i := 1
	for dec.More() {
		var record interface{}

		if err = dec.Decode(&record); err != nil {
			if err == io.EOF {
				p.Logger.Debugln("error while http read docoding : ", err)
			} else {
				p.Logger.Errorln("error while http read docoding : ", err)
				break
			}
		}

		var e processors.IPacket
		switch v := record.(type) {
		case nil:
			e = p.NewPacket(map[string]interface{}{
				"request": req,
			})
		case string:
			e = p.NewPacket(map[string]interface{}{
				"message": v,
				"request": req,
			})
		case map[string]interface{}:
			e = p.NewPacket(v)
			e.Fields().SetValueForPath(req, "request")
		case []interface{}:
			e = p.NewPacket(map[string]interface{}{
				"request": req,
				"data":    v,
			})
		default:
			p.Logger.Errorf("Unknow structure %#v", v)
		}

		id, _ := uuid.NewV4()
		e.Fields().SetValueForPath(id.String(), "uuid")
		p.opt.ProcessCommonOptions(e.Fields())

		firstAgent.Processor().Receive(e)
		i = i + 1
	}

	if err != nil && err != io.EOF {
		p.Logger.Errorln("error : ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(err.Error()))
		close(back)
		return
	}

	ppl.Stop()
	<-done
}

func (p *processor) Stop(e processors.IPacket) error {
	p.wg.Wait()
	return nil
}
