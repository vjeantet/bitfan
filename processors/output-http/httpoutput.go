//go:generate bitfanDoc

package httpoutput

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/facebookgo/muster"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	httpClient *http.Client
	muster     muster.Client
	processors.Base
	opt      *options
	shutdown bool
}

type options struct {
	// Add a field to an event. Default value is {}
	AddField map[string]interface{} `mapstructure:"add_field"`

	// This output lets you send events to a generic HTTP(S) endpoint
	// This setting can be dynamic using the %{foo} syntax.
	URL string `mapstructure:"url" validate:"required"`

	// Custom headers to use format is headers => {"X-My-Header", "%{host}"}. Default value is {}
	// This setting can be dynamic using the %{foo} syntax.
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

	// Set the format of the http body. Now supports only "json_lines"
	// @Default "json_lines"
	Format string `mapstructure:"format"`

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

	// Add any number of arbitrary tags to your event. There is no default value for this setting.
	// This can help with processing later. Tags can be dynamic and include parts of the event using the %{field} syntax.
	// Tags []string `mapstructure:"tags"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		HTTPMethod:     "post",
		KeepAlive:      true,
		PoolMax:        1,
		ConnectTimeout: 5,
		RequestTimeout: 30,
		Format:         "json_lines",
		RetryableCodes: []int{429, 500, 502, 503, 504},
		IgnorableCodes: []int{},
		BatchInterval:  5,
		BatchSize:      100,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	// Convert dinamycs fields
	processors.Dynamic(&p.opt.URL, e.Fields())
	headers := make(map[string]string)
	for k, v := range p.opt.Headers {
		processors.Dynamic(&k, e.Fields())
		processors.Dynamic(&v, e.Fields())
		headers[k] = v
	}
	p.opt.Headers = headers

	var (
		eventBytes []byte
		err        error
	)
	switch p.opt.Format {
	case "json_lines":
		if eventBytes, err = e.Fields().Json(true); err != nil {
			return err
		}
		eventBytes = append(eventBytes, "\n"...)
	default:
		return fmt.Errorf("HTTP Output: invalid format '%s'", p.opt.Format)
	}

	p.muster.Work <- eventBytes
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
	p.shutdown = true
	return p.muster.Stop()
}

type batch struct {
	p     *processor
	Items bytes.Buffer
	size  uint
}

func (b *batch) Add(item interface{}) {
	b.Items.Write(item.([]byte))
	b.size = b.size + 1
}

// Once a Batch is ready, it will be Fired. It must call notifier.Done once the
// batch has been processed.
func (b *batch) Fire(notifier muster.Notifier) {
	defer notifier.Done()
	var (
		err  error
		req  *http.Request
		resp *http.Response
	)
	for {
		req, err = http.NewRequest(b.p.opt.HTTPMethod, b.p.opt.URL, &b.Items)
		if err != nil {
			b.p.Logger.Errorf("Create request failed with: %s\n", err.Error())
			return
		}
		for hName, hValue := range b.p.opt.Headers {
			req.Header.Set(hName, hValue)
		}
		for {
			if resp, err = b.p.httpClient.Do(req); err == nil {
				break
			}
			b.p.Logger.Error(err)
			time.Sleep(time.Second)
			if b.p.shutdown {
				return
			}
		}

		io.Copy(ioutil.Discard, resp.Body)

		for _, ignoreCode := range b.p.opt.IgnorableCodes {
			if resp.StatusCode == ignoreCode {
				b.p.Logger.Debugf("Successfully sent %d messages with status %s\n", b.size, resp.Status)
				resp.Body.Close()
				return
			}
		}
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			b.p.Logger.Debugf("Successfully sent %d messages with status %s\n", b.size, resp.Status)
			resp.Body.Close()
			return
		}

		retry := false
		for _, retryCode := range b.p.opt.RetryableCodes {
			if resp.StatusCode == retryCode {
				retry = true
				break
			}
		}
		if retry {
			b.p.Logger.Warnf("Server returned %s. Retry send\n", resp.Status)
			resp.Body.Close()
			req.Body.Close()
			time.Sleep(time.Second * 10)
			if b.p.shutdown {
				return
			}
			continue
		}
		b.p.Logger.Errorf("Server returned %s, %d messages was be lost\n", resp.Status, b.size)
		return
	}
}
