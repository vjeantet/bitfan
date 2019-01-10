package processors

import (
	"testing"

	"github.com/awillis/bitfan/codecs"
	"github.com/awillis/bitfan/processors/doc"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/validator.v8"
)

func TestBaseNew(t *testing.T) {
	p := &Base{}
	p.SetPipelineUUID("123456")
	assert.Implements(t, (*Processor)(nil), p)
	assert.Equal(t, "123456", p.PipelineUUID)
	assert.IsType(t, (*doc.Processor)(nil), p.Doc())
	assert.ObjectsAreEqual(p, p.B())
}

func TestBaseMaxConcurentIsZero(t *testing.T) {
	p := &Base{}
	assert.Equal(t, 0, p.MaxConcurrent())
}

func TestBaseMethods(t *testing.T) {
	p := &Base{}
	ctx := newProcessorContext()
	conf := map[string]interface{}{}

	assert.NoError(t, p.Configure(ctx, conf))
	assert.NoError(t, p.Receive((IPacket)(nil)))
	assert.NoError(t, p.Tick((IPacket)(nil)))
	assert.NoError(t, p.Start((IPacket)(nil)))
	assert.NoError(t, p.Stop((IPacket)(nil)))
}

func TestBaseConfigureAndValidate(t *testing.T) {
	p := &Base{}
	ctx := newProcessorContext()

	rawVal := &struct {
		Numbeurre int
		Stringue  string `mapstructure:"driver"`
		Boule     bool
		Flaut     float64
	}{
		1,
		"hello",
		true,
		1.45,
	}

	conf := map[string]interface{}{
		"driver": "world",
	}

	err := p.ConfigureAndValidate(ctx, conf, rawVal)
	assert.NoError(t, err)
	assert.Equal(t, 1, rawVal.Numbeurre)
	assert.Equal(t, true, rawVal.Boule)
	assert.Equal(t, "world", rawVal.Stringue)
}

func TestBaseConfigureAndValidateValidationError(t *testing.T) {
	p := &Base{}
	ctx := newProcessorContext()

	rawVal := &struct {
		Numbeurre int
		Stringue  string `validate:"required"`
		Boule     bool
		Flaut     float64
	}{
		Boule: false,
	}

	conf := map[string]interface{}{
		"flaut": 1.67,
	}

	err := p.ConfigureAndValidate(ctx, conf, rawVal)

	assert.Error(t, err)
	assert.IsType(t, (validator.ValidationErrors)(nil), err)
	assert.Contains(t, err.(validator.ValidationErrors), ".Stringue")
}

func TestBaseConfigureAndValidateRawConfDecodingError(t *testing.T) {
	p := &Base{}
	ctx := newProcessorContext()

	rawVal := &struct {
		Numbeurre int
		Stringue  string `mapstructure:"driver"`
		Boule     bool
		Flaut     float64
	}{
		Boule: false,
	}

	conf := map[string]interface{}{
		"flaut":  "1.67",
		"driver": ctx,
	}

	err := p.ConfigureAndValidate(ctx, conf, rawVal)
	assert.Error(t, err)
	assert.IsType(t, (*mapstructure.Error)(nil), err)
}

func TestBaseConfigureAndValidateCodecs(t *testing.T) {
	p := &Base{}
	ctx := newProcessorContext()

	rawVal := &struct {
		Numbeurre int
		Stringue  string `mapstructure:"driver"`
		Boule     bool
		Flaut     float64
		Codec     codecs.CodecCollection
	}{
		Boule: false,
	}

	conf := map[string]interface{}{
		"flaut": "1.67",
		"codecs": map[int]interface{}{
			0: &dummyCodec{name: "codecType1"},
			1: &dummyCodec{name: "csv", role: "encoder"},
			2: &dummyCodec{name: "json", role: "decoder"},
		},
	}

	err := p.ConfigureAndValidate(ctx, conf, rawVal)
	assert.NoError(t, err)
	assert.Equal(t, "codecType1", rawVal.Codec.Default.Name)
	assert.Equal(t, "csv", rawVal.Codec.Enc.Name)
	assert.Equal(t, "json", rawVal.Codec.Dec.Name)

}

type dummyCodec struct {
	name string
	role string
}

func (d *dummyCodec) GetName() string {
	return d.name
}
func (d *dummyCodec) GetRole() string {
	return d.role
}
func (d *dummyCodec) GetOptions() map[string]interface{} {
	return map[string]interface{}{}
}

type dummyProcessorContext struct {
	logger        Logger
	packetSender  PacketSender
	packetBuilder PacketBuilder
	sentPackets   map[int][]IPacket
	builtPackets  []IPacket
	memory        Memory
	store         IStore
	mock.Mock
}

func newProcessorContext() *dummyProcessorContext {
	dp := &dummyProcessorContext{}
	dp.logger = logrus.New()
	dp.sentPackets = map[int][]IPacket{}
	dp.packetSender = nil
	dp.packetBuilder = nil
	dp.memory = nil
	dp.store = nil
	return dp
}

func (p dummyProcessorContext) Log() Logger {
	return p.logger
}
func (p dummyProcessorContext) PacketSender() PacketSender {
	return p.packetSender
}
func (p dummyProcessorContext) PacketBuilder() PacketBuilder {
	return p.packetBuilder
}
func (d *dummyProcessorContext) SentPackets(portNumber int) []IPacket {
	return d.sentPackets[portNumber]
}
func (d *dummyProcessorContext) BuiltPackets() []IPacket {
	return d.builtPackets
}
func (d *dummyProcessorContext) BuiltPacketsCount() int {
	return 0
}
func (d *dummyProcessorContext) SentPacketsCount(portNumber int) int {
	return 0
}
func (d *dummyProcessorContext) ConfigWorkingLocation() string {
	return ""
}
func (d *dummyProcessorContext) DataLocation() string {
	return ""
}
func (d *dummyProcessorContext) Memory() Memory {
	return d.memory
}
func (d *dummyProcessorContext) Store() IStore {
	return d.store
}
func (d *dummyProcessorContext) WebHook() WebHook {
	return nil
}
