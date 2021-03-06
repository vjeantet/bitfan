// Code generated by "bitfanDoc "; DO NOT EDIT
package sleepprocessor

import "github.com/vjeantet/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
  Name:       "sleepprocessor",
  ImportPath: "github.com/vjeantet/bitfan/processors/sleep",
  Doc:        "Sleep a given amount of time.\n\nThis will cause bitfan to stall for the given amount of time.\n\nThis is useful for rate limiting, etc.",
  DocShort:   "",
  Options:    &doc.ProcessorOptions{
    Doc:     "",
    Options: []*doc.ProcessorOption{
      &doc.ProcessorOption{
        Name:           "processors.CommonOptions",
        Alias:          ",squash",
        Doc:            "",
        Required:       false,
        Type:           "processors.CommonOptions",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Time",
        Alias:          "",
        Doc:            "The length of time to sleep, in Millisecond, for every event.",
        Required:       false,
        Type:           "int",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
    },
  },
  Ports: []*doc.ProcessorPort{},
}
}