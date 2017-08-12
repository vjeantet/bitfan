# INPUTEVENTPROCESSOR
Generate a blank event on interval

## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| Message  | string | false    | ""            |
| interval | string | true     | ""            |


## Details

### Message
* Value type is string
* Default value is `""`

string value to put in event

### interval
* This is a required setting.
* Value type is string
* Default value is `""`

Use CRON or BITFAN notation



## Configuration blueprint

```
inputeventprocessor{
	message => ""
	interval => "@every 10s"
}
```
