package whitelist

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
func TestMaxConcurent(t *testing.T) {
	max := New().(*processor).MaxConcurent()
	assert.Equal(t, 0, max, "this processor support concurency with no limit")
}

func TestConfigureNoCompareFields(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{
		"terms": []string{"val"},
	}
	ctx := testutils.NewProcessorContext()
	err := p.Configure(ctx, conf)
	assert.Error(t, err, "configuration is not correct")
}

func TestConfigureNoTerms(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{
		"Compare_Field": "message",
		"terms":         []string{},
	}
	ctx := testutils.NewProcessorContext()
	err := p.Configure(ctx, conf)
	assert.Error(t, err, "configuration is not correct")
}

func TestReceiveMatch(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"terms":         []string{"val1", "val2"},
		},
	)

	p.Receive(testutils.NewPacketOld("test", nil))
	p.Receive(testutils.NewPacketOld("fqsdf", nil))
	p.Receive(testutils.NewPacketOld("valo", nil))
	p.Receive(testutils.NewPacketOld("al1", nil))
	p.Receive(testutils.NewPacketOld("val1", nil))
	p.Receive(testutils.NewPacketOld("val3", nil))
	p.Receive(testutils.NewPacketOld("val2", nil))
	assert.Equal(t, 5, ctx.SentPacketsCount(0), "5 events should pass")
}

func TestReceiveMissingFieldIgnoreTrue(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "toto",
			// "Ignore_Missing": true, // default value
			"terms": []string{"val1", "val2"},
		},
	)

	p.Receive(testutils.NewPacketOld("val1", nil))
	p.Receive(testutils.NewPacketOld("val2", nil))
	p.Receive(testutils.NewPacketOld("val3", nil))
	p.Receive(testutils.NewPacketOld("val4", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "no event pass")
}

func TestReceiveMissingFieldIgnoreFalse(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field":  "toto",
			"Ignore_Missing": false,
			"terms":          []string{"val1", "val2"},
		},
	)

	p.Receive(testutils.NewPacketOld("val2d", nil))
	p.Receive(testutils.NewPacketOld("val1", nil))
	p.Receive(testutils.NewPacketOld("val1", nil))
	p.Receive(testutils.NewPacketOld("val1", nil))
	assert.Equal(t, 4, ctx.SentPacketsCount(0), "all events pass")
}
