# STOPPROCESSOR
Stop after emiting a blank event on start
Allow you to put first event and then stop processors as soon as they finish their job.

Permit to launch bitfan with a pipeline and quit when work is done.

## Synopsys


|  SETTING   |  TYPE  | REQUIRED | DEFAULT VALUE |
|------------|--------|----------|---------------|
| Add_field  | hash   | false    | {}            |
| Tags       | array  | false    | []            |
| Type       | string | false    | ""            |
| ExitBitfan | bool   | false    | true          |


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

### ExitBitfan
* Value type is bool
* Default value is `true`

Stop bitfan with the pipeline ending ?



## Configuration blueprint

```
stopprocessor{
	add_field => {}
	tags => []
	type => ""
	exitbitfan => true
}
```
