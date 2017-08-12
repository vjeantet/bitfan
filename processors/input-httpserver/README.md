# HTTPSERVERPROCESSOR
Listen and read a http request to build events with it.

Processor respond with a HTTP code as :

* `202` when request has been accepted, in body : the total number of event created
* `500` when an error occurs, in body : an error description

Use codecs to process body content as json / csv / lines / json lines / ....

URL is available as http://webhookhost/pluginLabel/URI

* webhookhost is defined by bitfan at startup
* pluginLabel is defined in pipeline configuration, it's the named processor if you put one, or `input_httpserver` by default
* URI is defined in plugin configuration (see below)

## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| Codec   | codec  | false    | "plain"       |
| Uri     | string | false    | "events"      |


## Details

### Codec
* Value type is codec
* Default value is `"plain"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

### Uri
* Value type is string
* Default value is `"events"`

URI path



## Configuration blueprint

```
httpserverprocessor{
	codec => "plain"
	uri => "events"
}
```
