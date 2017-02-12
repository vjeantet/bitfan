package whitelist

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/processors/testutils"
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

	p.Receive(testutils.NewPacket("test", nil))
	p.Receive(testutils.NewPacket("fqsdf", nil))
	p.Receive(testutils.NewPacket("valo", nil))
	p.Receive(testutils.NewPacket("al1", nil))
	p.Receive(testutils.NewPacket("val1", nil))
	p.Receive(testutils.NewPacket("val3", nil))
	p.Receive(testutils.NewPacket("val2", nil))
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

	p.Receive(testutils.NewPacket("val1", nil))
	p.Receive(testutils.NewPacket("val2", nil))
	p.Receive(testutils.NewPacket("val3", nil))
	p.Receive(testutils.NewPacket("val4", nil))
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

	p.Receive(testutils.NewPacket("val2d", nil))
	p.Receive(testutils.NewPacket("val1", nil))
	p.Receive(testutils.NewPacket("val1", nil))
	p.Receive(testutils.NewPacket("val1", nil))
	assert.Equal(t, 4, ctx.SentPacketsCount(0), "all events pass")
}
