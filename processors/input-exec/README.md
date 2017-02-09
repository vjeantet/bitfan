# EXECINPUT


## Synopsys


|  SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------|--------|----------|---------------|
| Command   | string | false    | ""            |
| Args      | array  | false    | []            |
| Add_field | hash   | false    | {}            |
| Interval  | string | false    | ""            |
| Codec     | string | false    | ""            |
| Tags      | array  | false    | []            |
| Type      | string | false    | ""            |


## Details

### Command
* Value type is string
* Default value is `""`



### Args
* Value type is array
* Default value is `[]`



### Add_field
* Value type is hash
* Default value is `{}`



### Interval
* Value type is string
* Default value is `""`



### Codec
* Value type is string
* Default value is `""`



### Tags
* Value type is array
* Default value is `[]`



### Type
* Value type is string
* Default value is `""`





## Configuration blueprint

```
execinput{
	command => ""
	args => []
	add_field => {}
	interval => ""
	codec => ""
	tags => []
	type => ""
}
```
