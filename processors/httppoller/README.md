# HTTPPOLLER
HTTPPoller allows you to intermittently poll remote HTTP URL, decode the output into an event

## Synopsys


|    SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------------|--------|----------|---------------|
| codec          | codec  | false    | "plain"       |
| interval       | string | false    | ""            |
| method         | string | false    | "GET"         |
| headers        | hash   | false    | {}            |
| url            | string | true     | ""            |
| target         | string | false    | ""            |
| ignore_failure | bool   | false    | true          |


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



## Configuration blueprint

```
httppoller{
	codec => "plain"
	interval => "every_10s"
	method => "GET"
	headers => {"User-Agent":"Bitfan","Accept":"application/json"}
	url=> "http://google.fr"
	target => ""
	ignore_failure => true
}
```
