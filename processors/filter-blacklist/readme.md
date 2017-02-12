# BLACKLIST
The blacklist rule will check a certain field against a blacklist, and match if it is in the blacklist.

## Synopsys


|    SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------------|--------|----------|---------------|
| add_field     | hash   | false    | {}            |
| add_tag       | array  | false    | []            |
| remove_field  | array  | false    | []            |
| remove_tag    | array  | false    | []            |
| compare_field | string | true     | ""            |
| blacklist     | array  | true     | []            |


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

### compare_field
* This is a required setting.
* Value type is string
* Default value is `""`

The name of the field to use to compare to the blacklist.
If the field is null, those events will be ignored.

### blacklist
* This is a required setting.
* Value type is array
* Default value is `[]`

A list of blacklisted values.
The compare_field term must be equal to one of these values for it to match.



## Configuration blueprint

```
blacklist{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	compare_field => "message"
	blacklist => ["val1","val2","val3"]
}
```
