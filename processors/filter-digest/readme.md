# DIGEST


## Synopsys


|   SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|--------|----------|---------------|
| add_field    | hash   | false    | {}            |
| add_tag      | array  | false    | []            |
| remove_field | array  | false    | []            |
| remove_tag   | array  | false    | []            |
| key_map      | string | false    | ""            |
| interval     | string | true     | ""            |


## Details

### add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### add_tag
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.

### remove_field
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary fields from this event. Example:
` kv {
`   remove_field => [ "foo_%{somefield}" ]
` }

### remove_tag
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary tags from the event. Tags can be dynamic and include parts of the event using the %{field} syntax.
Example:
` kv {
`   remove_tag => [ "foo_%{somefield}" ]
` }
If the event has field "somefield" == "hello" this filter, on success, would remove the tag foo_hello if it is present. The second example would remove a sad, unwanted tag as well.

### key_map
* Value type is string
* Default value is `""`

Add received event fields to the digest field named with the key map_key
When this setting is empty, digest will merge fields from coming events

### interval
* This is a required setting.
* Value type is string
* Default value is `""`

When should Digest send a digested event ?
Use CRON or BITFAN notation



## Configuration blueprint

```
digest{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	key_map => "type"
	interval => "every_10s"
}
```
