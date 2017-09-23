# HTTPPOLLER
HTTPPoller allows you to intermittently poll remote HTTP URL, decode the output into an event

## Synopsys


|     SETTING      |  TYPE  | REQUIRED | DEFAULT VALUE |
|------------------|--------|----------|---------------|
| codec            | codec  | false    | "plain"       |
| interval         | string | false    | ""            |
| method           | string | false    | "GET"         |
| url              | string | true     | ""            |
| target           | string | false    | ""            |
| failure_severity | int    | false    |             0 |
| tag_on_failure   | array  | false    | []            |


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

### url
* This is a required setting.
* Value type is string
* Default value is `""`

URL

### target
* Value type is string
* Default value is `""`

When data is an array it stores the resulting data into the given target field.

### failure_severity
* Value type is int
* Default value is `0`

Level of failure

1 - noFailures
2 - unsuccessful HTTP requests (unreachable connections)
3 - unreachable connections and HTTP responses > 400 of successful HTTP requests
4 - unreachable connections and non-2xx HTTP responses of successful HTTP requests

### tag_on_failure
* Value type is array
* Default value is `[]`

When set, http failures will pass the received event and
append values to the tags field when there has been an failure



## Configuration blueprint

```
httppoller{
	codec => "plain"
	interval => "every_10s"
	method => "GET"
	url=> "http://google.fr"
	target => ""
	failure_severity => 123
	tag_on_failure => ["_httprequestfailure"]
}
```
