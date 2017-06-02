# SLEEPPROCESSOR
This is useful for rate limiting, etc.

## Synopsys


|  SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------|--------|----------|---------------|
| Add_field | hash   | false    | {}            |
| Tags      | array  | false    | []            |
| Type      | string | false    | ""            |
| Time      | int    | false    |             0 |


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

### Time
* Value type is int
* Default value is `0`

The length of time to sleep, in Millisecond, for every event.



## Configuration blueprint

```
sleepprocessor{
	add_field => {}
	tags => []
	type => ""
	time => 123
}
```
