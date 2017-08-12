# DIGEST


## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| key_map  | string | false    | ""            |
| interval | string | false    | ""            |
| count    | int    | false    |             0 |


## Details

### key_map
* Value type is string
* Default value is `""`

Add received event fields to the digest field named with the key map_key
When this setting is empty, digest will merge fields from coming events

### interval
* Value type is string
* Default value is `""`

When should Digest send a digested event ?
Use CRON or BITFAN notation

### count
* Value type is int
* Default value is `0`

With min > 0, digest will not fire an event if less than min events were digested



## Configuration blueprint

```
digest{
	key_map => "type"
	interval => "every_10s"
	count => 123
}
```
