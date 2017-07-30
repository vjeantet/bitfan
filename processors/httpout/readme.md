# HTTPOUTPROCESSOR
Display on http the last received event

URL is available as http://webhookhost/pluginLabel/URI

* webhookhost is defined by bitfan at startup
* pluginLabel is defined in pipeline configuration, it's the named processor if you put one, or `input_httpserver` by default
* URI is defined in plugin configuration (see below)

## Synopsys


|  SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------|--------|----------|---------------|
| Add_field | hash   | false    | {}            |
| Tags      | array  | false    | []            |
| Type      | string | false    | ""            |
| Codec     | codec  | false    | "plain"       |
| Uri       | string | false    | "out"         |
| Headers   | hash   | false    | {}            |


## Details

### Add_field
* Value type is hash
* Default value is `{}`

Add a field to an event

### Tags
* Value type is array
* Default value is `[]`

Add any number of arbitrary tags to your event.
This can help with processing later.

### Type
* Value type is string
* Default value is `""`

Add a type field to all events handled by this input

### Codec
* Value type is codec
* Default value is `"plain"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

### Uri
* Value type is string
* Default value is `"out"`

URI path

### Headers
* Value type is hash
* Default value is `{}`

Add headers to output



## Configuration blueprint

```
httpoutprocessor{
	add_field => {}
	tags => []
	type => ""
	codec => "plain"
	uri => "out"
	headers => {}
}
```
