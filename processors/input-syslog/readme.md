# SYSLOGINPUT


## Synopsys


|  SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------|--------|----------|---------------|
| add_field | hash   | false    | {}            |
| format    | string | false    | ""            |
| port      | int    | false    |             0 |
| protocol  | string | false    | ""            |
| tags      | array  | false    | []            |
| type      | string | false    | ""            |


## Details

### add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### format
* Value type is string
* Default value is `""`

Which format to use to decode syslog messages. Default value is "automatic"
Value can be "automatic","rfc3164", "rfc5424" or "rfc6587"

Note on "automatic" format: if you don't know which format to select,
or have multiple incoming formats, this is the one to go for.
There is a theoretical performance penalty (it has to look at a few bytes
at the start of the frame), and a risk that you may parse things you don't want to parse
(rogue syslog clients using other formats), so if you can be absolutely sure of your syslog
format, it would be best to select it explicitly.

### port
* Value type is int
* Default value is `0`

Port number to listen on

### protocol
* Value type is string
* Default value is `""`

Protocol to use to listen to syslog messages
Value can be either "tcp" or "udp"

### tags
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.

### type
* Value type is string
* Default value is `""`

Add a type field to all events handled by this input



## Configuration blueprint

```
sysloginput{
	add_field => {}
	format => ""
	port => 123
	protocol => ""
	tags => []
	type => ""
}
```
