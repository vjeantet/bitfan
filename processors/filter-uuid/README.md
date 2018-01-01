# UUID
The uuid filter allows you to generate a UUID and add it as a field to each processed event.

This is useful if you need to generate a string that’s unique for every event, even if the same input is processed multiple times. If you want to generate strings that are identical each time a event with a given content is processed (i.e. a hash) you should use the fingerprint filter instead.

The generated UUIDs follow the version 4 definition in RFC 4122).

## Synopsys


|  SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------|--------|----------|---------------|
| Overwrite | bool   | false    | false         |
| target    | string | true     | ""            |


## Details

### Overwrite
* Value type is bool
* Default value is `false`

If the value in the field currently (if any) should be overridden by the generated UUID.
Defaults to false (i.e. if the field is present, with ANY value, it won’t be overridden)

### target
* This is a required setting.
* Value type is string
* Default value is `""`

Add a UUID to a field



## Configuration blueprint

```
uuid{
	overwrite => bool
	target => ""
}
```
