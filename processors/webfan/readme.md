# WEBFAN


## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| Codec   | codec  | false    | "json"        |
| uri     | string | true     | ""            |
| conf    | string | true     | ""            |
| headers | hash   | false    | {}            |


## Details

### Codec
* Value type is codec
* Default value is `"json"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

Default decode http request as json, response is json encoded.
Set multiple codec with role to customize

### uri
* This is a required setting.
* Value type is string
* Default value is `""`

URI path /_/path

### conf
* This is a required setting.
* Value type is string
* Default value is `""`

Path to pipeline's configuration to execute on request
This configuration should contains only a filter section an a output like ```output{pass{}}```

### headers
* Value type is hash
* Default value is `{}`

Headers to send back into each outgoing response



## Configuration blueprint

```
webfan{
	codec => "json"
	uri => ""
	conf => ""
	{"X-Processor" => "bitfan"}
}
```
