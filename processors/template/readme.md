# TEMPLATEPROCESSOR


## Synopsys


| SETTING  |   TYPE   | REQUIRED | DEFAULT VALUE |
|----------|----------|----------|---------------|
| location | location | true     | ?             |
| var      | hash     | false    | {}            |
| target   | string   | false    | "output"      |


## Details

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
* Default value is `"output"`

Define the target field for placing the template execution result. If this setting is omitted,
the data will be stored in the "output" field



## Configuration blueprint

```
templateprocessor{
	location => "test.tpl"
	var => {"hostname"=>"myhost","varname"=>"varvalue"}
	target => "mydata"
}
```
