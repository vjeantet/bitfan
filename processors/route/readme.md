# ROUTE
This processor, depending on the condition evaluation, will route message to
one or more different pipelines and/or pass the message through the processor to the next one.
Behavior :

* WHEN Condition is evaluated to true THEN the message go to the pipelines set in Path
* WHEN Condition is evaluated to true AND Fork set to true THEN the message go to the pipeline set in Path AND pass through.
* WHEN Condition is evaluated to false THEN the message pass through.
* WHEN Condition is evaluated to false AND Fork set to true THEN the message  pass through.

## Synopsys


|   SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|--------|----------|---------------|
| Add_field    | hash   | false    | {}            |
| Add_tag      | array  | false    | []            |
| condition    | string | false    | ""            |
| fork         | bool   | false    | false         |
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

### condition
* Value type is string
* Default value is `""`

set a condition to fork and route message
when false, message is routed to trunk
By default condition is evaluated to true

### fork
* Value type is bool
* Default value is `false`

Fork mode disabled by default

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

Path to configuration to send the incomming message, it could be a local file or an url
can be relative path to the current configuration.

### var
* Value type is hash
* Default value is `{}`

You can set variable references in the used configuration by using ${var}.
each reference will be replaced by the value of the variable found in this option
The replacement is case-sensitive.



## Configuration blueprint

```
route{
	add_field => {}
	add_tag => []
	condition => ""
	fork => false
	remove_field => []
	remove_tag => []
	type => ""
	path=> ["error.conf"]
	var => {"hostname"=>"myhost","varname"=>"varvalue"}
}
```
