package fileinput

import (
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/mitchellh/mapstructure"
	"github.com/veino/field"
	"github.com/veino/veino"

	"github.com/hpcloud/tail"
	"github.com/hpcloud/tail/watch"
)

func New(l veino.Logger) veino.Processor {
	return &processor{Logger: l}
}

type processor struct {
	Logger    veino.Logger
	Send      veino.PacketSender
	NewPacket veino.PacketBuilder
	// TODO : mettre les options logstash
	filepath            string
	opt                 *options
	SinceDBInfos        map[string]*SinceDBInfo `json:"-"`
	sinceDBLastInfosRaw []byte                  `json:"-"`
	SinceDBLastSaveTime time.Time               `json:"-"`
	q                   chan bool
}

type options struct {
	Add_field              map[string]interface{}
	Close_older            int // 3600
	Codec                  string
	Delimiter              string // \n
	Discover_interval      int    // 15
	Exclude                []string
	Ignore_older           int // 86400
	Max_open_files         string
	Path                   []string `validate:"required"`
	Sincedb_path           string
	Sincedb_write_interval int    // 15
	Start_position         string // end
	Stat_interval          int    // 1
	Tags                   []string
	Type                   string
}

func (p *processor) Configure(conf map[string]interface{}) error {
	cf := options{
		Start_position:         "end",
		Sincedb_path:           ".sincedb.json",
		Sincedb_write_interval: 15,
		Stat_interval:          1,
	}

	if mapstructure.Decode(conf, &cf) != nil {
		return nil
	}
	p.opt = &cf
	p.filepath = p.opt.Path[0]

	return validator.New(&validator.Config{TagName: "validate"}).Struct(p.opt)
}
func (p *processor) Start(e veino.IPacket) error {
	watch.POLL_DURATION = time.Second * time.Duration(p.opt.Stat_interval)
	p.q = make(chan bool)
	p.tailFile(p.filepath, e, p.q)
	return nil
}

func (p *processor) Stop(e veino.IPacket) error {
	p.q <- true
	<-p.q
	return nil
}

func (p *processor) Tick(e veino.IPacket) error    { return nil }
func (p *processor) Receive(e veino.IPacket) error { return nil }

func (p *processor) tailFile(path string, packet veino.IPacket, q chan bool) {
	t, _ := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true, Poll: true})
	go func() {
		<-q
		go t.Stop()
		close(q)
	}()

	host, err := os.Hostname()
	if err != nil {
		p.Logger.Printf("can not get hostname : %s", err.Error())
	}

	go func() {
		for line := range t.Lines {

			e := p.NewPacket(line.Text, map[string]interface{}{
				"host":       host,
				"path":       path,
				"@timestamp": line.Time.Format(veino.VeinoTime),
			})

			field.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
			p.Send(e)

		}
	}()

	return
}
