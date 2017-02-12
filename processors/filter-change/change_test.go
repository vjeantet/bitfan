package change

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

func TestReceiveDropAll(t *testing.T) {
	p := New().(*processor)
	ctx := testutils.NewProcessorContext()
	p.Configure(
		ctx,
		map[string]interface{}{
			"Compare_Field": "message",
		},
	)

	p.Receive(testutils.NewPacket("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 0")
	p.Receive(testutils.NewPacket("test", nil))
	p.Receive(testutils.NewPacket("test", nil))
	p.Receive(testutils.NewPacket("test", nil))
	p.Receive(testutils.NewPacket("test", nil))
	assert.Equal(t, 0, ctx.SentPacketsCount(0), "changed ! 1")
	p.Receive(testutils.NewPacket("toto", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 2")
	p.Receive(testutils.NewPacket("toto", nil))
	p.Receive(testutils.NewPacket("toto", nil))
	p.Receive(testutils.NewPacket("toto", nil))
	assert.Equal(t, 1, ctx.SentPacketsCount(0), "changed ! 3")
	p.Receive(testutils.NewPacket("test", nil))
	assert.Equal(t, 2, ctx.SentPacketsCount(0), "changed ! 4")

}
