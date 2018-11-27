package digest

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"bitfan/processors/doc"
	"bitfan/processors/testutils"
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
	conf := map[string]interface{}{}
	ctx := testutils.NewProcessorContext()
	err := p.Configure(ctx, conf)
	assert.EqualError(t, err, "no interval and no Count settings set")
}

func TestConfigureNegativeCount(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{
		"count": -10,
	}
	ctx := testutils.NewProcessorContext()
	err := p.Configure(ctx, conf)
	assert.EqualError(t, err, "Negative count setting")
}

func TestConfigureInvalidCount(t *testing.T) {
	p := New().(*processor)
	conf := map[string]interface{}{
		"count": "hello",
	}
	ctx := testutils.NewProcessorContext()
	err := p.Configure(ctx, conf)
	assert.EqualError(t, err, "1 error(s) decoding:\n\n* cannot parse 'count' as int: strconv.ParseInt: parsing \"hello\": invalid syntax")
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

	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "TEST"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "No match")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"type": "a random value"}))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "One match")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"MyField": "azerty"}))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "No match")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"type": "a random value"}))
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

	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"type": "first_value", "key": "value1"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "Not enough packets")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"type": "second_value", "key": "value2"}))
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
			"count": 2,
		},
	)

	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"key1": "value1", "key2": "value2"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "Not enough packets")
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"key3": "value3", "key4": "value4"}))
	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "Two match") {
		expected := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
			"key4": "value4",
		}
		testutils.AssertValuesForPaths(t, ctx, expected)
	}
}

func TestReceiveNoMatchBeforeCount(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"count": 100,
		},
	)

	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"key1": "value1", "key2": "value2"}))
	p.Receive(testutils.NewPacketOld("hello", map[string]interface{}{"key3": "value3", "key4": "value4"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "Two match")
}

func TestReceiveTickEverySecond(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			//"key_map": "",
			"count":    2,
			"interval": "every_1s",
		},
	)

	// RECEIVE
	p.Receive(testutils.NewPacketOld("hello1", map[string]interface{}{"key1": "value1", "key2": "value2"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "No time elapsed and not enough packets")

	// TICK !
	time.Sleep(time.Second)
	p.Tick(testutils.NewPacketOld("", map[string]interface{}{}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "1 second elapsed but not enough packets")

	// RECEIVE
	p.Receive(testutils.NewPacketOld("hello2", map[string]interface{}{"key3": "value3", "key4": "value4"}))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "Enough packets but not enough time elapsed")

	// TICK !
	time.Sleep(time.Second)
	p.Tick(testutils.NewPacketOld("", map[string]interface{}{}))

	if assert.Equal(t, 1, ctx.SentPacketsCount(0), "Enough packets and enough time sleeping: Go !") {
		expected := map[string]interface{}{
			"message": "hello2",
			"key1":    "value1",
			"key2":    "value2",
			"key3":    "value3",
			"key4":    "value4",
		}
		testutils.AssertValuesForPaths(t, ctx, expected)
	}
}
