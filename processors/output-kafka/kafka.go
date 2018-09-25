//go:generate bitfanDoc
package kafka


import (
	"github.com/vjeantet/bitfan/processors"
	"gopkg.in/mgo.v2"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt        *options
}

type options struct {


}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {
	var err error
	p.session, err = mgo.Dial(p.opt.Uri)
	if err != nil {
		return err
	}

	// Optional. Switch the session to a monotonic behavior.
	p.session.SetMode(mgo.Monotonic, true)
	p.collection = p.session.DB(p.opt.Database).C(p.opt.Collection)

	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	err := p.collection.Insert(e.Fields())

	return err
}

func (p *processor) Stop(e processors.IPacket) error {
	p.session.Close()
	return nil
}
