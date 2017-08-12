# HTTPOUTPROCESSOR
Display on http the last received event

URL is available as http://webhookhost/pluginLabel/URI

* webhookhost is defined by bitfan at startup
* pluginLabel is defined in pipeline configuration, it's the named processor if you put one, or `input_httpserver` by default
* URI is defined in plugin configuration (see below)

## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| Codec   | codec  | false    | "json"        |
| Uri     | string | false    | "out"         |
| Headers | hash   | false    | {}            |


## Details

### Codec
* Value type is codec
* Default value is `"json"`

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
	codec => "json"
	uri => "out"
	headers => {}
}
```
