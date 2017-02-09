//go:generate bitfanDoc
package elasticsearch

import (
	"fmt"
	"time"

	"github.com/vjeantet/bitfan/processors"
	"gopkg.in/olivere/elastic.v2"
)

var lines = map[string][]string{}

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

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

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt.Protocol = "http"
	p.opt.Port = 9200
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
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

func (p *processor) Start(e processors.IPacket) error {
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
