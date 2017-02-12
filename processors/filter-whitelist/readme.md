# WHITELIST
Similar to blacklist, this processor will compare a certain field to a whitelist, and match
if the list does not contain the term

## Synopsys


|    SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------------|--------|----------|---------------|
| add_field     | hash   | false    | {}            |
| add_tag       | array  | false    | []            |
| remove_field  | array  | false    | []            |
| remove_tag    | array  | false    | []            |
| compare_field | string | true     | ""            |
| ignore_null   | bool   | false    | true          |
| list          | array  | true     | []            |


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

The name of the field to use to compare to the whitelist.
If the field is null, those events will be ignored.

### ignore_null
* Value type is bool
* Default value is `true`

If true, events without a compare_key field will not match.

### list
* This is a required setting.
* Value type is array
* Default value is `[]`

A list of whitelisted values.
The compare_field term must be in this list or else it will match.



## Configuration blueprint

```
whitelist{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	compare_field => "message"
	ignore_null => true
	list => ["val1","val2","val3"]
}
```
