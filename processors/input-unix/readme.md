# UNIXINPUT


## Synopsys


|   SETTING    |     TYPE      | REQUIRED | DEFAULT VALUE |
|--------------|---------------|----------|---------------|
| add_field    | hash          | false    | {}            |
| data_timeout | time.Duration | false    |               |
| force_unlink | bool          | false    | ?             |
| mode         | string        | false    | ""            |
| path         | string        | true     | ""            |
| tags         | array         | false    | []            |
| type         | string        | false    | ""            |
| codec        | string        | false    | ""            |


## Details

### add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### data_timeout
* Value type is time.Duration
* Default value is ``

The read timeout in seconds. If a particular connection is idle for more than this timeout period, we will assume it is dead and close it.
If you never want to timeout, use 0.
Default value is 0

### force_unlink
* Value type is bool
* Default value is `?`

Remove socket file in case of EADDRINUSE failure
Default value is false

### mode
* Value type is string
* Default value is `""`

Mode to operate in. server listens for client connections, client connects to a server.
Value can be any of: "server", "client"
Default value is "server"

### path
* This is a required setting.
* Value type is string
* Default value is `""`

When mode is server, the path to listen on. When mode is client, the path to connect to.

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
unixinput{
	add_field => {}
	data_timeout => 30
	force_unlink => bool
	mode => ""
	path => ""
	tags => []
	type => ""
	codec => ""
}
```
