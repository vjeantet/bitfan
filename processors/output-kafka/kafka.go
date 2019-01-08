//go:generate bitfanDoc
package kafkaoutput

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/awillis/bitfan/processors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/gzip"
	"github.com/segmentio/kafka-go/lz4"
	"github.com/segmentio/kafka-go/snappy"
)

func New() processors.Processor {
	return &processor{
		msgs: make(chan []byte),
		wg:   new(sync.WaitGroup),
		opt:  &options{},
	}
}

type processor struct {
	processors.Base
	msgs   chan []byte
	wg     *sync.WaitGroup
	writer *kafka.Writer
	opt    *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// Bootstrap Server ( "host:port" )
	BootstrapServer string `mapstructure:"bootstrap_server"`
	// Broker list
	Brokers []string `mapstructure:"brokers"`
	// Kafka topic
	TopicID string `mapstructure:"topic_id" validate:"required"`
	// Kafka client id
	ClientID string `mapstructure:"client_id"`
	// Balancer ( roundrobin, hash or leastbytes )
	Balancer string `mapstructure:"balancer"`
	// Compression algorithm ( 'gzip', 'snappy', or 'lz4' )
	Compression string `mapstructure:"compression"`
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
	RequiredAcks int `mapstructure:"acks"`
	// Periodic Flush ( length of time in seconds a partially written buffer will live before being flushed )
	PeriodicFlush int `mapstructure:"pflush"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {

	defaults := options{
		Brokers:       []string{"localhost:9092"},
		ClientID:      "bitfan",
		MaxAttempts:   10,
		QueueSize:     1024,
		BatchSize:     256,
		KeepAlive:     180,
		IOTimeout:     30,
		RequiredAcks:  -1,
		PeriodicFlush: 20,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {

	var err error
	var balancer kafka.Balancer
	var codec kafka.CompressionCodec

	// lookup bootstrap server
	if p.opt.BootstrapServer != "" {
		brokers, err := bootstrapLookup(p.opt.BootstrapServer)
		if err != nil {
			p.Logger.Errorf("error getting bootstrap servers: %v", err)
		} else {
			p.opt.Brokers = brokers
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

	switch p.opt.Compression {
	case "gzip":
		codec = gzip.NewCompressionCodec()
	case "lz4":
		codec = lz4.NewCompressionCodec()
	case "snappy":
		codec = snappy.NewCompressionCodec()
	default:
		codec = nil
	}

	p.Logger.Infof("using kafka brokers %v", p.opt.Brokers)

	p.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: p.opt.Brokers,
		Topic:   p.opt.TopicID,
		Dialer: &kafka.Dialer{
			ClientID:  p.opt.ClientID,
			DualStack: true, // RFC-6555 compliance
			KeepAlive: time.Second * time.Duration(p.opt.KeepAlive),
		},
		Balancer:         balancer,
		CompressionCodec: codec,
		MaxAttempts:      p.opt.MaxAttempts,
		QueueCapacity:    p.opt.QueueSize,
		BatchSize:        p.opt.BatchSize,
		ReadTimeout:      time.Second * time.Duration(p.opt.IOTimeout),
		WriteTimeout:     time.Second * time.Duration(p.opt.IOTimeout),
		RequiredAcks:     p.opt.RequiredAcks,
		Async:            false,
	})

	go func(p *processor) {

		batch := make([]kafka.Message, 0, p.opt.BatchSize)
		var shutdown = false
		var pflush = false
		var pftimer = time.NewTimer(time.Second * time.Duration(p.opt.PeriodicFlush))
		p.wg.Add(1)

		for {
			select {
			case message, ok := <-p.msgs:
				if ok {
					batch = append(batch, kafka.Message{Value: message})
				} else {
					shutdown = true
				}
			case <-pftimer.C:
				pflush = true
			}

			if len(batch) == p.opt.BatchSize || shutdown == true || pflush == true {

				if !pftimer.Stop() {
					pflush = false
				}

				err = p.writer.WriteMessages(context.Background(), batch...)

				if err != nil {
					p.Logger.Errorf("error writing to kafka: %v", err)
				}

				if shutdown {
					p.Logger.Infof("shutting down kafka writer, flushed %d records", len(batch))
					p.wg.Done()
					break
				}

				pftimer.Reset(time.Second * time.Duration(p.opt.PeriodicFlush))
				batch = nil
			}
		}
	}(p)

	return err
}

func (p *processor) Receive(e processors.IPacket) error {

	var err error

	message, err := e.Fields().Json(true)

	if err != nil {
		p.Logger.Errorf("json encoding error: %v", err)
	}

	p.msgs <- message

	return err
}

func (p *processor) Stop(e processors.IPacket) error {
	close(p.msgs)
	p.wg.Wait()
	return p.writer.Close()
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
