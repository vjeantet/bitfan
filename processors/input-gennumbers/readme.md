# GENNUMBERS


## Synopsys


| SETTING  |   TYPE   | REQUIRED | DEFAULT VALUE |
|----------|----------|----------|---------------|
| count    | int      | false    |       1000000 |
| interval | interval | false    | ?             |


## Details

### count
* Value type is int
* Default value is `1000000`

How many events to generate

### interval
* Value type is interval
* Default value is `?`





## Configuration blueprint

```
gennumbers{
	count => 1000000
	interval => "10"
}
```
