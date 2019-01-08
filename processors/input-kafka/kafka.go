//go:generate bitfanDoc
package kafkainput

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/awillis/bitfan/processors"
)

func New() processors.Processor {
	return &processor{
		opt: &options{},
		wg:  new(sync.WaitGroup),
	}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// Bootstrap Server ( "host:port" )
	BootstrapServer string `mapstructure:"bootstrap_server"`
	// Broker list
	Brokers []string `mapstructure:"brokers"`
	// Kafka topic
	TopicID string `mapstructure:"topic_id" validate:"required"`
	// Kafka group id
	GroupID string `mapstructure:"group_id" validate:"required"`
	// Kafka client id
	ClientID string `mapstructure:"client_id"`
	// Queue Size
	QueueSize int `mapstructure:"queue_size"`
	// Minimum amount of bytes to fetch per request
	RequestBytesMin int `mapstructure:"request_bytes_min"`
	// Maximum amount of bytes to fetch per request
	RequestBytesMax int `mapstructure:"request_bytes_max"`
	// Keep Alive ( in seconds )
	KeepAlive int `mapstructure:"keepalive"`
	// Max time to wait for new data when fetching batches ( in seconds )
	MaxWait int `mapstructure:"max_wait"`
	// Frequency at which the reader lag is updated. Negative value disables lag reporting.
	ReadLagInterval int `mapstructure:"read_lag_interval"`
}

type processor struct {
	processors.Base

	opt    *options
	msgs   chan []byte
	wg     *sync.WaitGroup
	reader *kafka.Reader
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {

	defaults := options{
		Brokers:         []string{"localhost:9092"},
		ClientID:        "bitfan",
		GroupID:         "bitfan",
		QueueSize:       1024,
		RequestBytesMin: 10e3,
		RequestBytesMax: 10e6,
		KeepAlive:       180,
		MaxWait:         30,
		ReadLagInterval: 20,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {

	var err error

	// lookup bootstrap server
	if p.opt.BootstrapServer != "" {
		brokers, err := bootstrapLookup(p.opt.BootstrapServer)
		if err != nil {
			p.Logger.Errorf("error getting bootstrap servers: %v", err)
		} else {
			p.opt.Brokers = brokers
		}
	}

	p.Logger.Infof("using kafka brokers %v", p.opt.Brokers)

	p.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: p.opt.Brokers,
		GroupID: p.opt.GroupID,
		Topic:   p.opt.TopicID,
		Dialer: &kafka.Dialer{
			ClientID:  p.opt.ClientID,
			DualStack: true,
			KeepAlive: time.Second * time.Duration(p.opt.KeepAlive),
		},
		QueueCapacity:   p.opt.QueueSize,
		MinBytes:        p.opt.RequestBytesMin,
		MaxBytes:        p.opt.RequestBytesMax,
		MaxWait:         time.Second * time.Duration(p.opt.MaxWait),
		ReadLagInterval: time.Second * time.Duration(p.opt.ReadLagInterval),
		// the following options depend on group id
		//GroupBalancers:    []kafka.GroupBalancer{new(kafka.RangeGroupBalancer)},
		//HeartbeatInterval: 0,
		//CommitInterval:    0,
		//SessionTimeout:    0,
		//RebalanceTimeout:  0,
		//RetentionTime:     0,
	})

	go func(p *processor) {

		for {
			msg, err := p.reader.ReadMessage(context.Background())
			if err == io.EOF {
				break
			} else if err != nil {
				p.Logger.Errorf("error reading from kafka: %s", err)
				continue
			}

			var key string

			if len(msg.Key) > 0 {
				key = string(msg.Key)
			} else {
				key = "message"
			}

			ne := p.NewPacket(map[string]interface{}{
				key: msg.Value,
			})

			p.opt.ProcessCommonOptions(ne.Fields())
			p.Send(ne)
		}
	}(p)

	return err
}

func (p *processor) Stop(e processors.IPacket) error {
	return p.reader.Close()
}

func bootstrapLookup(endpoint string) ([]string, error) {

	var err error
	var brokers []string

	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return brokers, err
	}

	addrs, err := net.LookupHost(host)

	if err != nil {
		return brokers, err
	}

	for _, ip := range addrs {
		brokers = append(brokers, strings.Join([]string{ip, port}, ":"))
	}

	return brokers, err
}
