# EXEC
Execute a command and use its stdout as event data

## Synopsys


|   SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|--------|----------|---------------|
| add_field    | hash   | false    | {}            |
| add_tag      | array  | false    | []            |
| remove_field | array  | false    | []            |
| remove_tag   | array  | false    | []            |
| command      | string | true     | ""            |
| args         | array  | false    | []            |
| stdin        | bool   | false    | false         |
| target       | string | false    | "stdout"      |
| codec        | codec  | false    | "plain"       |


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

### command
* This is a required setting.
* Value type is string
* Default value is `""`



### args
* Value type is array
* Default value is `[]`



### stdin
* Value type is bool
* Default value is `false`

Pass the complete event to stdin as a json string

### target
* Value type is string
* Default value is `"stdout"`

Where do the output should be stored
Set "." when output is json formated and want to replace current event fields with output
response. (useful)

### codec
* Value type is codec
* Default value is `"plain"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline



## Configuration blueprint

```
exec{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	command => ""
	args => []
	stdin => false
	target => "stdout"
	codec => "plain"
}
```
