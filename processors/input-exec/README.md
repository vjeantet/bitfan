# EXECINPUT


## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| Command  | string | false    | ""            |
| Args     | array  | false    | []            |
| Interval | string | false    | ""            |
| codec    | codec  | false    | "plain"       |


## Details

### Command
* Value type is string
* Default value is `""`



### Args
* Value type is array
* Default value is `[]`



### Interval
* Value type is string
* Default value is `""`



### codec
* Value type is codec
* Default value is `"plain"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline



## Configuration blueprint

```
execinput{
	command => ""
	args => []
	interval => ""
	codec => "plain"
}
```
