//go:generate bitfanDoc
// Read events from standard input.
// By default, each event is assumed to be one line. If you want to join lines, you’ll want to use the multiline filter.
package stdin

import (
	"io"
	"os"
	"time"

	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/core"
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
	// @Default "line"
	// @Type codec
	Codec codecs.Codec

	// Stop bitfan on stdin EOF ? (use it when you pipe data with |)
	// @Default false
	EofExit bool `mapstructure:"eof_exit"`
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
		Codec:   codecs.New("line", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		EofExit: false,
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.host, err = os.Hostname(); err != nil {
		p.Logger.Warnf("can not get hostname : %s", err.Error())
	}

	return err
}
func (p *processor) Start(e processors.IPacket) error {
	p.q = make(chan bool)

	var dec codecs.Decoder
	var err error

	if dec, err = p.opt.Codec.NewDecoder(os.Stdin); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		return err
	}

	stdinChan := make(chan interface{})
	go func(p *processor, ch chan interface{}) {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				p.Logger.Errorf("Panic ! stdin - %s", err.Error())
			}
		}()

		for {
			var record interface{}
			if err := dec.Decode(&record); err != nil {
				if err == io.EOF {
					p.Logger.Debugf("codec end of file", err.Error())
					if p.opt.EofExit {
						core.Stop()
						p, _ := os.FindProcess(os.Getpid())
						p.Signal(os.Interrupt)
					}

				} else {
					p.Logger.Errorln("codec error : ", err.Error())
				}
				return
			} else {
				ch <- record
			}
		}
	}(p, stdinChan)

	go func(ch chan interface{}) {
		for {
			select {
			case msg := <-ch:
				var ne processors.IPacket

				switch v := msg.(type) {
				case string:
					ne = p.NewPacket(v, map[string]interface{}{
						"host": p.host,
					})
				case map[string]interface{}:
					ne = p.NewPacket("", v)
					ne.Fields().SetValueForPath(p.host, "host")
				case []interface{}:
					ne = p.NewPacket("", map[string]interface{}{
						"host": p.host,
						"data": v,
					})
				default:
					p.Logger.Errorf("Unknow structure %#v", v)
				}

				processors.ProcessCommonFields(ne.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
				p.Send(ne)

			case <-time.After(1 * time.Second):

			case <-p.q:
				close(p.q)
				close(ch)
				return
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
