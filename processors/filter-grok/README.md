# GROK


## Synopsys


|       SETTING       | TYPE  | REQUIRED | DEFAULT VALUE |
|---------------------|-------|----------|---------------|
| break_on_match      | bool  | false    | ?             |
| keep_empty_captures | bool  | false    | ?             |
| match               | hash  | true     | {}            |
| named_capture_only  | bool  | false    | ?             |
| patterns_dir        | array | false    | []            |
| tag_on_failure      | array | false    | []            |


## Details

### break_on_match
* Value type is bool
* Default value is `?`

Break on first match. The first successful match by grok will result in the filter being
finished. If you want grok to try all patterns (maybe you are parsing different things),
then set this to false

### keep_empty_captures
* Value type is bool
* Default value is `?`

If true, keep empty captures as event fields

### match
* This is a required setting.
* Value type is hash
* Default value is `{}`

A hash of matches of field ⇒ value
@nodefault

For example:
```
    filter {
      grok { match => { "message" => "Duration: %{NUMBER:duration}" } }
    }
```
If you need to match multiple patterns against a single field, the value can be an array of patterns
```
    filter {
      grok { match => { "message" => [ "Duration: %{NUMBER:duration}", "Speed: %{NUMBER:speed}" ] } }
    }
```

### named_capture_only
* Value type is bool
* Default value is `?`

If true, only store named captures from grok.

### patterns_dir
* Value type is array
* Default value is `[]`

BitFan ships by default with a bunch of patterns, so you don’t necessarily need to
define this yourself unless you are adding additional patterns. You can point to
multiple pattern directories using this setting Note that Grok will read all files
in the directory and assume its a pattern file (including any tilde backup files)

### tag_on_failure
* Value type is array
* Default value is `[]`

Append values to the tags field when there has been no successful match



## Configuration blueprint

```
grok{
	break_on_match => bool
	keep_empty_captures => bool
	match => {}
	named_capture_only => bool
	patterns_dir => []
	tag_on_failure => []
}
```
