# DROP
Drops everything received
Drops everything that gets to this filter.

This is best used in combination with conditionals, for example:
```
filter {
  if [loglevel] == "debug" {
    drop { }
  }
}
```
The above will only pass events to the drop filter if the loglevel field is debug. This will cause all events matching to be dropped.

## Synopsys


|   SETTING    | TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|-------|----------|---------------|
| Add_field    | hash  | false    | {}            |
| Add_tag      | array | false    | []            |
| Remove_field | array | false    | []            |
| Remove_Tag   | array | false    | []            |
| Percentage   | int   | false    |             0 |


## Details

### Add_field
* Value type is hash
* Default value is `{}`

If this event survice to drop, add any arbitrary fields to this event.
Field names can be dynamic and include parts of the event using the %{field}.

### Add_tag
* Value type is array
* Default value is `[]`

If this event survice to drop, add arbitrary tags to the event.
Tags can be dynamic and include parts of the event using the %{field} syntax.

### Remove_field
* Value type is array
* Default value is `[]`

If this event survice to drop, remove arbitrary fields from this event.

### Remove_Tag
* Value type is array
* Default value is `[]`

If this event survice to drop, remove arbitrary tags from the event.
Tags can be dynamic and include parts of the event using the %{field} syntax

### Percentage
* Value type is int
* Default value is `0`

Drop all the events within a pre-configured percentage.
This is useful if you just need a percentage but not the whole.



## Configuration blueprint

```
drop{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	percentage => 123
}
```
