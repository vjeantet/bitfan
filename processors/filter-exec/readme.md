# EXEC
Execute a command and use its stdout as event data

## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| command | string | true     | ""            |
| args    | array  | false    | []            |
| stdin   | bool   | false    | false         |
| target  | string | false    | "stdout"      |
| codec   | codec  | false    | "plain"       |


## Details

### command
* This is a required setting.
* Value type is string
* Default value is `""`



### args
* Value type is array
* Default value is `[]`



### stdin
* Value type is bool
* Default value is `false`

Pass the complete event to stdin as a json string

### target
* Value type is string
* Default value is `"stdout"`

Where do the output should be stored
Set "." when output is json formated and want to replace current event fields with output
response. (useful)

### codec
* Value type is codec
* Default value is `"plain"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline



## Configuration blueprint

```
exec{
	command => ""
	args => []
	stdin => false
	target => "stdout"
	codec => "plain"
}
```
