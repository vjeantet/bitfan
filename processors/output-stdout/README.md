# STDOUT
A simple output which prints to the STDOUT of the shell running BitFan. This output can be quite convenient when debugging plugin configurations, by allowing instant access to the event data after it has passed through the inputs and filters.

For example, the following output configuration, in conjunction with the BitFan -e command-line flag, will allow you to see the results of your event pipeline for quick iteration.
```
output {
  stdout {}
}
```
Useful codecs include:

pp: outputs event data using the go "k0kubun/pp" package
if codec is rubydebug, it will treated as "pp"
```
output {
  stdout { codec => pp }
}
```
json: outputs event data in structured JSON format
```
output {
  stdout { codec => json }
}
```

## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| Codec   | string | false    | "line"        |


## Details

### Codec
* Value type is string
* Default value is `"line"`

Codec can be one of  "json", "line", "pp" or "rubydebug"



## Configuration blueprint

```
stdout{
	codec => "pp"
}
```
