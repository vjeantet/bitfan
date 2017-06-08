# EVALPROCESSOR
Modify or add event's field with the result of an expression (math or compare)

**Operators and types supported in expression :**

* Modifiers: `+` `-` `/` `*` `&` `|` `^` `**` `%` `>>` `<<`
* Comparators: `>` `>=` `<` `<=` `==` `!=` `=~` `!~`
* Logical ops: `||` `&&`
* Numeric constants, as 64-bit floating point (`12345.678`)
* String constants (single quotes: `'foobar'`)
* Date constants (single quotes, using any permutation of RFC3339, ISO8601, ruby date, or unix date; date parsing is automatically tried with any string constant)
* Boolean constants: `true` `false`
* Parenthesis to control order of evaluation `(` `)`
* Arrays (anything separated by `,` within parenthesis: `(1, 2, 'foo')`)
* Prefixes: `!` `-` `~`
* Ternary conditional: `?` `:`
* Null coalescence: `??`

## Synopsys


|   SETTING    | TYPE  | REQUIRED | DEFAULT VALUE |
|--------------|-------|----------|---------------|
| add_field    | hash  | false    | {}            |
| add_tag      | array | false    | []            |
| remove_field | array | false    | []            |
| remove_tag   | array | false    | []            |
| expressions  | hash  | true     | {}            |


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

### expressions
* This is a required setting.
* Value type is hash
* Default value is `{}`

list of field to set with expression's result



## Configuration blueprint

```
evalprocessor{
	add_field => {}
	add_tag => []
	remove_field => []
	remove_tag => []
	expressions => { "usage" => "[usage] * 100" }
}
```
