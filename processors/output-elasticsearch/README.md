# ELASTICSEARCH


## Synopsys


| SETTING  |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------|--------|----------|---------------|
| Host     | string | false    | ""            |
| Cluster  | string | false    | ""            |
| Protocol | string | false    | ""            |
| Port     | int    | false    |             0 |
| User     | string | false    | ""            |
| Password | string | false    | ""            |


## Details

### Host
* Value type is string
* Default value is `""`



### Cluster
* Value type is string
* Default value is `""`



### Protocol
* Value type is string
* Default value is `""`



### Port
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
elasticsearch{
	host => ""
	cluster => ""
	protocol => ""
	port => 123
	user => ""
	password => ""
}
```
