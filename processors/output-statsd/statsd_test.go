package statsd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/processors/testutils"
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
	assert.Equal(t, "200.response.total.200", p.dynamicKey("response.total.%{message}", testutils.NewPacket("200", nil)))
	assert.Equal(t, "400.response.total.400", p.dynamicKey("response.total.%{message}", testutils.NewPacket("400", nil)))
	assert.Equal(t, "message.message.message", p.dynamicKey("%{message}.%{message}", testutils.NewPacket("message", nil)))
}
