# ELASTICINPUT


## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| Hosts    | array  | false    | []            |
| Query    | string | false    | ""            |
| Size     | int    | false    |             0 |
| User     | string | false    | ""            |
| Password | string | false    | ""            |


## Details

### Hosts
* Value type is array
* Default value is `[]`



### Query
* Value type is string
* Default value is `""`



### Size
* Value type is int
* Default value is `0`



### User
* Value type is string
* Default value is `""`



### Password
* Value type is string
* Default value is `""`





## Configuration blueprint

```
elasticinput{
	hosts => []
	query => ""
	size => 123
	user => ""
	password => ""
}
```
