# SQLPROCESSOR


## Synopsys


|      SETTING      |  TYPE  | REQUIRED | DEFAULT VALUE |
|-------------------|--------|----------|---------------|
| Add_field         | hash   | false    | {}            |
| Tags              | array  | false    | []            |
| Type              | string | false    | ""            |
| Codec             | string | false    | ""            |
| driver            | string | true     | ""            |
| event_by          | string | false    | "row"         |
| statement         | string | true     | ""            |
| interval          | string | true     | ""            |
| connection_string | string | true     | ""            |
| var               | hash   | false    | {}            |
| target            | string | false    | "data"        |


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

### driver
* This is a required setting.
* Value type is string
* Default value is `""`

GOLANG driver class to load, for example, "mysql".

### event_by
* Value type is string
* Default value is `"row"`

Send an event row by row or one event with all results
possible values "row", "result"

### statement
* This is a required setting.
* Value type is string
* Default value is `""`

SQL Statement
When there is more than 1 statement, only data from the last one will generate events.

### interval
* This is a required setting.
* Value type is string
* Default value is `""`

Set an interval when this processor is used as a input

### connection_string
* This is a required setting.
* Value type is string
* Default value is `""`



### var
* Value type is hash
* Default value is `{}`

You can set variable to be used in Statements by using ${var}.
each reference will be replaced by the value of the variable found in Statement's content
The replacement is case-sensitive.

### target
* Value type is string
* Default value is `"data"`

Define the target field for placing the retrieved data. If this setting is omitted,
the data will be stored in the "data" field
Set the value to "." to store value to the root (top level) of the event



## Configuration blueprint

```
sqlprocessor{
	add_field => {}
	tags => []
	type => ""
	codec => ""
	driver => "mysql"
	event_by => "row"
	statement => "SELECT * FROM mytable"
	interval => "10"
	connection_string => "username:password@tcp(192.168.1.2:3306)/mydatabase?charset=utf8"
	var => {"hostname"=>"myhost","varname"=>"varvalue"}
	target => "data"
}
```
