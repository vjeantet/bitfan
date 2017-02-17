# HTTPOUTPUT


## Synopsys


|     SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------------|--------|----------|---------------|
| add_field       | hash   | false    | {}            |
| url             | string | true     | ""            |
| headers         | hash   | false    | {}            |
| http_method     | string | false    | ""            |
| keepalive       | bool   | false    | ?             |
| pool_max        | int    | false    |             0 |
| connect_timeout | uint   | false    | ?             |
| request_timeout | uint   | false    | ?             |
| format          | string | false    | ""            |
| retryable_codes | array  | false    | []            |
| ignorable_codes | array  | false    | []            |
| batch_interval  | uint   | false    | ?             |
| batch_size      | uint   | false    | ?             |


## Details

### add_field
* Value type is hash
* Default value is `{}`

Add a field to an event. Default value is {}

### url
* This is a required setting.
* Value type is string
* Default value is `""`

This output lets you send events to a generic HTTP(S) endpoint
This setting can be dynamic using the %{foo} syntax.

### headers
* Value type is hash
* Default value is `{}`

Custom headers to use format is headers => {"X-My-Header", "%{host}"}. Default value is {}
This setting can be dynamic using the %{foo} syntax.

### http_method
* Value type is string
* Default value is `""`

The HTTP Verb. One of "put", "post", "patch", "delete", "get", "head". Default value is "post"

### keepalive
* Value type is bool
* Default value is `?`

Turn this on to enable HTTP keepalive support. Default value is true

### pool_max
* Value type is int
* Default value is `0`

Max number of concurrent connections. Default value is 1

### connect_timeout
* Value type is uint
* Default value is `?`

Timeout (in seconds) to wait for a connection to be established. Default value is 10

### request_timeout
* Value type is uint
* Default value is `?`

Timeout (in seconds) for the entire request. Default value is 60

### format
* Value type is string
* Default value is `""`

Set the format of the http body. Now supports only "json_lines"

### retryable_codes
* Value type is array
* Default value is `[]`

If encountered as response codes this plugin will retry these requests

### ignorable_codes
* Value type is array
* Default value is `[]`

If you would like to consider some non-2xx codes to be successes
enumerate them here. Responses returning these codes will be considered successes

### batch_interval
* Value type is uint
* Default value is `?`



### batch_size
* Value type is uint
* Default value is `?`





## Configuration blueprint

```
httpoutput{
	add_field => {}
	url => ""
	headers => {}
	http_method => ""
	keepalive => bool
	pool_max => 123
	connect_timeout => uint
	request_timeout => uint
	format => ""
	retryable_codes => []
	ignorable_codes => []
	batch_interval => uint
	batch_size => uint
}
```
