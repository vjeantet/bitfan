//go:generate bitfanDoc
package elasticinput

import (
	"fmt"
	"io"
	"time"

	elastic "gopkg.in/olivere/elastic.v3"

	"github.com/clbanning/mxj"
	"github.com/k0kubun/pp"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	Hosts    []string
	Query    string
	Size     int
	User     string
	Password string
}

type processor struct {
	processors.Base

	opt    *options
	q      chan bool
	client *elastic.Client
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Hosts: []string{"localhost"},
		Query: "",
		Size:  200,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {
	pp.Println(p.opt)
	p.q = make(chan bool)
	var err error
	// Create a client
	p.client, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s", p.opt.Hosts[0])),
		elastic.SetBasicAuth(p.opt.User, p.opt.Password),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}

	responseChan := make(chan interface{})
	go func(p *processor, ch chan interface{}) {

		query := elastic.NewQueryStringQuery(p.opt.Query)
		fmt.Println(p.opt.Query)
		scroll := p.client.Scroll().
			Index("logstash-*").
			Query(query).
			Size(p.opt.Size)

		searchResult, err := scroll.Do()
		if err == io.EOF {
			fmt.Print("Found no tweets\n")
		} else {
			fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)
			time.Sleep(time.Second)
		}

		for {

			// fmt.Println(hits, "/", searchResult.TotalHits())

			// Iterate through results
			for _, hit := range searchResult.Hits.Hits {
				// hit.Index contains the name of the index

				// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
				var e processors.IPacket
				fields, err := mxj.NewMapJson(*hit.Source)
				if err != nil {
					p.Logger.Error(err.Error())
					e = p.NewPacket(string(*hit.Source), nil)
				} else {
					e = p.NewPacket("", fields)
				}

				// if err != nil {
				// 	// 	// Deserialization failed
				// 	panic(err)
				// }

				// // Work with tweet
				// fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
				p.opt.ProcessCommonOptions(e.Fields())
				p.Send(e)
			}

			searchResult, err = scroll.Do()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
		}

	}(p, responseChan)

	// host, err := os.Hostname()
	// if err != nil {
	// 	p.Logger.Printf("can not get hostname : %s", err.Error())
	// }

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.q <- true
	<-p.q
	return nil
}
