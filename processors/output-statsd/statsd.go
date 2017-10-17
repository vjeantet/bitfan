//go:generate bitfanDoc
package statsd

import (
	"fmt"
	"strings"

	"github.com/ShowMax/go-fqdn"
	"github.com/vjeantet/bitfan/processors"
	"gopkg.in/alexcesaro/statsd.v2"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt  *options
	conn *statsd.Client
}

type options struct {
	// The address of the statsd server.
	Host string `mapstructure:"host"`

	// The port to connect to on your statsd server.
	Port int `mapstructure:"port"`

	Protocol string `mapstructure:"protocol"`

	// The name of the sender. Dots will be replaced with underscores.
	Sender string `mapstructure:"sender"`

	// A count metric. metric_name => count as hash
	Count map[string]interface{} `mapstructure:"count"`

	// A decrement metric. Metric names as array.
	Decrement []string `mapstructure:"decrement"`

	// A gauge metric. metric_name => gauge as hash.
	Gauge map[string]interface{} `mapstructure:"gauge"`

	// An increment metric. Metric names as array.
	Increment []string `mapstructure:"increment"`

	// The statsd namespace to use for this metric.
	NameSpace string `mapstructure:"namespace"`

	// The sample rate for the metric.
	SampleRate float32 `mapstructure:"sample_rate"`

	// A set metric. metric_name => "string" to append as hash
	Set map[string]string `mapstructure:"set"`

	// A timing metric. metric_name => duration as hash
	Timing map[string]interface{} `mapstructure:"timing"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Host:       "localhost",
		Port:       8125,
		Protocol:   "udp",
		Sender:     fqdn.Get(),
		NameSpace:  "bitfan",
		SampleRate: 1.0,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	for key, value := range p.opt.Count {
		p.conn.Count(p.metricBuild(key, e), value)
	}
	for _, key := range p.opt.Increment {
		p.conn.Count(p.metricBuild(key, e), 1)
	}
	for _, key := range p.opt.Decrement {
		p.conn.Count(p.metricBuild(key, e), -1)
	}
	for key, value := range p.opt.Gauge {
		p.conn.Gauge(p.metricBuild(key, e) , value)
	}
	for key, value := range p.opt.Timing {
		p.conn.Timing(p.metricBuild(key, e), value)
	}
	for key, value := range p.opt.Set {
		p.conn.Unique(p.metricBuild(key, e), value)
	}
	return nil
}

func (p *processor) metricBuild(key string, e processors.IPacket) string {
	k, s := key, p.opt.Sender
	processors.Dynamic(&k, e.Fields())
	processors.Dynamic(&s, e.Fields())
	s = strings.Replace(s, ".", "_", -1)
	return fmt.Sprintf("%s.%s", s, k)
}

func (p *processor) Start(e processors.IPacket) error {
	url := fmt.Sprintf("%s:%d", p.opt.Host, p.opt.Port)
	var err error
	p.conn, err = statsd.New(
		statsd.Address(url),
		statsd.Prefix(p.opt.NameSpace),
		statsd.Network(p.opt.Protocol),
		statsd.SampleRate(p.opt.SampleRate),
		statsd.ErrorHandler(func(err error) {
			p.Logger.Error(err)
		}),
	)
	return err
}

func (p *processor) Stop(e processors.IPacket) error {
	p.conn.Close()
	return nil
}
