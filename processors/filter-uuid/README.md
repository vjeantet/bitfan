# UUID
The uuid filter allows you to generate a UUID and add it as a field to each processed event.

This is useful if you need to generate a string that’s unique for every event, even if the same input is processed multiple times. If you want to generate strings that are identical each time a event with a given content is processed (i.e. a hash) you should use the fingerprint filter instead.

The generated UUIDs follow the version 4 definition in RFC 4122).

## Synopsys


|   SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|--------|----------|---------------|
| Add_field    | hash   | false    | {}            |
| Add_tag      | array  | false    | []            |
| Remove_field | array  | false    | []            |
| Remove_Tag   | array  | false    | []            |
| Overwrite    | bool   | false    | ?             |
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

### Remove_Tag
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary tags from the event.
Tags can be dynamic and include parts of the event using the %{field} syntax

### Overwrite
* Value type is bool
* Default value is `?`

If the value in the field currently (if any) should be overridden by the generated UUID.
Defaults to false (i.e. if the field is present, with ANY value, it won’t be overridden)

### Target
* Value type is string
* Default value is `""`

Add a UUID to a field



## Configuration blueprint

```
uuid{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	overwrite => bool
	target => ""
}
```
