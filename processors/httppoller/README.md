# HTTPPOLLER
HTTPPoller allows you to call an HTTP Endpoint, decode the output into an event

## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| codec    | codec  | false    | "plain"       |
| interval | string | false    | ""            |
| method   | string | false    | "GET"         |
| url      | string | true     | ""            |
| target   | string | false    | ""            |


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



## Configuration blueprint

```
httppoller{
	codec => "plain"
	interval => "every_10s"
	method => "GET"
	url=> "http://google.fr"
	target => ""
}
```
