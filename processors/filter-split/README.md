# SPLIT
The split filter clones an event by splitting one of its fields and placing each value resulting from the split into a clone of the original event. The field being split can either be a string or an array.

An example use case of this filter is for taking output from the exec input plugin which emits one event for the whole output of a command and splitting that output by newline - making each line an event.

The end result of each split is a complete copy of the event with only the current split section of the given field changed.

## Synopsys


|  SETTING   |  TYPE  | REQUIRED | DEFAULT VALUE |
|------------|--------|----------|---------------|
| Field      | string | false    | ""            |
| Target     | string | false    | ""            |
| Terminator | string | false    | ""            |


## Details

### Field
* Value type is string
* Default value is `""`

The field which value is split by the terminator

### Target
* Value type is string
* Default value is `""`

The field within the new event which the value is split into. If not set, target field defaults to split field name

### Terminator
* Value type is string
* Default value is `""`

The string to split on. This is usually a line terminator, but can be any string
Default value is "\n"



## Configuration blueprint

```
split{
	field => ""
	target => ""
	terminator => ""
}
```
