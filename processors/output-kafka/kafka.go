//go:generate bitfanDoc
package kafkaoutput

import (
	"context"
	"github.com/miekg/dns"
	"github.com/segmentio/kafka-go"
	"github.com/vjeantet/bitfan/processors"
	"strconv"
	"strings"
	"time"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base
	writer *kafka.Writer
	opt    *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// Bootstrap Host
	BootstrapHost string `mapstructure:"bootstrap_host"`
	// Bootstrap Port
	BootstrapPort int `mapstructure:"bootstrap_port"`
	// Broker list
	Brokers []string `mapstructure:"brokers"`
	// Kafka topic
	Topic string `mapstructure:"topic" validate:"required"`
	// Balancer ( roundrobin, hash or leastbytes )
	Balancer string `mapstructure:"balancer"`
	// Max Attempts
	MaxAttempts int `mapstructure:"max_attempts"`
	// Queue Size
	QueueSize int `mapstructure:"queue_size"`
	// Batch Size
	BatchSize int `mapstructure:"batch_size"`
	// Keep Alive ( in seconds )
	KeepAlive int `mapstructure:"keepalive"`
	// IO Timeout ( in seconds )
	IOTimeout int `mapstructure:"io_timeout"`
	// Required Acks ( number of replicas that must acknowledge write. -1 for all replicas )
	RequiredAcks int `mapstructure:"required_acks"`
}

func bootstrapLookup(hostname string, port int) ([]string, error) {

	var result []string

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(hostname), dns.TypeA)

	dnsconf, err := dns.ClientConfigFromFile("/etc/resolv.conf")

	if err != nil {
		return result, err
	}

	ans, err := dns.Exchange(msg, strings.Join([]string{dnsconf.Servers[0], dnsconf.Port}, ":"))

	if err != nil {
		return result, err
	}

	for _, rr := range ans.Answer {
		if rr.Header().Rrtype == dns.TypeA {
			arec := rr.(*dns.A)
			broker := strings.Join([]string{arec.A.String(), strconv.Itoa(port)}, ":")
			result = append(result, broker)
		}
	}

	return result, err
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {

	defaults := options{
		Brokers:      []string{"localhost:9092"},
		MaxAttempts:  10,
		QueueSize:    10e3,
		BatchSize:    10e2,
		KeepAlive:    180,
		IOTimeout:    10,
		RequiredAcks: -1,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {

	var err error
	var balancer kafka.Balancer

	if p.opt.BootstrapHost != "" && p.opt.BootstrapPort > 0 {
		p.opt.Brokers, err = bootstrapLookup(p.opt.BootstrapHost, p.opt.BootstrapPort)
		if err != nil {
			p.Logger.Error(err)
		}
	}

	switch p.opt.Balancer {
	case "roundrobin":
		balancer = &kafka.RoundRobin{}
	case "hash":
		balancer = &kafka.Hash{}
	case "leastbytes":
		balancer = &kafka.LeastBytes{}
	default:
		balancer = &kafka.RoundRobin{}
	}

	p.Logger.Infof("using kafka brokers %v", p.opt.Brokers)

	p.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: p.opt.Brokers,
		Topic:   p.opt.Topic,
		Dialer: &kafka.Dialer{
			DualStack: true, // RFC-6555 compliance
			KeepAlive: time.Second * time.Duration(p.opt.KeepAlive),
		},
		Balancer:      balancer,
		MaxAttempts:   p.opt.MaxAttempts,
		QueueCapacity: p.opt.QueueSize,
		BatchSize:     p.opt.BatchSize,
		ReadTimeout:   time.Second * time.Duration(p.opt.IOTimeout),
		WriteTimeout:  time.Second * time.Duration(p.opt.IOTimeout),
		RequiredAcks:  p.opt.RequiredAcks,
		Async:         false,
	})

	return err
}

func (p *processor) Receive(e processors.IPacket) error {
	var err error

	err = p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("message"),
			Value: []byte(e.Message()),
		})

	return err
}

func (p *processor) Stop(e processors.IPacket) error {
	return p.writer.Close()
}
