# WEBSOCKETINPUT
Receive event on a ws connection

## Synopsys


|     SETTING      |  TYPE  | REQUIRED | DEFAULT VALUE |
|------------------|--------|----------|---------------|
| Codec            | codec  | false    | "json"        |
| Uri              | string | false    | "wsin"        |
| max_message_size | int    | false    |             0 |


## Details

### Codec
* Value type is codec
* Default value is `"json"`

The codec used for outputed data.

### Uri
* Value type is string
* Default value is `"wsin"`

URI path

### max_message_size
* Value type is int
* Default value is `0`

Maximum message size allowed from peer.



## Configuration blueprint

```
websocketinput{
	codec => "json"
	uri => "wsin"
	max_message_size => 123
}
```
