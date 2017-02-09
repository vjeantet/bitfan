# MUTATE
mutate filter allows to perform general mutations on fields. You can rename, remove, replace, and modify fields in your event.

## Synopsys


|    SETTING     | TYPE  | REQUIRED | DEFAULT VALUE |
|----------------|-------|----------|---------------|
| Add_field      | hash  | false    | {}            |
| Add_tag        | array | false    | []            |
| Convert        | hash  | false    | {}            |
| Gsub           | array | false    | []            |
| Join           | hash  | false    | {}            |
| Lowercase      | array | false    | []            |
| Merge          | hash  | false    | {}            |
| Remove_field   | array | false    | []            |
| Remove_tag     | array | false    | []            |
| Rename         | hash  | false    | {}            |
| Replace        | hash  | false    | {}            |
| Split          | hash  | false    | {}            |
| Strip          | array | false    | []            |
| Update         | hash  | false    | {}            |
| Uppercase      | array | false    | []            |
| Remove_all_but | array | false    | []            |


## Details

### Add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### Add_tag
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event.
Tags can be dynamic and include parts of the event using the %{field} syntax.

### Convert
* Value type is hash
* Default value is `{}`

Convert a field’s value to a different type, like turning a string to an integer.
If the field value is an array, all members will be converted. If the field is a hash,
no action will be taken.
If the conversion type is boolean, the acceptable values are:
True: true, t, yes, y, and 1
False: false, f, no, n, and 0
If a value other than these is provided, it will pass straight through and log a warning message.
Valid conversion targets are: integer, float, string, and boolean.

### Gsub
* Value type is array
* Default value is `[]`

Convert a string field by applying a regular expression and a replacement. If the field is not a string, no action will be taken.
This configuration takes an array consisting of 3 elements per field/substitution.
Be aware of escaping any backslash in the config file.

### Join
* Value type is hash
* Default value is `{}`

Join an array with a separator character. Does nothing on non-array fields

### Lowercase
* Value type is array
* Default value is `[]`

Convert a value to its lowercase equivalent

### Merge
* Value type is hash
* Default value is `{}`

Merge two fields of arrays or hashes. String fields will be automatically be converted into an array

### Remove_field
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary fields from this event.

### Remove_tag
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary tags from the event.
Tags can be dynamic and include parts of the event using the %{field} syntax

### Rename
* Value type is hash
* Default value is `{}`

Rename key on one or more fields

### Replace
* Value type is hash
* Default value is `{}`

Replace a field with a new value. The new value can include %{foo} strings to
help you build a new value from other parts of the event

### Split
* Value type is hash
* Default value is `{}`

Split a field to an array using a separator character. Only works on string fields

### Strip
* Value type is array
* Default value is `[]`

Strip whitespace from processors. NOTE: this only works on leading and trailing whitespace

### Update
* Value type is hash
* Default value is `{}`

Update an existing field with a new value. If the field does not exist, then no action will be taken

### Uppercase
* Value type is array
* Default value is `[]`

Convert a value to its uppercase equivalent

### Remove_all_but
* Value type is array
* Default value is `[]`

remove all fields, except theses fields (work only with first level fields)



## Configuration blueprint

```
mutate{
	add_field => {}
	add_tag => []
	convert => {}
	gsub => []
	join => {}
	lowercase => []
	merge => {}
	remove_field => []
	remove_tag => []
	rename => {}
	replace => {}
	split => {}
	strip => []
	update => {}
	uppercase => []
	remove_all_but => []
}
```
