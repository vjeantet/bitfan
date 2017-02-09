package testutils

import (
	"github.com/stretchr/testify/mock"
	"github.com/vjeantet/bitfan/processors"
)

type Processor struct {
	processors.Processor
	mock.Mock
}

func NewProcessor(f func() Processor) processors.Processor {
	return Processor{Processor: f()}
}
