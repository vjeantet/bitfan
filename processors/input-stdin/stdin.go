//go:generate bitfanDoc
// Read events from standard input.
// By default, each event is assumed to be one line. If you want to join lines, you’ll want to use the multiline filter.
package stdin

import (
	"os"
	"time"

	"github.com/vjeantet/bitfan/codecs"
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
	Codec codecs.Codec
}

// Reads events from standard input
type processor struct {
	processors.Base

	opt  *options
	q    chan bool
	host string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec: codecs.New("line"),
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)

	if p.host, err = os.Hostname(); err != nil {
		p.Logger.Warnf("can not get hostname : %s", err.Error())
	}

	return err
}
func (p *processor) Start(e processors.IPacket) error {
	p.q = make(chan bool)

	var dec codecs.Decoder
	var err error

	if dec, err = p.opt.Codec.Decoder(os.Stdin); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		return err
	}

	stdinChan := make(chan string)
	go func(p *processor, ch chan string) {
		defer p.Logger.Errorln("XXXXXXXX")
		for {
			if record, err := dec.Decode(); err != nil {
				p.Logger.Errorln("codec error : ", err.Error())
				return
			} else {
				if record == nil {
					p.Logger.Debugln("waiting for more content...")
				} else {
					ch <- record["message"].(string)
				}
			}
		}
	}(p, stdinChan)

	go func(ch chan string) {
		for {
			select {
			case stdin, _ := <-ch:

				ne := p.NewPacket(stdin, map[string]interface{}{
					"host": p.host,
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
