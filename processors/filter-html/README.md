# HTML


## Synopsys


|    SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------------|--------|----------|---------------|
| source_field   | string | false    | ""            |
| text           | hash   | false    | {}            |
| size           | hash   | false    | {}            |
| tag_on_failure | array  | false    | []            |


## Details

### source_field
* Value type is string
* Default value is `""`

Which field contains the html document

### text
* Value type is hash
* Default value is `{}`

Add fields with the text of elements found with css selector

### size
* Value type is hash
* Default value is `{}`

Add fields with the number of elements found with css selector

### tag_on_failure
* Value type is array
* Default value is `[]`

Append values to the tags field when the html document can not be parsed



## Configuration blueprint

```
html{
	source_field => ""
	text => {}
	size => {}
	tag_on_failure => []
}
```
