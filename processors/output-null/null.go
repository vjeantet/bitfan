//go:generate bitfanDoc
// Drops everything received
package null

import "github.com/awillis/bitfan/processors"

func New() processors.Processor {
	return &processor{}
}

// Drops everything received
type processor struct {
	processors.Base
}
