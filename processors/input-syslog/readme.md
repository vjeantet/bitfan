# SYSLOGINPUT


## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| format   | string | false    | ""            |
| port     | int    | false    |             0 |
| protocol | string | false    | ""            |


## Details

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



## Configuration blueprint

```
sysloginput{
	format => ""
	port => 123
	protocol => ""
}
```
