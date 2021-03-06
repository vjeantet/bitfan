// Code generated by "bitfanDoc "; DO NOT EDIT
package digest

import "github.com/vjeantet/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
  Name:       "digest",
  ImportPath: "github.com/vjeantet/bitfan/processors/filter-digest",
  Doc:        "",
  DocShort:   "Digest events every x",
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
        Name:           "KeyMap",
        Alias:          "key_map",
        Doc:            "Add received event fields to the digest field named with the key map_key\nWhen this setting is empty, digest will merge fields from coming events",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "key_map => \"type\"",
      },
      &doc.ProcessorOption{
        Name:           "Interval",
        Alias:          "interval",
        Doc:            "When should Digest send a digested event ?\nUse CRON or BITFAN notation",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "interval => \"every_10s\"",
      },
      &doc.ProcessorOption{
        Name:           "Count",
        Alias:          "count",
        Doc:            "With min > 0, digest will not fire an event if less than min events were digested",
        Required:       false,
        Type:           "int",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
    },
  },
  Ports: []*doc.ProcessorPort{
    &doc.ProcessorPort{
      Default: true,
      Name:    "PORT_SUCCESS",
      Number:  0,
      Doc:     "",
    },
  },
}
}