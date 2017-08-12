# STDIN
Read events from standard input.
By default, each event is assumed to be one line. If you want to join lines, you’ll want to use the multiline filter.

## Synopsys


| SETTING  | TYPE  | REQUIRED | DEFAULT VALUE |
|----------|-------|----------|---------------|
| Codec    | codec | false    | "line"        |
| eof_exit | bool  | false    | false         |


## Details

### Codec
* Value type is codec
* Default value is `"line"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

### eof_exit
* Value type is bool
* Default value is `false`

Stop bitfan on stdin EOF ? (use it when you pipe data with |)



## Configuration blueprint

```
stdin{
	codec => "line"
	eof_exit => false
}
```
