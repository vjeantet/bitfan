+++
date = "2017-05-16T21:08:31+02:00"
description = ""
title = "Common options"
weight = 70
+++

## Synopsys

|       SETTING       | TYPE  | REQUIRED | DEFAULT VALUE |
|---------------------|-------|----------|---------------|
| add_field           | hash  | false    | {}            |
| add_tag             | array | false    | []            |
| remove_field        | array | false    | []            |
| remove_tag          | array | false    | []            |
| trace               | bool  | false    | false         |
| type               | string | false    | ""            |


## Details

### Add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### Add_tag
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.


### Remove_field
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary fields from this event.

### Remove_tag
* Value type is array
* Default value is `[]`

### Trace
* Value type is bool
* Default value is `false`

Log each event produced by the processor.
Usefull while building or debugging a configuration pipeline.

Trace event are logged at INFO level.


### Type
* Value type is string
* Default value is `""`

Add a type field to all events handled by this input



