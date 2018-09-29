package rabbitmqoutput

import (
	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/bitfan/processors/doc"
	"github.com/vjeantet/bitfan/processors/testutils"
	"testing"
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
	assert.Equal(t, 0, max, "this processor does support concurency")
}
func TestConfigure(t *testing.T) {
	conf := map[string]interface{}{
		"exchange": "exchange",
		"exchange_type": "exchange_type",
	}
	ctx := testutils.NewProcessorContext()
	p := New()
	err := p.Configure(ctx, conf)
	assert.Nil(t, err, "Configure() processor without error")
}
