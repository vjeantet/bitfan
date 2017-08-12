# BLACKLIST
The blacklist rule will check a certain field against a blacklist, and match if it is in the blacklist.

## Synopsys


|    SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------------|--------|----------|---------------|
| compare_field | string | true     | ""            |
| terms         | array  | true     | []            |


## Details

### compare_field
* This is a required setting.
* Value type is string
* Default value is `""`

The name of the field to use to compare to the blacklist.
If the field is null, those events will be ignored.

### terms
* This is a required setting.
* Value type is array
* Default value is `[]`

List of blacklisted terms.
The compare_field term must be equal to one of these values for it to match.



## Configuration blueprint

```
blacklist{
	compare_field => "message"
	terms => ["val1","val2","val3"]
}
```
