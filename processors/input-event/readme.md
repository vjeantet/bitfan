# INPUTEVENTPROCESSOR
Generate a blank event on interval

## Synopsys


|  SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------|--------|----------|---------------|
| Add_field | hash   | false    | {}            |
| Tags      | array  | false    | []            |
| Type      | string | false    | ""            |
| Message   | string | false    | ""            |
| interval  | string | true     | ""            |


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

### Message
* Value type is string
* Default value is `""`



### interval
* This is a required setting.
* Value type is string
* Default value is `""`

Use CRON or BITFAN notation



## Configuration blueprint

```
inputeventprocessor{
	add_field => {}
	tags => []
	type => ""
	message => ""
	interval => "every_10s"
}
```
