package processors

import (
	"github.com/mitchellh/mapstructure"
	"bitfan/codecs"
	"bitfan/processors/doc"
	"gopkg.in/go-playground/validator.v8"
)

type Base struct {
	Send                  PacketSender
	NewPacket             PacketBuilder
	Logger                Logger
	Memory                Memory
	WebHook               WebHook
	Store                 IStore
	ConfigWorkingLocation string
	DataLocation          string
	PipelineUUID          string
	Name                  string
	Label                 string
}

// B returns the Base Processor
func (b *Base) B() *Base {
	return b
}

func (b *Base) Doc() *doc.Processor {
	return &doc.Processor{}
}

func (b *Base) SetPipelineUUID(uuid string) {
	b.PipelineUUID = uuid
}

func (b *Base) SetProcessorIdentifiers(name, label string) {
	b.Name = name
	b.Label = label
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

	b.Store = ctx.Store()

	//Codecs
	if v, ok := conf["codecs"]; ok {
		codecCollection := &codecs.CodecCollection{}
		for _, vc := range v.(map[int]interface{}) {
			vcodec := vc.(ICodec)

			c := codecs.New(vcodec.GetName(), vcodec.GetOptions(), ctx.Log(), ctx.ConfigWorkingLocation())
			switch vcodec.GetRole() {
			case "encoder":
				codecCollection.Enc = c
			case "decoder":
				codecCollection.Dec = c
			default:
				codecCollection.Default = c
			}

		}
		conf["codec"] = codecCollection
		delete(conf, "codecs")
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
