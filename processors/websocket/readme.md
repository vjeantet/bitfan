# WEBSOCKET
Send event received over a ws connection to connected clients

## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| Codec   | codec  | false    | "json"        |
| Uri     | string | false    | "wsout"       |


## Details

### Codec
* Value type is codec
* Default value is `"json"`

The codec used for outputed data.

### Uri
* Value type is string
* Default value is `"wsout"`

URI path



## Configuration blueprint

```
websocket{
	codec => "json"
	uri => "wsout"
}
```
