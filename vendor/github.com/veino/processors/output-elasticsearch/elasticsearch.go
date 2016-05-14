package elasticsearch

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/veino/veino"
	"gopkg.in/olivere/elastic.v2"
)

var lines = map[string][]string{}

func New(l veino.Logger) veino.Processor {
	return &processor{}
}

type processor struct {
	client *elastic.Client
	opt    *options
}

type options struct {
	Host     string
	Cluster  string
	Protocol string
	Port     int
	User     string
	Password string
}

func (p *processor) Configure(conf map[string]interface{}) error {
	cf := options{Protocol: "http", Port: 9200}

	if err := mapstructure.Decode(conf, &cf); err != nil {
		return err
	}
	p.opt = &cf

	return nil
}

func (p *processor) Receive(e veino.IPacket) error {
	t := time.Now()
	index := fmt.Sprintf("logstash-%d.%02d.%02d", t.Year(), t.Month(), t.Day())
	// Add a document to the index
	data := e.Fields()
	_, err := p.client.Index().
		Index(index).
		Type("logs").
		BodyJson(data).
		Do()

	if err != nil {
		// Handle error
		panic(err)
	}

	return nil
}

func (p *processor) Tick(e veino.IPacket) error { return nil }

func (p *processor) Start(e veino.IPacket) error {
	var err error

	p.client, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("%s://%s:%d", p.opt.Protocol, p.opt.Host, p.opt.Port)),
		elastic.SetBasicAuth(p.opt.User, p.opt.Password),
		elastic.SetSniff(false),
	)

	if err != nil {
		// Handle error
		panic(err)
	}

	return nil
}

func (p *processor) Stop(e veino.IPacket) error { return nil }
