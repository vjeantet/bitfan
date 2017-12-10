# WEBSOCKETINPUT
Receive event on a ws connection

## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| Codec   | codec  | false    | "json"        |
| Uri     | string | false    | "wsin"        |


## Details

### Codec
* Value type is codec
* Default value is `"json"`

The codec used for outputed data.

### Uri
* Value type is string
* Default value is `"wsin"`

URI path



## Configuration blueprint

```
websocketinput{
	codec => "json"
	uri => "wsin"
}
```
