# UDPINPUT


## Synopsys


|  SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------|--------|----------|---------------|
| add_field | hash   | false    | {}            |
| port      | int    | false    |             0 |
| tags      | array  | false    | []            |
| type      | string | false    | ""            |
| codec     | string | false    | ""            |


## Details

### add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### port
* Value type is int
* Default value is `0`

UDP port number to listen on

### tags
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.

### type
* Value type is string
* Default value is `""`

Add a type field to all events handled by this input

### codec
* Value type is string
* Default value is `""`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline



## Configuration blueprint

```
udpinput{
	add_field => {}
	port => 123
	tags => []
	type => ""
	codec => ""
}
```
