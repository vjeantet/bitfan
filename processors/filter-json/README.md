# JSON


## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| Source  | string | false    | ""            |
| Target  | string | false    | ""            |


## Details

### Source
* Value type is string
* Default value is `""`

The configuration for the JSON filter

### Target
* Value type is string
* Default value is `""`

Define the target field for placing the parsed data. If this setting is omitted,
the JSON data will be stored at the root (top level) of the event



## Configuration blueprint

```
json{
	source => ""
	target => ""
}
```
