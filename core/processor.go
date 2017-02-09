package core

import "github.com/vjeantet/bitfan/processors"

type ProcessorFactory func() processors.Processor
