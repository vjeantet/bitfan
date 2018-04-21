# TCPOUTPUT


## Synopsys


|     SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------------|--------|----------|---------------|
| codec           | codec  | false    | "line"        |
| host            | string | true     | ""            |
| port            | uint   | true     | ?             |
| keepalive       | bool   | false    | true          |
| request_timeout | uint   | false    |            30 |
| retry_interval  | uint   | false    |            10 |


## Details

### codec
* Value type is codec
* Default value is `"line"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

### host
* This is a required setting.
* Value type is string
* Default value is `""`



### port
* This is a required setting.
* Value type is uint
* Default value is `?`



### keepalive
* Value type is bool
* Default value is `true`

Turn this on to enable HTTP keepalive support. Default value is true

### request_timeout
* Value type is uint
* Default value is `30`

Timeout (in seconds) for the entire request. Default value is 60

### retry_interval
* Value type is uint
* Default value is `10`





## Configuration blueprint

```
tcpoutput{
	codec => "line"
	host => ""
	port => uint
	keepalive => true
	request_timeout => 30
	retry_interval => 10
}
```
