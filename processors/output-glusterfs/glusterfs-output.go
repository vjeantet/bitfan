//go:generate bitfanDoc
// +build !extra

package glusterfsoutput

import (
	"fmt"

	"bitfan/processors"
)

func New() processors.Processor {
	return &processor{}
}

type processor struct {
	processors.Base
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	return fmt.Errorf("this plugin requires the tag `extra` to be passed at build time (e.g. `go build -tags extra`)")
}
