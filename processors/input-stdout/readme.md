# INPUTSTDOUT


## Synopsys


| SETTING | TYPE  | REQUIRED | DEFAULT VALUE |
|---------|-------|----------|---------------|
| codec   | codec | false    | "line"        |


## Details

### codec
* Value type is codec
* Default value is `"line"`

Codec can be one of  "json", "line", "pp" or "rubydebug"



## Configuration blueprint

```
inputstdout{
	codec => "pp"
}
```
