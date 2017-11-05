# HTTPOUTPUT


## Synopsys


|     SETTING     |  TYPE  | REQUIRED |         DEFAULT VALUE          |
|-----------------|--------|----------|--------------------------------|
| codec           | codec  | false    | "json"                         |
| add_field       | hash   | false    | {}                             |
| url             | string | true     | ""                             |
| headers         | hash   | false    | {"Content-Type" =>             |
|                 |        |          | "application/json"}            |
| http_method     | string | false    | "post"                         |
| keepalive       | bool   | false    | true                           |
| pool_max        | int    | false    |                              1 |
| connect_timeout | uint   | false    |                              5 |
| request_timeout | uint   | false    |                             30 |
| retryable_codes | array  | false    | [429, 500, 502, 503, 504]      |
| ignorable_codes | array  | false    | []                             |
| batch_interval  | uint   | false    |                              5 |
| batch_size      | uint   | false    |                            100 |
| retry_interval  | uint   | false    |                             10 |


## Details

### codec
* Value type is codec
* Default value is `"json"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

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
* Default value is `{"Content-Type" => "application/json"}`

Custom headers to use format is headers => {"X-My-Header", "%{host}"}. Default value is {}
This setting can be dynamic using the %{foo} syntax.

### http_method
* Value type is string
* Default value is `"post"`

The HTTP Verb. One of "put", "post", "patch", "delete", "get", "head". Default value is "post"

### keepalive
* Value type is bool
* Default value is `true`

Turn this on to enable HTTP keepalive support. Default value is true

### pool_max
* Value type is int
* Default value is `1`

Max number of concurrent connections. Default value is 1

### connect_timeout
* Value type is uint
* Default value is `5`

Timeout (in seconds) to wait for a connection to be established. Default value is 10

### request_timeout
* Value type is uint
* Default value is `30`

Timeout (in seconds) for the entire request. Default value is 60

### retryable_codes
* Value type is array
* Default value is `[429, 500, 502, 503, 504]`

If encountered as response codes this plugin will retry these requests

### ignorable_codes
* Value type is array
* Default value is `[]`

If you would like to consider some non-2xx codes to be successes
enumerate them here. Responses returning these codes will be considered successes

### batch_interval
* Value type is uint
* Default value is `5`



### batch_size
* Value type is uint
* Default value is `100`



### retry_interval
* Value type is uint
* Default value is `10`





## Configuration blueprint

```
httpoutput{
	codec => "json"
	add_field => {}
	url => ""
	headers => {"Content-Type" => "application/json"}
	http_method => "post"
	keepalive => true
	pool_max => 1
	connect_timeout => 5
	request_timeout => 30
	retryable_codes => [429, 500, 502, 503, 504]
	ignorable_codes => []
	batch_interval => 5
	batch_size => 100
	retry_interval => 10
}
```
