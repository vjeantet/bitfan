package blacklist

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
	assert.EqualError(t, err, "Key: 'options.CompareField' Error:Field validation for 'CompareField' failed on the 'required' tag")
}

func TestConfigureNoTerms(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{
		"Compare_Field": "message",
		"terms":         []string{},
	}
	ctx := testutils.NewProcessorContext()
	err := p.Configure(ctx, conf)
	assert.EqualError(t, err, "blacklist option should have at least one value")
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
	if assert.Equal(t, 2, ctx.SentPacketsCount(0), "2 events should pass") {
		assert.Equal(t, "val1", ctx.SentPackets(0)[0].Message())
		assert.Equal(t, "val2", ctx.SentPackets(0)[1].Message())
	}
}

func TestReceiveDuplicateTermsInConfig(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"terms":         []string{"val1", "val1", "val1", "val1", "val1"},
		},
	)
	p.Receive(testutils.NewPacketOld("test", nil))
	p.Receive(testutils.NewPacketOld("fqsdf", nil))
	p.Receive(testutils.NewPacketOld("valo", nil))
	p.Receive(testutils.NewPacketOld("al1", nil))
	p.Receive(testutils.NewPacketOld("val1", nil))
	p.Receive(testutils.NewPacketOld("val3", nil))
	p.Receive(testutils.NewPacketOld("val2", nil))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "2 events should pass") {
		assert.Equal(t, "val1", ctx.SentPackets(0)[0].Message())
	}
}

func TestReceiveAllMessagesMatch(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"terms":         []string{"valo", "test", "val2", "al1", "val3", "val1", "fqsdf"},
		},
	)
	p.Receive(testutils.NewPacketOld("test", nil))
	p.Receive(testutils.NewPacketOld("fqsdf", nil))
	p.Receive(testutils.NewPacketOld("valo", nil))
	p.Receive(testutils.NewPacketOld("al1", nil))
	p.Receive(testutils.NewPacketOld("val1", nil))
	p.Receive(testutils.NewPacketOld("val3", nil))
	p.Receive(testutils.NewPacketOld("val2", nil))
	if assert.Equal(t, 7, ctx.SentPacketsCount(0), "2 events should pass") {
		assert.Equal(t, "test", ctx.SentPackets(0)[0].Message())
		assert.Equal(t, "fqsdf", ctx.SentPackets(0)[1].Message())
		assert.Equal(t, "valo", ctx.SentPackets(0)[2].Message())
		assert.Equal(t, "al1", ctx.SentPackets(0)[3].Message())
		assert.Equal(t, "val1", ctx.SentPackets(0)[4].Message())
		assert.Equal(t, "val3", ctx.SentPackets(0)[5].Message())
		assert.Equal(t, "val2", ctx.SentPackets(0)[6].Message())
	}
}

func TestReceiveMessageIncludingTermsDoNotMatch(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
			"terms":         []string{"test"},
		},
	)
	p.Receive(testutils.NewPacketOld("testtest", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "0 event should pass")
}

func TestReceiveFieldIncludingTermsDoNotMatch(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "MyField",
			"terms":         []string{"test"},
		},
	)

	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "testtest"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "0 event should pass")
}

func TestReceiveFieldNamesAreCaseSensitive(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "MyField",
			"terms":         []string{"test"},
		},
	)

	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"myfield": "test"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "myfield != MyField")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "test"}))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "match !")
}

func TestReceiveFieldValuesAreCaseSensitive(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "MyField",
			"terms":         []string{"test"},
		},
	)

	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "TEST"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "TEST != test")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "Test"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "TEST != test")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "test"}))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "match !")
}
func TestReceiveLongValue(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "MyField",
			"terms":         []string{"azertyuiopqsdfghjklmwxcvbnazertyuiopqsdfghjklmwxcvbnazertyuiopqsdfghjklmwxcvbnazertyuiopqsdfghjklmwxcvbn"},
		},
	)

	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "TEST"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "No match")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "Test"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "No match")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "azertyuiopqsdfghjklmwxcvbnazertyuiopqsdfghjklmwxcvbnazertyuiopqsdfghjklmwxcvbnazertyuiopqsdfghjklmwxcvbn"}))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "match !")
}
