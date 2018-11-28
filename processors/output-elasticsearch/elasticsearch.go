//go:generate bitfanDoc
package elasticsearch

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bitfan/processors"
	"github.com/jehiah/go-strftime"
	els5 "gopkg.in/olivere/elastic.v5"
	els6 "gopkg.in/olivere/elastic.v6"
)

var lines = map[string][]string{}

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	bulkProcessor6 *els6.BulkProcessor
	client6        *els6.Client

	bulkProcessor5 *els5.BulkProcessor
	client5        *els5.Client

	opt       *options
	lastIndex string
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
	// @Default true
	SSL bool `mapstructure:"ssl"`

	// ElasticSearch server version (6 or 5)
	// @Default 6
	Version int `mapstructure:"version"`
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
		SSL:           true,
		DocumentType:  "%{type}",
		Version:       6,
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.opt.Version > 6 || p.opt.Version < 5 {
		return fmt.Errorf("This processor support elasticsearch version 6 or 5, version %d is not supported", p.opt.Version)
	}

	return p.startBulkProcessor()
}

func (p *processor) Receive(e processors.IPacket) error {
	defer func() {
		if r := recover(); r != nil {
			p.Logger.Errorf("PANIC %s", r)
		}
	}()

	name := p.opt.Index
	processors.Dynamic(&name, e.Fields())

	// use @timestamp to compute index name, on error use time.Now()
	t, err := e.Fields().ValueForPath("@timestamp")
	if err != nil {
		t = time.Now()
	}
	index := strftime.Format(name, t.(time.Time))

	// Create Index if it does not exists
	p.checkIndex(index)

	// https://www.elastic.co/guide/en/logstash/current/plugins-outputs-elasticsearch.html#plugins-outputs-elasticsearch-document_type
	documentType := p.opt.DocumentType
	processors.Dynamic(&documentType, e.Fields())
	if documentType == "" {
		documentType = "logs"
	}

	switch p.opt.Version {
	case 6:
		event := els6.NewBulkIndexRequest().
			Index(index).
			Type(documentType).
			Doc(e.Fields().Old())
		p.bulkProcessor6.Add(event)
	case 5:
		event := els5.NewBulkIndexRequest().
			Index(index).
			Type(documentType).
			Doc(e.Fields().Old())
		p.bulkProcessor6.Add(event)
	}

	return nil
}

func (p *processor) startBulkProcessor() (err error) {
	scheme := map[bool]string{true: "https", false: "http"}[p.opt.SSL]

RECONNECT:
	switch p.opt.Version {
	case 6:
		p.client6, err = els6.NewClient(
			els6.SetURL(fmt.Sprintf("%s://%s:%d%s", scheme, p.opt.Host, p.opt.Port, p.opt.Path)),
			els6.SetBasicAuth(p.opt.User, p.opt.Password),
			els6.SetSniff(false),
		)
		if version, err := p.client6.ElasticsearchVersion(fmt.Sprintf("%s://%s:%d%s", scheme, p.opt.Host, p.opt.Port, p.opt.Path)); err == nil {
			if string(version[0]) != strconv.Itoa(p.opt.Version) {
				p.client6.Stop()
				if string(version[0]) == "5" {
					p.opt.Version = 5
					p.Logger.Infof("This elasticsearch server is v%s", version)
					goto RECONNECT
				}
				return fmt.Errorf("This elasticsearch server is v%s and does not match %d", version, p.opt.Version)
			}
		}
	case 5:
		p.client5, err = els5.NewClient(
			els5.SetURL(fmt.Sprintf("%s://%s:%d%s", scheme, p.opt.Host, p.opt.Port, p.opt.Path)),
			els5.SetBasicAuth(p.opt.User, p.opt.Password),
			els5.SetSniff(false),
		)
		if version, err := p.client5.ElasticsearchVersion(fmt.Sprintf("%s://%s:%d%s", scheme, p.opt.Host, p.opt.Port, p.opt.Path)); err == nil {
			if string(version[0]) != strconv.Itoa(p.opt.Version) {
				p.client5.Stop()
				if string(version[0]) == "6" {
					p.Logger.Infof("This elasticsearch server is v%s", version)
					p.opt.Version = 6
					goto RECONNECT
				}
				return fmt.Errorf("This elasticsearch server is v%s and does not match %d", version, p.opt.Version)
			}
		}
	}

	if err != nil {
		return err
	}
	fn := func(executionId int64, requests []els6.BulkableRequest, response *els6.BulkResponse, err error) {
		p.Logger.Debugf("commited %d requests ", len(requests))

	}

	switch p.opt.Version {
	case 6:
		p.bulkProcessor6, err = p.client6.BulkProcessor().
			BulkActions(p.opt.FlushCount).
			BulkSize(p.opt.FlushSize).
			FlushInterval(time.Duration(p.opt.IdleFlushTime) * time.Second).
			After(fn).
			Do(context.Background())
	case 5:
		p.bulkProcessor5, err = p.client5.BulkProcessor().
			BulkActions(p.opt.FlushCount).
			BulkSize(p.opt.FlushSize).
			FlushInterval(time.Duration(p.opt.IdleFlushTime) * time.Second).
			Do(context.Background())
	}

	return err
}

func (p *processor) checkIndex(name string) error {
	// alreadyseen index ?
	if p.lastIndex == name {
		return nil
	}
	// Check if the index exists
	switch p.opt.Version {
	case 6:
		if exists, err := p.client6.IndexExists(name).Do(context.Background()); err != nil {
			return err
		} else if !exists {
			p.client6.CreateIndex(name).Do(context.Background())
		}
	case 5:
		if exists, err := p.client5.IndexExists(name).Do(context.Background()); err != nil {
			return err
		} else if !exists {
			p.client5.CreateIndex(name).Do(context.Background())
		}
	}

	p.lastIndex = name
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	switch p.opt.Version {
	case 6:
		p.bulkProcessor6.Close()
	case 5:
		p.bulkProcessor5.Close()
	}

	return nil
}
