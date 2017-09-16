# INPUTEVENTPROCESSOR
Generate a blank event on interval

## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| Message  | string | false    | ""            |
| count    | int    | false    |             1 |
| interval | string | false    | ""            |


## Details

### Message
* Value type is string
* Default value is `""`

string value to put in event

### count
* Value type is int
* Default value is `1`

How many events to generate

### interval
* Value type is string
* Default value is `""`

Use CRON or BITFAN notation
When omited, event will be generated on start



## Configuration blueprint

```
inputeventprocessor{
	message => ""
	count => 1000000
	interval => "@every 10s"
}
```
