# CHANGE
This rule will monitor a certain field and match if that field changes. The field must change with respect to the last event

## Synopsys


|    SETTING     |  TYPE  | REQUIRED |  DEFAULT VALUE   |
|----------------|--------|----------|------------------|
| add_field      | hash   | false    | {}               |
| add_tag        | array  | false    | []               |
| remove_field   | array  | false    | []               |
| remove_tag     | array  | false    | []               |
| compare_field  | string | true     | ""               |
| ignore_missing | bool   | false    | true             |
| timeframe      | int    | false    | 0 (no timeframe) |


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

### ignore_missing
* Value type is bool
* Default value is `true`

If true, events without a compare_key field will not count as changed.

### timeframe
* Value type is int
* Default value is `0 (no timeframe)`

The maximum time in seconds between changes. After this time period, Bitfan will forget the old value of the compare_field field.



## Configuration blueprint

```
change{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	compare_field => "message"
	ignore_missing => true
	timeframe => 10
}
```
