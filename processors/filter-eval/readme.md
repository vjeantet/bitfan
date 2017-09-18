# EVALPROCESSOR
Modify or add event's field with the result of

* an expression (math or compare)
* an go template

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


|   SETTING   | TYPE | REQUIRED | DEFAULT VALUE |
|-------------|------|----------|---------------|
| expressions | hash | false    | {}            |
| templates   | hash | false    | {}            |
| var         | hash | false    | {}            |


## Details

### expressions
* Value type is hash
* Default value is `{}`

list of field to set with expression's result

### templates
* Value type is hash
* Default value is `{}`

list of field to set with a go template location

### var
* Value type is hash
* Default value is `{}`

You can set variable to be used in template by using ${var}.
each reference will be replaced by the value of the variable found in Template's content
The replacement is case-sensitive.



## Configuration blueprint

```
evalprocessor{
	expressions => { "usage" => "[usage] * 100" }
	expressions => { "count" => "{{len .data}}", "mail"=>"mytemplate.tpl" }
	var => {"hostname"=>"myhost","varname"=>"varvalue"}
}
```
