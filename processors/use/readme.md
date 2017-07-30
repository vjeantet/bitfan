# USE
When used in input (input->filter->o) the processor will receive events from the last filter from the pipeline used in configuration file.

When used in filter (i->filter->o) the processor will

* pass the event to the first filter plugin found in the used configuration file
* receive events from the last filter plugin found in the used configuration file

When used in output (i->filter->output->o) the processor will

* pass the event to the first filter plugin found in the used configuration file

## Synopsys


|   SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|--------|----------|---------------|
| Add_field    | hash   | false    | {}            |
| Add_tag      | array  | false    | []            |
| Remove_field | array  | false    | []            |
| Remove_tag   | array  | false    | []            |
| Type         | string | false    | ""            |
| path         | array  | true     | []            |
| var          | hash   | false    | {}            |


## Details

### Add_field
* Value type is hash
* Default value is `{}`

If this processor is successful, add any arbitrary fields to this event.

### Add_tag
* Value type is array
* Default value is `[]`

If this processor is successful, add arbitrary tags to the event.
Tags can be dynamic and include parts of the event using the %{field} syntax.

### Remove_field
* Value type is array
* Default value is `[]`

If this processor is successful, remove arbitrary fields from this event.

### Remove_tag
* Value type is array
* Default value is `[]`

If this processor is successful, remove arbitrary tags from the event.
Tags can be dynamic and include parts of the event using the %{field} syntax

### Type
* Value type is string
* Default value is `""`

Add a type field to all events handled by this processor

### path
* This is a required setting.
* Value type is array
* Default value is `[]`

Path to configuration to import in this pipeline, it could be a local file or an url
can be relative path to the current configuration.

SPLIT and JOIN : in filter Section, set multiples path to make a split and join into your pipeline

### var
* Value type is hash
* Default value is `{}`

You can set variable references in the used configuration by using ${var}.
each reference will be replaced by the value of the variable found in this option

The replacement is case-sensitive.



## Configuration blueprint

```
use{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	type => ""
	path=> ["meteo-input.conf"]
	var => {"hostname"=>"myhost","varname"=>"varvalue"}
}
```
