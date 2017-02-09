# JSON


## Synopsys


|   SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|--------|----------|---------------|
| Add_field    | hash   | false    | {}            |
| Add_tag      | array  | false    | []            |
| Remove_field | array  | false    | []            |
| Remove_tag   | array  | false    | []            |
| Source       | string | false    | ""            |
| Target       | string | false    | ""            |


## Details

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

### Remove_tag
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary tags from the event.
Tags can be dynamic and include parts of the event using the %{field} syntax

### Source
* Value type is string
* Default value is `""`

The configuration for the JSON filter

### Target
* Value type is string
* Default value is `""`

Define the target field for placing the parsed data. If this setting is omitted,
the JSON data will be stored at the root (top level) of the event



## Configuration blueprint

```
json{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	source => ""
	target => ""
}
```
