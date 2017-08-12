# NEWTERM
This processor matches when a new value appears in a field that has never been seen before.

## Synopsys


|    SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------------|--------|----------|---------------|
| compare_field  | string | true     | ""            |
| ignore_missing | bool   | false    | true          |
| terms          | array  | false    | []            |


## Details

### compare_field
* This is a required setting.
* Value type is string
* Default value is `""`

The name of the field to use to compare to terms list.
If the field is null, those events will be ignored.

### ignore_missing
* Value type is bool
* Default value is `true`

If true, events without a compare_field field will be ignored.

### terms
* Value type is array
* Default value is `[]`

A list of initial terms to consider now new.
The compare_field term must be in this list or else it will match.



## Configuration blueprint

```
newterm{
	compare_field => "message"
	ignore_missing => true
	terms => ["val1","val2","val3"]
}
```
