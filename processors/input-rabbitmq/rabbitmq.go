//go:generate bitfanDoc
package rabbitmqinput

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/clbanning/mxj"
	"github.com/streadway/amqp"
	"github.com/vjeantet/bitfan/processors"
)

// New returns a rabbimq processor
func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt  *options
	conn *amqp.Connection
	ch   *amqp.Channel
}

type options struct {
	// Enable message acknowledgements. Default value is true
	//
	// With acknowledgements messages fetched but not yet sent into the pipeline will be requeued by the server if BitFan shuts down.
	// Acknowledgements will however hurt the message throughput.
	// This will only send an ack back every prefetch_count messages. Working in batches provides a performance boost.
	Ack bool `mapstructure:"ack"`

	// Acknowledge messages in batch of value.
	// Default value is 1 (acknowledge each message individually)
	AckBatchSize uint64 `mapstructure:"ack_batch_size"`

	// Add a field to an event. Default value is {}
	AddField map[string]interface{} `mapstructure:"add_field"`

	// Extra queue arguments as an array. Default value is {}
	//
	// E.g. to make a RabbitMQ queue mirrored, use: {"x-ha-policy" => "all"}
	Arguments amqp.Table `mapstructure:"arguments"`

	// Should the queue be deleted on the broker when the last consumer disconnects? Default value is false
	//
	// Set this option to false if you want the queue to remain on the broker, queueing up messages until a consumer comes along to consume them.
	AutoDelete bool `mapstructure:"auto_delete"`

	// The codec used for input data. Default value is "json"
	//
	// Input codecs are a convenient method for decoding your data before it enters the input, without needing a separate filter in your BitFan pipeline.
	Codec string `mapstructure:"codec"`

	// Time in seconds to wait before retrying a connection. Default value is 1
	ConnectRetryInterval uint64 `mapstructure:"connect_retry_interval"`

	// Is this queue durable (a.k.a "Should it survive a broker restart?"")?  Default value is false
	Durable bool `mapstructure:"durable"`

	// The name of the exchange to bind the queue to. There is no default value for this setting.
	Exchange string `mapstructure:"exchange"`

	// Is the queue exclusive? Default value is false
	//
	//Exclusive queues can only be used by the connection that declared them and will be deleted when it is closed (e.g. due to a BitFan restart).
	Exclusive bool `mapstructure:"exclusive"`

	// Heartbeat delay in seconds. If unspecified no heartbeats will be sent
	Heartbeat int `mapstructure:"heartbeat"`

	// RabbitMQ server address. There is no default value for this setting.
	Host string `mapstructure:"host"`

	// The routing key to use when binding a queue to the exchange. Default value is ""
	//
	// This is only relevant for direct or topic exchanges.
	Key string `mapstructure:"key"`

	// Not implemented! Enable the storage of message headers and properties in @metadata. Default value is false
	//
	// This may impact performance
	MetadataEnabled bool `mapstructure:"metadata_enabled"`

	// Use queue passively declared, meaning it must already exist on the server. Default value is false
	//
	//
	// To have BitFan create the queue if necessary leave this option as false.
	// If actively declaring a queue that already exists, the queue options for this plugin (durable etc) must match those of the existing queue.
	Passive bool `mapstructure:"passive"`

	// RabbitMQ password. Default value is "guest"
	Password string `mapstructure:"password"`

	// RabbitMQ port to connect on. Default value is 5672
	Port int `mapstructure:"port"`

	// Prefetch count. Default value is 256
	//
	// If acknowledgements are enabled with the ack option, specifies the number of outstanding unacknowledged
	PrefetchCount int `mapstructure:"prefetch_count"`

	// The name of the queue BitFan will consume events from. If left empty, a transient queue with an randomly chosen name will be created.
	Queue string `mapstructure:"queue"`

	// Enable or disable SSL. Default value is false
	SSL bool `mapstructure:"ssl"`

	// Add any number of arbitrary tags to your event. There is no default value for this setting.
	//
	// This can help with processing later. Tags can be dynamic and include parts of the event using the %{field} syntax.
	Tags []string `mapstructure:"tags"`

	// RabbitMQ username. Default value is "guest"
	User string `mapstructure:"user"`

	// Validate SSL certificate. Default value is false
	VerifySSL bool `mapstructure:"verify_ssl"`

	// The vhost to use. Default value is "/"
	Vhost string `mapstructure:"vhost"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Ack:                  true,
		AckBatchSize:         1,
		AutoDelete:           false,
		ConnectRetryInterval: 1,
		Codec:                "json",
		Durable:              false,
		Exclusive:            false,
		MetadataEnabled:      false, // Not implemented
		Passive:              false,
		Password:             "guest",
		Port:                 5672,
		PrefetchCount:        256,
		SSL:                  false,
		User:                 "guest",
		VerifySSL:            false,
		Vhost:                "/",
	}

	p.opt = &defaults
	if err := p.ConfigureAndValidate(ctx, conf, p.opt); err != nil {
		return err
	}

	go func() {
		for {
			deliveries, err := p.consume()
			if err == nil {
				for msg := range deliveries {
					event := p.parse(msg.Body)
					processors.AddFields(p.opt.AddField, event.Fields())

					if len(p.opt.Tags) > 0 {
						processors.AddTags(p.opt.Tags, event.Fields())
					}

					if p.Send(event, 0) {
						if p.opt.Ack && (msg.DeliveryTag%p.opt.AckBatchSize) == 0 {
							p.ch.Ack(msg.DeliveryTag, true)
						}
					}
				}
			} else {
				p.Logger.Error(err)
			}
			time.Sleep(time.Duration(p.opt.ConnectRetryInterval) * time.Second)
		}
	}()

	return nil
}

func (p *processor) setup() (err error) {
	scheme := map[bool]string{true: "amqps", false: "amqp"}[p.opt.SSL]
	url := fmt.Sprintf("%s://%s:%s@%s:%d/%s", scheme, p.opt.User, p.opt.Password, p.opt.Host, p.opt.Port, p.opt.Vhost)

	p.Logger.Infoln("Connecting to " + url)

	amqpConfig := amqp.Config{Heartbeat: time.Duration(p.opt.Heartbeat) * time.Second}
	if p.opt.SSL {
		amqpConfig.TLSClientConfig = &tls.Config{InsecureSkipVerify: !p.opt.VerifySSL}
	}

	p.conn, err = amqp.DialConfig(url, amqpConfig)
	if err != nil {
		return err
	}

	p.ch, err = p.conn.Channel()
	if err != nil {
		return err
	}

	if !p.opt.Passive {
		_, err = p.ch.QueueDeclare(
			p.opt.Queue,
			p.opt.Durable,
			p.opt.AutoDelete,
			p.opt.Exclusive,
			false, // no-wait
			p.opt.Arguments,
		)
		if err != nil {
			return err
		}

		if err = p.ch.QueueBind(
			p.opt.Queue,
			p.opt.Key,
			p.opt.Exchange,
			false,
			nil,
		); err != nil {
			return err
		}
	}

	p.Logger.Infoln("Connected to " + url)
	return nil
}

func (p *processor) consume() (deliveries <-chan amqp.Delivery, err error) {
	if err = p.setup(); err != nil {
		return nil, err
	}

	if err := p.ch.Qos(p.opt.PrefetchCount, 0, false); err != nil {
		return nil, err
	}

	deliveries, err = p.ch.Consume(
		p.opt.Queue,
		"", // consumer
		!p.opt.Ack,
		p.opt.Exclusive,
		false, // no-local
		false, // no-wait
		p.opt.Arguments,
	)

	return deliveries, err
}

func (p *processor) parse(message []byte) processors.IPacket {
	var event processors.IPacket

	switch p.opt.Codec {
	case "json":
		fields, err := mxj.NewMapJson(message)
		if err != nil {
			p.Logger.Errorf(err.Error())
			event = p.NewPacket(string(message), nil)
		} else {
			event = p.NewPacket("", fields)
		}

	default:
		event = p.NewPacket(string(message), nil)
	}

	return event
}

func (p *processor) Stop(e processors.IPacket) error {
	p.ch.Close()
	p.conn.Close()
	return nil
}
