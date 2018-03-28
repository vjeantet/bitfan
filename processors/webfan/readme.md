# WEBFAN
Example
```
input{
  webhook{
        uri => "toto/titi"
        pipeline=> "test.conf"
        codec => plain{
            role => "decoder"
        }
        codec => plain{
            role => "encoder"
            format=> "<h1>Hello {{.request.querystring.name}}</h1>"
        }
        headers => {
            "Content-Type" => "text/html"
        }
    }
}
```

## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| Codec    | codec  | false    | ?             |
| uri      | string | true     | ""            |
| pipeline | string | true     | ""            |
| headers  | hash   | false    | {}            |


## Details

### Codec
* Value type is codec
* Default value is `?`

The codec used for posted data. Input codecs are a convenient method for decoding
your data before it enters the pipeline, without needing a separate filter in your bitfan pipeline

Default decode http request as plain text, response is json encoded.
Set multiple codec with role to customize

### uri
* This is a required setting.
* Value type is string
* Default value is `""`

URI path /_/path

### pipeline
* This is a required setting.
* Value type is string
* Default value is `""`

Path to pipeline's configuration to execute on request
This configuration should contains only a filter section an a output like ```output{pass{}}```

### headers
* Value type is hash
* Default value is `{}`

Headers to send back into outgoing response



## Configuration blueprint

```
webfan{
	codec => plain { role=>"encoder"} codec => json { role=>"decoder"}
	uri => ""
	pipeline => ""
	{"X-Processor" => "bitfan"}
}
```
