package processors

import (
	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/core/config"
	"github.com/vjeantet/bitfan/processors/doc"
	"gopkg.in/go-playground/validator.v8"
)

type Base struct {
	Send                  PacketSender
	NewPacket             PacketBuilder
	Logger                Logger
	Memory                Memory
	WebHook               WebHook
	ConfigWorkingLocation string
	DataLocation          string
	PipelineID            int
}

// B returns the Base Processor
func (b *Base) B() *Base {
	return b
}

func (b *Base) Doc() *doc.Processor {
	return &doc.Processor{}
}

func (b *Base) SetPipelineID(ID int) {
	b.PipelineID = ID
}

func (b *Base) Configure(ctx ProcessorContext, conf map[string]interface{}) error { return nil }

func (b *Base) Receive(e IPacket) error { return nil }

func (b *Base) Tick(e IPacket) error { return nil }

func (b *Base) Start(e IPacket) error { return nil }

func (b *Base) Stop(e IPacket) error { return nil }

func (b *Base) MaxConcurent() int {
	return 0
}

func (b *Base) ConfigureAndValidate(ctx ProcessorContext, conf map[string]interface{}, rawVal interface{}) error {

	// Logger
	b.Logger = ctx.Log()

	// Configuration location dir
	b.ConfigWorkingLocation = ctx.ConfigWorkingLocation()

	// Packet Sender func
	b.Send = ctx.PacketSender()

	// Packet Builder func
	b.NewPacket = ctx.PacketBuilder()

	// StoreSpace
	b.Memory = ctx.Memory()

	// WebHook
	b.WebHook = ctx.WebHook()

	// Datalocation
	b.DataLocation = ctx.DataLocation()

	//Codecs
	if v, ok := conf["codec"]; ok {
		switch vcodec := v.(type) {
		case *config.Codec:
			conf["codec"] = codecs.New(vcodec.Name, vcodec.Options, ctx.Log(), ctx.ConfigWorkingLocation())
		}
	}

	// Set processor's user options
	if err := mapstructure.WeakDecode(conf, rawVal); err != nil {
		return err
	}

	// validates processor's user options
	if err := validator.New(&validator.Config{TagName: "validate"}).Struct(rawVal); err != nil {
		return err
	}

	return nil
}
