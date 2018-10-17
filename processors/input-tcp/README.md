# TCPINPUT


## Synopsys


|     SETTING      | TYPE | REQUIRED | DEFAULT VALUE |
|------------------|------|----------|---------------|
| port             | int  | false    |             0 |
| read_buffer_size | int  | false    |             0 |


## Details

### port
* Value type is int
* Default value is `0`

TCP port number to listen on

### read_buffer_size
* Value type is int
* Default value is `0`

Message buffer size



## Configuration blueprint

```
tcp {
	port => 123
	read_buffer_size => 123
}
```
