package testutils

import (
	"github.com/stretchr/testify/mock"
	"github.com/vjeantet/bitfan/processors"
	"testing"
	"github.com/stretchr/testify/assert"
)

type Processor struct {
	processors.Processor
	mock.Mock
}

func NewProcessor(f func() Processor) processors.Processor {
	return Processor{Processor: f()}
}

func AssertValuesForPaths(t *testing.T, ctx *DummyProcessorContext, pathValues map[string]interface{}) {
	for path, expectedVal := range pathValues {
		value, err := ctx.SentPackets(0)[0].Fields().ValueForPath(path)
		assert.Nil(t, err, "Unknown path: "+path)
		assert.Equal(t, expectedVal, value, "Invalid value for path "+path)
	}
}
