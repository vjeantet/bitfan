package date

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"bitfan/processors/testutils"
)

type msg struct {
	Tags      []string `json:"tags"`
	Timestamp string   `json:"@timestamp"`
	TS        string   `json:"ts"`
}

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor")
}

func TestMatchUnix(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":    []string{"ts", "UNIX"},
		"timezone": "UTC",
	}
	p.Configure(ctx, conf)
	em := testutils.NewPacketOld("", map[string]interface{}{"ts": "1499254601"})
	p.Receive(em)
	var m msg
	em.Fields().Struct(&m)
	assert.NotContains(t, m.Tags, "_dateparsefailure")
	assert.Equal(t, "2017-07-05T11:36:41Z", m.Timestamp)
}

func TestMatchUnixWithMS(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":    []string{"ts", "UNIX"},
		"timezone": "UTC",
	}
	p.Configure(ctx, conf)
	em := testutils.NewPacketOld("", map[string]interface{}{"ts": "1499254601.343"})
	p.Receive(em)
	var m msg
	em.Fields().Struct(&m)
	assert.NotContains(t, m.Tags, "_dateparsefailure")
	assert.Equal(t, "2017-07-05T11:36:41.000000343Z", m.Timestamp)
}

func TestMatchUnixMS(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":    []string{"ts", "UNIX_MS"},
		"timezone": "UTC",
	}
	p.Configure(ctx, conf)
	em := testutils.NewPacketOld("", map[string]interface{}{"ts": "1499254601343"})
	p.Receive(em)
	var m msg
	em.Fields().Struct(&m)
	assert.NotContains(t, m.Tags, "_dateparsefailure")
	assert.Equal(t, "2017-07-05T11:36:41.000000343Z", m.Timestamp)
}

func TestMatchJODATime(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"match":    []string{"ts", "YYYY-MM-ddTHH:mm:ss"},
		"timezone": "Europe/Paris",
	}
	p.Configure(ctx, conf)
	em := testutils.NewPacketOld("", map[string]interface{}{"ts": "2017-07-05T11:36:41"})
	p.Receive(em)
	var m msg
	em.Fields().Struct(&m)
	assert.NotContains(t, m.Tags, "_dateparsefailure")
	assert.Equal(t, "2017-07-05T11:36:41+02:00", m.Timestamp)
}
