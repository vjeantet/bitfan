//go:generate bitfanDoc

package httpoutput

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/awillis/bitfan/codecs"
	"github.com/awillis/bitfan/processors"
	"github.com/clbanning/mxj"
	"github.com/facebookgo/muster"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	httpClient *http.Client
	muster     muster.Client
	processors.Base
	enc      codecs.Encoder
	opt      *options
	shutdown chan struct{}
}

type options struct {
	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @Default "json"
	// @Enum "json","line","pp","rubydebug"
	// @Type codec
	Codec codecs.CodecCollection `mapstructure:"codec"`

	// Add a field to an event. Default value is {}
	AddField map[string]interface{} `mapstructure:"add_field"`

	// This output lets you send events to a generic HTTP(S) endpoint
	// This setting can be dynamic using the %{foo} syntax.
	URL string `mapstructure:"url" validate:"required"`

	// Custom headers to use format is headers => {"X-My-Header", "%{host}"}. Default value is {}
	// This setting can be dynamic using the %{foo} syntax.
	// @Default {"Content-Type" => "application/json"}
	Headers map[string]string `mapstructure:"headers"`

	// The HTTP Verb. One of "put", "post", "patch", "delete", "get", "head". Default value is "post"
	// @Default "post"
	HTTPMethod string `mapstructure:"http_method"`

	// Turn this on to enable HTTP keepalive support. Default value is true
	// @Default true
	KeepAlive bool `mapstructure:"keepalive"`

	// Max number of concurrent connections. Default value is 1
	// @Default 1
	PoolMax int `mapstructure:"pool_max"`

	// Timeout (in seconds) to wait for a connection to be established. Default value is 10
	// @Default 5
	ConnectTimeout uint `mapstructure:"connect_timeout"`

	// Timeout (in seconds) for the entire request. Default value is 60
	// @Default 30
	RequestTimeout uint `mapstructure:"request_timeout"`

	// If encountered as response codes this plugin will retry these requests
	// @Default [429, 500, 502, 503, 504]
	RetryableCodes []int `mapstructure:"retryable_codes"`

	// If you would like to consider some non-2xx codes to be successes
	// enumerate them here. Responses returning these codes will be considered successes
	IgnorableCodes []int `mapstructure:"ignorable_codes"`

	// @Default 5
	BatchInterval uint `mapstructure:"batch_interval"`

	// @Default 100
	BatchSize uint `mapstructure:"batch_size"`

	// @Default 10
	RetryInterval uint `mapstructure:"retry_interval"`

	// Add any number of arbitrary tags to your event. There is no default value for this setting.
	// This can help with processing later. Tags can be dynamic and include parts of the event using the %{field} syntax.
	// Tags []string `mapstructure:"tags"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		HTTPMethod:     "POST",
		KeepAlive:      true,
		PoolMax:        1,
		ConnectTimeout: 5,
		RequestTimeout: 30,
		Codec: codecs.CodecCollection{
			Enc: codecs.New("json", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		RetryableCodes: []int{429, 500, 502, 503, 504},
		IgnorableCodes: []int{},
		BatchInterval:  5,
		BatchSize:      100,
		RetryInterval:  10,
	}
	p.opt = &defaults
	p.opt.HTTPMethod = strings.ToUpper(p.opt.HTTPMethod)
	p.shutdown = make(chan struct{}, 1)
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	p.muster.Work <- e.Fields()
	return nil
}

func (p *processor) Start(e processors.IPacket) error {
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(p.opt.ConnectTimeout) * time.Second,
			KeepAlive: time.Duration(time.Second * 300),
		}).Dial,
		TLSClientConfig:       &tls.Config{},
		DisableCompression:    true,
		DisableKeepAlives:     !p.opt.KeepAlive,
		MaxIdleConns:          p.opt.PoolMax,
		MaxIdleConnsPerHost:   p.opt.PoolMax,
		ExpectContinueTimeout: time.Duration(time.Second * 3),
	}
	p.httpClient = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(p.opt.RequestTimeout) * time.Second,
	}
	p.muster.MaxBatchSize = p.opt.BatchSize
	p.muster.BatchTimeout = time.Duration(p.opt.BatchInterval) * time.Second
	p.muster.MaxConcurrentBatches = uint(p.opt.PoolMax)
	p.muster.PendingWorkCapacity = 0
	p.muster.BatchMaker = func() muster.Batch { return &batch{p: p} }
	err := p.muster.Start()
	return err
}

func (p *processor) Stop(e processors.IPacket) error {
	close(p.shutdown)
	return p.muster.Stop()
}

type batch struct {
	p       *processor
	Items   []*mxj.Map
	url     *string
	headers map[string]string
}

func (b *batch) Add(item interface{}) {
	fields := item.(*mxj.Map)
	if b.url == nil {
		url := b.p.opt.URL
		processors.Dynamic(&url, fields)
		b.url = &url
		b.headers = make(map[string]string)
		for k, v := range b.p.opt.Headers {
			processors.Dynamic(&k, fields)
			processors.Dynamic(&v, fields)
			b.headers[k] = v
		}
	}
	b.Items = append(b.Items, item.(*mxj.Map))
}

// Once a Batch is ready, it will be Fired. It must call notifier.Done once the
// batch has been processed.
func (b *batch) Fire(notifier muster.Notifier) {
	defer notifier.Done()
	var body bytes.Buffer
	writer := bufio.NewWriter(&body)
	enc, err := b.p.opt.Codec.NewEncoder(writer)
	if err != nil {
		b.p.Logger.Errorf("Lost %d messages. Codec failed with: %d", len(b.Items), err)
		return
	}
	for i := range b.Items {
		if err := enc.Encode(b.Items[i].Old()); err != nil {
			b.p.Logger.Errorf("Can't encode item with error: %v", err)
		}
	}
	if err := writer.Flush(); err != nil {
		b.p.Logger.Errorf("Lost %d messages with: %v", len(b.Items), err)
		return
	}

	for {
		retry, err := b.send(body.Bytes())
		if err == nil {
			b.p.Logger.Debugf("Successfully sent %d messages", len(b.Items))
			return
		}
		if !retry {
			b.p.Logger.Errorf("Lost %d messages. %v", len(b.Items), err)
			return
		}
		b.p.Logger.Warnf("Can't sent %d messages. %v. Retry after %d seconds", len(b.Items), err, b.p.opt.RetryInterval)
		select {
		case <-b.p.shutdown:
			b.p.Logger.Errorf("Shutdown. Lost %d messages", len(b.Items))
			return
		case <-time.NewTimer(time.Duration(b.p.opt.RetryInterval) * time.Second).C:
			continue
		}
	}
}

func (b *batch) send(body []byte) (retry bool, err error) {
	req, err := http.NewRequest(b.p.opt.HTTPMethod, *b.url, bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("Create request failed with: %v", err)
	}
	defer req.Body.Close()
	for hName, hValue := range b.headers {
		req.Header.Set(hName, hValue)
	}
	resp, err := b.p.httpClient.Do(req)
	if err != nil {
		return true, fmt.Errorf("Send request failed with: %v", err)
	}
	defer resp.Body.Close()

	io.Copy(ioutil.Discard, resp.Body)
	for _, ignoreCode := range b.p.opt.IgnorableCodes {
		if resp.StatusCode == ignoreCode {
			return false, nil
		}
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return false, nil
	}

	for _, retryCode := range b.p.opt.RetryableCodes {
		if resp.StatusCode == retryCode {
			return true, fmt.Errorf("Server returned %s", resp.Status)
		}
	}
	return false, fmt.Errorf("Server returned %s", resp.Status)
}
