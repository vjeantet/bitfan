# UNIXINPUT


## Synopsys


|   SETTING    |     TYPE      | REQUIRED | DEFAULT VALUE |
|--------------|---------------|----------|---------------|
| data_timeout | time.Duration | false    |               |
| force_unlink | bool          | false    | false         |
| mode         | string        | false    | ""            |
| path         | string        | true     | ""            |
| codec        | string        | false    | ""            |


## Details

### data_timeout
* Value type is time.Duration
* Default value is ``

The read timeout in seconds. If a particular connection is idle for more than this timeout period, we will assume it is dead and close it.
If you never want to timeout, use 0.
Default value is 0

### force_unlink
* Value type is bool
* Default value is `false`

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

### codec
* Value type is string
* Default value is `""`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline



## Configuration blueprint

```
unixinput{
	data_timeout => 30
	force_unlink => bool
	mode => ""
	path => ""
	codec => ""
}
```
