# JSON


## Synopsys


|       SETTING        |  TYPE  | REQUIRED |     DEFAULT VALUE     |
|----------------------|--------|----------|-----------------------|
| skip_on_invalid_json | bool   | false    | false                 |
| source               | string | true     | ""                    |
| target               | string | false    | ""                    |
| tag_on_failure       | array  | false    | ["_jsonparsefailure"] |


## Details

### skip_on_invalid_json
* Value type is bool
* Default value is `false`

Allow to skip filter on invalid json

### source
* This is a required setting.
* Value type is string
* Default value is `""`

The configuration for the JSON filter

### target
* Value type is string
* Default value is `""`

Define the target field for placing the parsed data. If this setting is omitted,
the JSON data will be stored at the root (top level) of the event

### tag_on_failure
* Value type is array
* Default value is `["_jsonparsefailure"]`

Append values to the tags field when there has been no successful match



## Configuration blueprint

```
json{
	skip_on_invalid_json => false
	source => ""
	target => ""
	tag_on_failure => ["_jsonparsefailure"]
}
```
