# CHANGE
This rule will monitor a certain field and match if that field changes. The field must change with respect to the last event

## Synopsys


|    SETTING     |  TYPE  | REQUIRED |  DEFAULT VALUE   |
|----------------|--------|----------|------------------|
| compare_field  | string | true     | ""               |
| ignore_missing | bool   | false    | true             |
| timeframe      | int    | false    | 0 (no timeframe) |


## Details

### compare_field
* This is a required setting.
* Value type is string
* Default value is `""`

The name of the field to use to compare to the blacklist.

If the field is null, those events will be ignored.

### ignore_missing
* Value type is bool
* Default value is `true`

If true, events without a compare_key field will not count as changed.

### timeframe
* Value type is int
* Default value is `0 (no timeframe)`

The maximum time in seconds between changes. After this time period, Bitfan will forget the old value of the compare_field field.



## Configuration blueprint

```
change{
	compare_field => "message"
	ignore_missing => true
	timeframe => 10
}
```
