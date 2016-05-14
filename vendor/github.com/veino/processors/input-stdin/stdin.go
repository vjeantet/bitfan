package stdin

import (
	"bufio"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/veino/field"
	"github.com/veino/veino"
)

func New(l veino.Logger) veino.Processor {
	return &processor{logger: l}
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	Tags []string

	// Add a type field to all events handled by this input
	Type string

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your veino pipeline
	Codec string
}

type processor struct {
	logger    veino.Logger
	NewPacket veino.PacketBuilder
	Send      veino.PacketSender
	opt       *options
	q         chan bool
}

func (p *processor) Configure(conf map[string]interface{}) error {
	cf := options{}
	if mapstructure.Decode(conf, &cf) != nil {
		return nil
	}
	p.opt = &cf
	return nil
}
func (p *processor) Start(e veino.IPacket) error {
	p.q = make(chan bool)

	stdinChan := make(chan string)
	go func(p *processor, ch chan string) {
		bio := bufio.NewReader(os.Stdin)
		for {

			line, hasMoreInLine, err := bio.ReadLine()
			if err == nil && hasMoreInLine == false {
				ch <- string(line)
			}
		}
	}(p, stdinChan)

	host, err := os.Hostname()
	if err != nil {
		p.logger.Printf("can not get hostname : %s", err.Error())
	}

	go func(ch chan string) {
		for {
			select {
			case stdin, _ := <-ch:

				ne := p.NewPacket(stdin, map[string]interface{}{
					"host": host,
				})

				field.ProcessCommonFields(ne.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
				p.Send(ne)

			case <-time.After(5 * time.Second):

			}

			select {
			case <-p.q:
				close(p.q)
				close(ch)
				return
			default:
			}
		}
	}(stdinChan)

	return nil
}

func (p *processor) Stop(e veino.IPacket) error {
	p.q <- true
	<-p.q
	return nil
}
func (p *processor) Tick(e veino.IPacket) error    { return nil }
func (p *processor) Receive(e veino.IPacket) error { return nil }
