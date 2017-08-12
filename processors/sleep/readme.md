# SLEEPPROCESSOR
Sleep a given amount of time.

This will cause bitfan to stall for the given amount of time.

This is useful for rate limiting, etc.

## Synopsys


| SETTING | TYPE | REQUIRED | DEFAULT VALUE |
|---------|------|----------|---------------|
| Time    | int  | false    |             0 |


## Details

### Time
* Value type is int
* Default value is `0`

The length of time to sleep, in Millisecond, for every event.



## Configuration blueprint

```
sleepprocessor{
	time => 123
}
```
