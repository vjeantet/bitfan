//go:generate bitfanDoc
// Read events from standard input.
// By default, each event is assumed to be one line. If you want to join lines, you’ll want to use the multiline filter.
package stdin

import (
	"bufio"
	"os"
	"time"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// Add a field to an event
	Add_field map[string]interface{}

	// Add any number of arbitrary tags to your event.
	// This can help with processing later.
	Tags []string

	// Add a type field to all events handled by this input
	Type string

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @default "line"
	Codec string
}

// Reads events from standard input
type processor struct {
	processors.Base

	opt *options
	q   chan bool
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}
func (p *processor) Start(e processors.IPacket) error {
	p.q = make(chan bool)

	stdinChan := make(chan string)
	go func(p *processor, ch chan string) {
		defer func() {
			if err := recover(); err != nil {

			}
		}()
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
		p.Logger.Warnf("can not get hostname : %s", err.Error())
	}

	go func(ch chan string) {
		for {
			select {
			case stdin, _ := <-ch:

				ne := p.NewPacket(stdin, map[string]interface{}{
					"host": host,
				})

				processors.ProcessCommonFields(ne.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
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

func (p *processor) Stop(e processors.IPacket) error {
	p.q <- true
	<-p.q
	return nil
}
