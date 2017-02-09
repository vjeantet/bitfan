# HTML


## Synopsys


|    SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------------|--------|----------|---------------|
| Add_field      | hash   | false    | {}            |
| Tags           | array  | false    | []            |
| Type           | string | false    | ""            |
| Codec          | string | false    | ""            |
| source_field   | string | false    | ""            |
| text           | hash   | false    | {}            |
| size           | hash   | false    | {}            |
| tag_on_failure | array  | false    | []            |


## Details

### Add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### Tags
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.

### Type
* Value type is string
* Default value is `""`

Add a type field to all events handled by this input

### Codec
* Value type is string
* Default value is `""`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

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
	add_field => {}
	tags => []
	type => ""
	codec => ""
	source_field => ""
	text => {}
	size => {}
	tag_on_failure => []
}
```
