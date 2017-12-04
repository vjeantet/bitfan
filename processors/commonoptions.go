package processors

import "github.com/clbanning/mxj"

type CommonOptions struct {
	// If this filter is successful, add any arbitrary fields to this event.
	AddField map[string]interface{} `mapstructure:"add_field"`

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	AddTag []string `mapstructure:"add_tag"`

	// Add a type field to all events handled by this input
	Type string `mapstructure:"type"`

	// If this filter is successful, remove arbitrary fields from this event. Example:
	// ` kv {
	// `   remove_field => [ "foo_%{somefield}" ]
	// ` }
	RemoveField []string `mapstructure:"remove_field"`

	// If this filter is successful, remove arbitrary tags from the event. Tags can be dynamic and include parts of the event using the %{field} syntax.
	// Example:
	// ` kv {
	// `   remove_tag => [ "foo_%{somefield}" ]
	// ` }
	// If the event has field "somefield" == "hello" this filter, on success, would remove the tag foo_hello if it is present. The second example would remove a sad, unwanted tag as well.
	RemoveTag []string `mapstructure:"remove_tag"`

	// Log each event produced by the processor (usefull while building or debugging a pipeline)
	Trace bool `mapstructure:"trace"`
}

func (c *CommonOptions) ProcessCommonOptions(data *mxj.Map) {
	if len(c.AddField) > 0 {
		AddFields(c.AddField, data)
	}
	if len(c.AddTag) > 0 {
		AddTags(c.AddTag, data)
	}
	if c.Type != "" {
		SetType(c.Type, data)
	}
	if len(c.RemoveField) > 0 {
		RemoveFields(c.RemoveField, data)
	}
	if len(c.RemoveTag) > 0 {
		RemoveTags(c.RemoveTag, data)
	}
}
