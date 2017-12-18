package digest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/processors/testutils"
	"time"
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
	assert.Equal(t, 1, max, "this processor does not support concurency")
}
func TestConfigureEmptyConfigIsInvalid(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{
	}
	ctx := testutils.NewProcessorContext()
	err := p.Configure(ctx, conf)
	assert.EqualError(t, err, "no interval and no Count settings set")
}

func TestReceiveSimple(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"key_map": "type",
		},
	)

	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"MyField": "TEST"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "No match")
	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"type": "a random value"}))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "One match")
	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"MyField": "azerty"}))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "No match")
	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"type": "a random value"}))
	assert.Equal(t, 2, ctx.SentPacketsCount(0), "Two match")
}

func TestReceiveMergeTwoEventsWithKeyMap(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"key_map": "type",
			"count":   2,
		},
	)

	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"type": "first_value", "key": "value1"}))
	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"type": "second_value", "key": "value2"}))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "Two match") {
		firstValue, err := ctx.SentPackets(0)[0].Fields().ValueForPath("first_value.key")
		assert.Nil(t, err, "No error")
		assert.Equal(t, "value1", firstValue)

		secondValue, err := ctx.SentPackets(0)[0].Fields().ValueForPath("second_value.key")
		assert.Nil(t, err, "No error")
		assert.Equal(t, "value2", secondValue)

	}
}
func TestReceiveMergeTwoEventsWithoutKeyMap(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"count": 6, // messages + fields
		},
	)

	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"key1": "value1", "key2": "value2"}))
	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"key3": "value3", "key4": "value4"}))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "Two match") {
		expected := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
			"key4": "value4",
		}
		AssertValuesForPaths(t, ctx, expected)
	}
}

func TestReceiveNoMatchBeforeCount(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"count": 100, // messages + fields
		},
	)

	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"key1": "value1", "key2": "value2"}))
	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"key3": "value3", "key4": "value4"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "Two match")
}

func TestReceiveSendEveryTwoSeconds(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			//"key_map": "",
			"count":    6,
			"interval": "every_1s",
		},
	)

	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"key1": "value1", "key2": "value2"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0))
	p.Receive(testutils.NewPacket("hello", map[string]interface{}{"key3": "value3", "key4": "value4"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0))
	time.Sleep(time.Second * 3)
	assert.Equal(t, 2, ctx.SentPacketsCount(0))

}

func AssertValuesForPaths(t *testing.T, ctx *testutils.DummyProcessorContext, pathValues map[string]string) {
	for path, expectedVal := range pathValues {
		value, err := ctx.SentPackets(0)[0].Fields().ValueForPath(path)
		assert.Nil(t, err, "Unknown path: "+path)
		assert.Equal(t, expectedVal, value, "Invalid value for path "+path)
	}
}
