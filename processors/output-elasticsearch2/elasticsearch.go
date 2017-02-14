//go:generate bitfanDoc
package elasticsearch2

import (
	"fmt"
	"time"

	"github.com/jehiah/go-strftime"
	"github.com/vjeantet/bitfan/processors"
	"gopkg.in/olivere/elastic.v3"
)

var lines = map[string][]string{}

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	bulkProcessor *elastic.BulkProcessor
	client        *elastic.Client
	opt           *options
	lastIndex     string
}

type options struct {
	// The document type to write events to. There is no default value for this setting.
	//
	// Generally you should try to write only similar events to the same type.
	// String expansion %{foo} works here. Unless you set document_type, the event type will
	// be used if it exists otherwise the document type will be assigned the value of logs
	// @Default "%{type}"
	DocumentType string `mapstructure:"document_type"`

	// The number of requests that can be enqueued before flushing them. Default value is 1000
	// @Default 1000
	FlushCount int `mapstructure:"flush_count"`

	// The number of bytes that the bulk requests can take up before the bulk processor decides to flush. Default value is 5242880 (5MB).
	// @Default 5242880
	FlushSize int `mapstructure:"flush_size"`

	// Host of the remote instance. Default value is "localhost"
	// @Default "localhost"
	Host string `mapstructure:"host"`

	// The amount of seconds since last flush before a flush is forced. Default value is 1
	//
	// This setting helps ensure slow event rates don’t get stuck.
	// For example, if your flush_size is 100, and you have received 10 events,
	// and it has been more than idle_flush_time seconds since the last flush,
	// those 10 events will be flushed automatically.
	// This helps keep both fast and slow log streams moving along in near-real-time.
	// @Default 1
	IdleFlushTime int `mapstructure:"idle_flush_time"`

	// The index to write events to. Default value is "logstash-%Y.%m.%d"
	//
	// This can be dynamic using the %{foo} syntax and strftime syntax (see http://strftime.org/).
	// The default value will partition your indices by day.
	// @Default "logstash-%Y.%m.%d"
	Index string `mapstructure:"index"`

	// Password to authenticate to a secure Elasticsearch cluster. There is no default value for this setting.
	Password string `mapstructure:"password"`

	// HTTP Path at which the Elasticsearch server lives. Default value is "/"
	//
	// Use this if you must run Elasticsearch behind a proxy that remaps the root path for the Elasticsearch HTTP API lives.
	// @Default "/"
	Path string `mapstructure:"path"`

	// ElasticSearch port to connect on. Default value is 9200
	// @Default 9200
	Port int `mapstructure:"port"`

	// Username to authenticate to a secure Elasticsearch cluster. There is no default value for this setting.
	User string `mapstructure:"user"`

	// Enable SSL/TLS secured communication to Elasticsearch cluster. Default value is false
	// @Default false
	SSL bool `mapstructure:"ssl"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		FlushCount:    1000,
		FlushSize:     5242880,
		Host:          "localhost",
		IdleFlushTime: 1,
		Index:         "logstash-%Y.%m.%d",
		Path:          "/",
		Port:          9200,
		SSL:           false,
		DocumentType:  "%{type}",
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	return p.startBulkProcessor()
}

func (p *processor) Receive(e processors.IPacket) error {
	name := p.opt.Index
	processors.Dynamic(&name, e.Fields())

	// use @timestamp to compute index name, on error use time.Now()
	t, err := time.Parse(processors.TimeFormat, e.Fields().ValueOrEmptyForPathString("@timestamp"))
	if err != nil {
		t = time.Now()
	}
	index := strftime.Format(name, t)

	// Create Index if it does not exists
	p.checkIndex(index)

	// https://www.elastic.co/guide/en/logstash/current/plugins-outputs-elasticsearch.html#plugins-outputs-elasticsearch-document_type
	documentType := p.opt.DocumentType
	processors.Dynamic(&documentType, e.Fields())
	if documentType == "" {
		documentType = "logs"
	}

	event := elastic.NewBulkIndexRequest().
		Index(index).
		Type(documentType).
		Doc(e.Fields().Old())

	p.bulkProcessor.Add(event)
	return nil
}

func (p *processor) startBulkProcessor() (err error) {
	scheme := map[bool]string{true: "https", false: "http"}[p.opt.SSL]
	p.client, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("%s://%s:%d%s", scheme, p.opt.Host, p.opt.Port, p.opt.Path)),
		elastic.SetBasicAuth(p.opt.User, p.opt.Password),
		elastic.SetSniff(false),
	)

	if err != nil {
		return err
	}

	p.bulkProcessor, err = p.client.BulkProcessor().
		BulkActions(p.opt.FlushCount).
		BulkSize(p.opt.FlushSize).
		FlushInterval(time.Duration(p.opt.IdleFlushTime) * time.Second).
		Do()

	return err
}

func (p *processor) checkIndex(name string) error {
	// alreadyseen index ?
	if p.lastIndex == name {
		return nil
	}
	// Check if the index exists
	exists, err := p.client.IndexExists(name).Do()
	if err != nil {
		return err
	}
	if !exists {
		p.client.CreateIndex(name).Do()
	}
	p.lastIndex = name
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.bulkProcessor.Close()
	return nil
}
