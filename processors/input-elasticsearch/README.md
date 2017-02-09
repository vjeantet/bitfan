# ELASTICINPUT


## Synopsys


|  SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------|--------|----------|---------------|
| Add_field | hash   | false    | {}            |
| Tags      | array  | false    | []            |
| Type      | string | false    | ""            |
| Codec     | string | false    | ""            |
| Hosts     | array  | false    | []            |
| Query     | string | false    | ""            |
| Size      | int    | false    |             0 |
| User      | string | false    | ""            |
| Password  | string | false    | ""            |


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

### Hosts
* Value type is array
* Default value is `[]`



### Query
* Value type is string
* Default value is `""`



### Size
* Value type is int
* Default value is `0`



### User
* Value type is string
* Default value is `""`



### Password
* Value type is string
* Default value is `""`





## Configuration blueprint

```
elasticinput{
	add_field => {}
	tags => []
	type => ""
	codec => ""
	hosts => []
	query => ""
	size => 123
	user => ""
	password => ""
}
```
