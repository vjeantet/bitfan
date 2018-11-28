package statsd

import (
	"testing"

	"bitfan/processors/doc"
	"bitfan/processors/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p := New()
	_, ok := p.(*processor)
	assert.Equal(t, ok, true, "New() should return a processor struct")
}
func TestDoc(t *testing.T) {
	assert.IsType(t, &doc.Processor{}, New().(*processor).Doc())
}

func TestMetricBuild(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	conf := map[string]interface{}{
		"sender": "%{message}",
	}
	assert.NoError(t, p.Configure(ctx, conf), "configuration is correct, error should be nil")
	assert.Equal(t, "200.response.total.200", p.dynamicKey("response.total.%{message}", testutils.NewPacketOld("200", nil)))
	assert.Equal(t, "400.response.total.400.100", p.dynamicKey("response.total.%{message}.%{int}", testutils.NewPacketOld("400", map[string]interface{}{"int": 100})))
	assert.Equal(t, "message.message.message", p.dynamicKey("%{message}.%{message}", testutils.NewPacketOld("message", nil)))

	v, err := p.dynamicValue("%{float}", testutils.NewPacketOld("message", map[string]interface{}{"float": 12.123}))
	assert.NoError(t, err)
	assert.Equal(t, 12.123, v)
	v, err = p.dynamicValue("%{int}", testutils.NewPacketOld("message", map[string]interface{}{"int": 123}))
	assert.NoError(t, err)
	assert.Equal(t, 123.0, v)
	v, err = p.dynamicValue("%{str}", testutils.NewPacketOld("message", map[string]interface{}{"str": "4444.99"}))
	assert.NoError(t, err)
	assert.Equal(t, 4444.99, v)
}
