# KV
This filter helps automatically parse messages (or specific event fields)
which are of the foo=bar variety.

## Synopsys


|        SETTING         |  TYPE  | REQUIRED | DEFAULT VALUE |
|------------------------|--------|----------|---------------|
| add_field              | hash   | false    | {}            |
| add_tag                | array  | false    | []            |
| allow_duplicate_values | bool   | false    | ?             |
| default_keys           | hash   | false    | {}            |
| exclude_keys           | array  | false    | []            |
| field_split            | string | false    | ""            |
| include_brackets       | bool   | false    | ?             |
| include_keys           | array  | false    | []            |
| Prefix                 | string | false    | ""            |
| Recursive              | bool   | false    | ?             |
| remove_field           | array  | false    | []            |
| remove_tag             | array  | false    | []            |
| Source                 | string | false    | ""            |
| Target                 | string | false    | ""            |
| Trim                   | string | false    | ""            |
| trimkey                | string | false    | ""            |
| value_split            | string | false    | ""            |


## Details

### add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### add_tag
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.

### allow_duplicate_values
* Value type is bool
* Default value is `?`

A bool option for removing duplicate key/value pairs.
When set to false, only one unique key/value pair will be preserved.
For example, consider a source like from=me from=me.
[from] will map to an Array with two elements: ["me", "me"].
to only keep unique key/value pairs, you could use this configuration
` kv {
`   allow_duplicate_values => false
` }

### default_keys
* Value type is hash
* Default value is `{}`

A hash specifying the default keys and their values which should be added
to the event in case these keys do not exist in the source field being parsed.
Example
`kv {
`  default_keys => { "from"=> "logstash@example.com",
`                   "to"=> "default@dev.null" }
`}

### exclude_keys
* Value type is array
* Default value is `[]`

An array specifying the parsed keys which should not be added to the event.
By default no keys will be excluded.
For example, consider a source like Hey, from=<abc>, to=def foo=bar.
To exclude from and to, but retain the foo key, you could use this configuration:
`kv {
`  exclude_keys => [ "from", "to" ]
`}

### field_split
* Value type is string
* Default value is `""`

A string of characters to use as delimiters for parsing out key-value pairs.
These characters form a regex character class and thus you must escape special regex characters like [ or ] using \.
## Example with URL Query Strings
For example, to split out the args from a url query string such as ?pin=12345~0&d=123&e=foo@bar.com&oq=bobo&ss=12345:
` kv {
`   field_split => "&?"
` }
The above splits on both & and ? characters, giving you the following fields:
* pin: 12345~0
* d: 123
* e: foo@bar.com
* oq: bobo
* ss: 12345

### include_brackets
* Value type is bool
* Default value is `?`

A boolean specifying whether to include brackets as value wrappers (the default is true)
` kv {
`   include_brackets => true
` }
For example, the result of this line: bracketsone=(hello world) bracketstwo=[hello world]
will be:
* bracketsone: hello world
* bracketstwo: hello world
instead of:
* bracketsone: (hello
* bracketstwo: [hello

### include_keys
* Value type is array
* Default value is `[]`

An array specifying the parsed keys which should be added to the event. By default all keys will be added.
For example, consider a source like Hey, from=<abc>, to=def foo=bar. To include from and to, but exclude the foo key, you could use this configuration:
` kv {
` include_keys => [ "from", "to" ]
` }

### Prefix
* Value type is string
* Default value is `""`

A string to prepend to all of the extracted keys.
For example, to prepend arg_ to all keys:
` kv {
`   prefix => "arg_" }
` }

### Recursive
* Value type is bool
* Default value is `?`

A boolean specifying whether to drill down into values and recursively get more key-value pairs from it. The extra key-value pairs will be stored as subkeys of the root key.
Default is not to recursive values.
` kv {
`  recursive => "true"
` }

### remove_field
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary fields from this event. Example:
` kv {
`   remove_field => [ "foo_%{somefield}" ]
` }

### remove_tag
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary tags from the event. Tags can be dynamic and include parts of the event using the %{field} syntax.
Example:
` kv {
`   remove_tag => [ "foo_%{somefield}" ]
` }
If the event has field "somefield" == "hello" this filter, on success, would remove the tag foo_hello if it is present. The second example would remove a sad, unwanted tag as well.

### Source
* Value type is string
* Default value is `""`

The field to perform key=value searching on
For example, to process the not_the_message field:
` kv { source => "not_the_message" }

### Target
* Value type is string
* Default value is `""`

The name of the container to put all of the key-value pairs into.
If this setting is omitted, fields will be written to the root of the event, as individual fields.
For example, to place all keys into the event field kv:
` kv { target => "kv" }

### Trim
* Value type is string
* Default value is `""`

A string of characters to trim from the value. This is useful if your values are wrapped in brackets or are terminated with commas (like postfix logs).
For example, to strip <, >, [, ] and , characters from values:
` kv {
`   trim => "<>[],"
` }

### trimkey
* Value type is string
* Default value is `""`

A string of characters to trim from the key. This is useful if your keys are wrapped in brackets or start with space.
For example, to strip < > [ ] and , characters from keys:
` kv {
`   trimkey => "<>[],"
` }

### value_split
* Value type is string
* Default value is `""`

A string of characters to use as delimiters for identifying key-value relations.
These characters form a regex character class and thus you must escape special regex characters like [ or ] using \.
For example, to identify key-values such as key1:value1 key2:value2:
` { kv { value_split => ":" }



## Configuration blueprint

```
kv{
	add_field => {}
	add_tag => []
	allow_duplicate_values => bool
	default_keys => {}
	exclude_keys => []
	field_split => ""
	include_brackets => bool
	include_keys => []
	prefix => ""
	recursive => bool
	remove_field => []
	remove_tag => []
	source => ""
	target => ""
	trim => ""
	trimkey => ""
	value_split => ""
}
```
