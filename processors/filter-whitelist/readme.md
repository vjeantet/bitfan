# WHITELIST
Similar to blacklist, this processor will compare a certain field to a whitelist, and match
if the list does not contain the term

## Synopsys


|    SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------------|--------|----------|---------------|
| compare_field  | string | true     | ""            |
| ignore_missing | bool   | false    | true          |
| terms          | array  | true     | []            |


## Details

### compare_field
* This is a required setting.
* Value type is string
* Default value is `""`

The name of the field to use to compare to the whitelist.
If the field is null, those events will be ignored.

### ignore_missing
* Value type is bool
* Default value is `true`

If true, events without a compare_key field will not match.

### terms
* This is a required setting.
* Value type is array
* Default value is `[]`

A list of whitelisted terms.
The compare_field term must be in this list or else it will match.



## Configuration blueprint

```
whitelist{
	compare_field => "message"
	ignore_missing => true
	terms => ["val1","val2","val3"]
}
```
