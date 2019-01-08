package drop

import (
	"testing"

	"github.com/awillis/bitfan/processors/doc"
	"github.com/awillis/bitfan/processors/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor")
}
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}
func getExampleConfiguration() map[string]interface{} {
	return map[string]interface{}{
		"add_field": map[string]interface{}{
			"test1": "myvalue",
		},
	}
}

func TestConfigureDefaults(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{}
	ctx := testutils.NewProcessorContext()
	ret := p.Configure(ctx, conf)
	assert.Equal(t, nil, ret, "configuration is correct, it should return nil")
	assert.Equal(t, 100, p.opt.Percentage)
}

func TestConfigure(t *testing.T) {
	p := New().(*processor)
	conf := getExampleConfiguration()
	ctx := testutils.NewProcessorContext()
	ret := p.Configure(ctx, conf)
	assert.Equal(t, nil, ret, "configuration is correct, it should return nil")

	assert.Equal(t, len(p.opt.AddField), 1)
}

func TestReceiveDropAll(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(ctx, getExampleConfiguration())
	for i := 1; i <= 1000; i++ {
		em := testutils.NewPacketOld("", nil)
		p.Receive(em)
	}

	assert.Equal(t, 0, ctx.BuiltPacketsCount(), "unexpected event was created by the processor")
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "all events should be dropped")
}
func TestReceiveDrop99p(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(ctx, getExampleConfiguration())
	p.opt.Percentage = 99
	total := 1000
	for i := 1; i <= total; i++ {
		em := testutils.NewPacketOld("", nil)
		p.Receive(em)
	}

	expected := total / 100 * (100 - p.opt.Percentage)
	delta := total / 100
	assert.Equal(t, 0, ctx.BuiltPacketsCount(), "unexpected event was created by the processor")
	assert.InDelta(t, expected, ctx.SentPacketsCount(0), float64(delta))
	assert.Equal(t, "myvalue", ctx.SentPackets(0)[0].Fields().ValueOrEmptyForPathString("test1"), "a new should field be added")
}
func TestReceiveDrop80p(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(ctx, getExampleConfiguration())
	p.opt.Percentage = 80
	total := 10000
	for i := 1; i <= total; i++ {
		em := testutils.NewPacketOld("", nil)
		p.Receive(em)
	}

	expected := total / 100 * (100 - p.opt.Percentage)
	delta := total / 100
	assert.Equal(t, 0, ctx.BuiltPacketsCount(), "unexpected event was created by the processor")
	assert.InDelta(t, expected, ctx.SentPacketsCount(0), float64(delta))
	assert.Equal(t, "myvalue", ctx.SentPackets(0)[0].Fields().ValueOrEmptyForPathString("test1"), "a new should field be added")
}
