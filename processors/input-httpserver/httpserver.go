//go:generate bitfanDoc
// Listen and read a http request to build events with it.
//
// Processor respond with a HTTP code as :
//
// * `202` when request has been accepted, in body : the total number of event created
// * `500` when an error occurs, in body : an error description
//
// Use codecs to process body content as json / csv / lines / json lines / ....
//
// URL is available as http://webhookhost/pluginLabel/URI
//
// * webhookhost is defined by bitfan at startup
// * pluginLabel is defined in pipeline configuration, it's the named processor if you put one, or `input_httpserver` by default
// * URI is defined in plugin configuration (see below)
package httpserverprocessor

import (
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	uuid "github.com/nu7hatch/gouuid"
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
	//
	// Default decode http request as plain data, response is json encoded.
	// Set multiple codec with role to customize
	// @Default "plain"
	// @Type codec
	Codec codecs.CodecCollection

	// URI path
	// @Default "events"
	Uri string

	// Headers to send back into each outgoing response
	// @LSExample {"X-Processor" => "bitfan"}
	Headers map[string]string `mapstructure:"headers"`

	// What to send back to client ?
	// @Default ["uuid"]
	Body []string `mapstructure:"body"`
}

// Reads events from standard input
type processor struct {
	processors.Base

	opt  *options
	q    chan bool
	host string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Uri:  "events",
		Body: []string{"uuid"},
	}
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
	p.q = make(chan bool)
	p.WebHook.Add(p.opt.Uri, p.HttpHandler)
	return nil
}

// Handle Request received by bitfan for this agent (url hook should be registered during p.Start)
func (p *processor) HttpHandler(w http.ResponseWriter, r *http.Request) {
	p.Logger.Debug("reading request")

	// Create a reader
	var dec codecs.Decoder
	var err error

	if dec, err = p.opt.Codec.NewDecoder(r.Body); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
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

	var responseData map[string]interface{}
	responseData = map[string]interface{}{}

	p.Logger.Debug("request = ", req)
	p.Logger.Debug("start reading body content")
	i := 1
	for dec.More() {
		var record interface{}
		var body map[string]interface{}
		body = map[string]interface{}{}
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
			e = p.NewPacket("", map[string]interface{}{
				"request": req,
			})
		case string:
			e = p.NewPacket(v, map[string]interface{}{
				"request": req,
			})
		case map[string]interface{}:
			e = p.NewPacket("", v)
			e.Fields().SetValueForPath(req, "request")
		case []interface{}:
			e = p.NewPacket("", map[string]interface{}{
				"request": req,
				"data":    v,
			})
		default:
			p.Logger.Errorf("Unknow structure %#v", v)
		}

		id, _ := uuid.NewV4()
		e.Fields().SetValueForPath(id.String(), "uuid")
		p.opt.ProcessCommonOptions(e.Fields())

		for _, path := range p.opt.Body {
			value, err := e.Fields().ValueForPath(path)
			if err != nil {
				p.Logger.Errorf("ValueForPath %s - %v", path, err)
				continue
			}
			body[path] = value
		}

		responseData[strconv.Itoa(i)] = body

		p.Send(e)
		i = i + 1
		select {
		case <-p.q:
			return
		default:
		}
	}

	if err != nil && err != io.EOF {
		p.Logger.Errorln("error : ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(err.Error()))
		return
	}

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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for hn, hv := range p.opt.Headers {
		w.Header().Set(hn, hv)
		p.Logger.Debugf("added header : %s -> %s", hn, hv)
	}
	w.WriteHeader(http.StatusAccepted)
	enc.Encode(responseData)
}

func (p *processor) Stop(e processors.IPacket) error {
	close(p.q)
	return nil
}
