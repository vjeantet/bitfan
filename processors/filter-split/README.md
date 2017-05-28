# SPLIT
The split filter clones an event by splitting one of its fields and placing each value resulting from the split into a clone of the original event. The field being split can either be a string or an array.

An example use case of this filter is for taking output from the exec input plugin which emits one event for the whole output of a command and splitting that output by newline - making each line an event.

The end result of each split is a complete copy of the event with only the current split section of the given field changed.

## Synopsys


|   SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|--------|----------|---------------|
| Field        | string | false    | ""            |
| Target       | string | false    | ""            |
| Terminator   | string | false    | ""            |
| Add_field    | hash   | false    | {}            |
| Add_tag      | array  | false    | []            |
| Remove_field | array  | false    | []            |
| Remove_Tag   | array  | false    | []            |


## Details

### Field
* Value type is string
* Default value is `""`

The field which value is split by the terminator

### Target
* Value type is string
* Default value is `""`

The field within the new event which the value is split into. If not set, target field defaults to split field name

### Terminator
* Value type is string
* Default value is `""`

The string to split on. This is usually a line terminator, but can be any string
Default value is "\n"

### Add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.
Field names can be dynamic and include parts of the event using the %{field}.

### Add_tag
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event.
Tags can be dynamic and include parts of the event using the %{field} syntax.

### Remove_field
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary fields from this event.

### Remove_Tag
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary tags from the event.
Tags can be dynamic and include parts of the event using the %{field} syntax



## Configuration blueprint

```
split{
	field => ""
	target => ""
	terminator => ""
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
}
```
