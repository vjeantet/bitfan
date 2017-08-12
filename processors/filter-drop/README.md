# DROP
Drops everything received
Drops everything that gets to this filter.

This is best used in combination with conditionals, for example:
```
filter {
  if [loglevel] == "debug" {
    drop { }
  }
}
```
The above will only pass events to the drop filter if the loglevel field is debug. This will cause all events matching to be dropped.

## Synopsys


|  SETTING   | TYPE | REQUIRED | DEFAULT VALUE |
|------------|------|----------|---------------|
| Percentage | int  | false    |             0 |


## Details

### Percentage
* Value type is int
* Default value is `0`

Drop all the events within a pre-configured percentage.
This is useful if you just need a percentage but not the whole.



## Configuration blueprint

```
drop{
	percentage => 123
}
```
