# TEMPLATEPROCESSOR


## Synopsys


|  SETTING  |   TYPE   | REQUIRED | DEFAULT VALUE |
|-----------|----------|----------|---------------|
| Add_field | hash     | false    | {}            |
| Tags      | array    | false    | []            |
| Type      | string   | false    | ""            |
| location  | location | true     | ?             |
| var       | hash     | false    | {}            |
| target    | string   | false    | "generated"   |


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

### location
* This is a required setting.
* Value type is location
* Default value is `?`

Go Template content

set inline content, a path or an url to the template content

Go template : https://golang.org/pkg/html/template/

### var
* Value type is hash
* Default value is `{}`

You can set variable to be used in template by using ${var}.
each reference will be replaced by the value of the variable found in Template's content
The replacement is case-sensitive.

### target
* Value type is string
* Default value is `"generated"`

Define the target field for placing the template execution result. If this setting is omitted,
the data will be stored in the "data" field



## Configuration blueprint

```
templateprocessor{
	add_field => {}
	tags => []
	type => ""
	location => "test.tpl"
	var => {"hostname"=>"myhost","varname"=>"varvalue"}
	target => "mydata"
}
```
