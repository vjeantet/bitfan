// mutate filter allows to perform general mutations on fields. You can rename, remove, replace, and modify fields in your event.
package mutate

import (
	"github.com/mitchellh/mapstructure"
	"github.com/veino/field"
	"github.com/veino/veino"
)

const (
	PORT_SUCCESS = 0
)

func New(l veino.Logger) veino.Processor {
	return &processor{}
}

type processor struct {
	Send veino.PacketSender

	// If this filter is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax.
	Add_tag []string

	// Convert a field’s value to a different type, like turning a string to an integer.
	// If the field value is an array, all members will be converted. If the field is a hash,
	// no action will be taken.
	// If the conversion type is boolean, the acceptable values are:
	// True: true, t, yes, y, and 1
	// False: false, f, no, n, and 0
	// If a value other than these is provided, it will pass straight through and log a warning message.
	// Valid conversion targets are: integer, float, string, and boolean.
	Convert map[string]string

	// Convert a string field by applying a regular expression and a replacement. If the field is not a string, no action will be taken.
	// This configuration takes an array consisting of 3 elements per field/substitution.
	// Be aware of escaping any backslash in the config file.
	Gsub []string

	// Join an array with a separator character. Does nothing on non-array fields
	Join map[string]string

	// Convert a value to its lowercase equivalent
	Lowercase []string

	// Merge two fields of arrays or hashes. String fields will be automatically be converted into an array
	Merge map[string]string

	// If this filter is successful, remove arbitrary fields from this event.
	Remove_field []string

	// If this filter is successful, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	Remove_tag []string

	// Rename key on one or more fields
	Rename map[string]string

	// Replace a field with a new value. The new value can include %{foo} strings to
	// help you build a new value from other parts of the event
	Replace map[string]interface{}

	// Split a field to an array using a separator character. Only works on string fields
	Split map[string]string

	// Strip whitespace from field. NOTE: this only works on leading and trailing whitespace
	Strip []string

	// Update an existing field with a new value. If the field does not exist, then no action will be taken
	Update map[string]interface{}

	// Convert a value to its uppercase equivalent
	Uppercase []string

	// remove all fields, except theses fields (work only with first level fields)
	Remove_all_but []string
}

func (p *processor) Configure(conf map[string]interface{}) error {
	if err := mapstructure.Decode(conf, p); err != nil {
		return err
	}

	return nil
}

func (p *processor) Receive(e veino.IPacket) error {
	field.AddFields(p.Add_field, e.Fields())
	field.AddTags(p.Add_tag, e.Fields())
	field.UpdateFields(p.Update, e.Fields())
	field.UpdateFields(p.Replace, e.Fields())
	field.RemoveFields(p.Remove_field, e.Fields())
	field.RenameFields(p.Rename, e.Fields())
	field.UpperCaseFields(p.Uppercase, e.Fields())
	field.LowerCaseFields(p.Lowercase, e.Fields())
	field.RemoveAllButFields(p.Remove_all_but, e.Fields())
	field.Convert(p.Convert, e.Fields())
	field.Join(p.Join, e.Fields())
	field.RemoveTags(p.Remove_tag, e.Fields())
	field.Gsub(p.Gsub, e.Fields())
	field.Split(p.Split, e.Fields())
	field.Strip(p.Strip, e.Fields())
	field.Merge(p.Merge, e.Fields())

	p.Send(e, PORT_SUCCESS)

	return nil
}

func (p *processor) Tick(e veino.IPacket) error { return nil }

func (p *processor) Start(e veino.IPacket) error { return nil }

func (p *processor) Stop(e veino.IPacket) error { return nil }
