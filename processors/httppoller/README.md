# HTTPPOLLER
HTTPPoller allows you to intermittently poll remote HTTP URL, decode the output into an event

## Synopsys


|    SETTING     |   TYPE   | REQUIRED | DEFAULT VALUE |
|----------------|----------|----------|---------------|
| codec          | codec    | false    | "plain"       |
| interval       | string   | false    | ""            |
| method         | string   | false    | "GET"         |
| headers        | hash     | false    | {}            |
| body           | location | false    | ?             |
| url            | string   | true     | ""            |
| target         | string   | false    | ""            |
| ignore_failure | bool     | false    | true          |
| var            | hash     | false    | {}            |


## Details

### codec
* Value type is codec
* Default value is `"plain"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

### interval
* Value type is string
* Default value is `""`

Use CRON or BITFAN notation

### method
* Value type is string
* Default value is `"GET"`

Http Method

### headers
* Value type is hash
* Default value is `{}`

Define headers for the request.

### body
* Value type is location
* Default value is `?`

The request body (e.g. for an HTTP POST request). No default body is specified

### url
* This is a required setting.
* Value type is string
* Default value is `""`

URL

### target
* Value type is string
* Default value is `""`

When data is an array it stores the resulting data into the given target field.

### ignore_failure
* Value type is bool
* Default value is `true`

When true, unsuccessful HTTP requests, like unreachable connections, will
not raise an event, but a log message.
When false an event is generated with a tag _http_request_failure

### var
* Value type is hash
* Default value is `{}`

You can set variable to be used in Body by using ${var}.
each reference will be replaced by the value of the variable found in Body's content
The replacement is case-sensitive.



## Configuration blueprint

```
httppoller{
	codec => "plain"
	interval => "every_10s"
	method => "GET"
	headers => {"User-Agent":"Bitfan","Accept":"application/json"}
	body => location
	url=> "http://google.fr"
	target => ""
	ignore_failure => true
	var => {"hostname"=>"myhost","varname"=>"varvalue"}
}
```
