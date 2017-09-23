# GLUSTERFSOUTPUT


## Synopsys


|      SETTING      |     TYPE      | REQUIRED | DEFAULT VALUE |
|-------------------|---------------|----------|---------------|
| codec             | string        | false    | ""            |
| create_if_deleted | bool          | false    | false         |
| dir_mode          | os.FileMode   | false    | ?             |
| file_mode         | os.FileMode   | false    | ?             |
| flush_interval    | time.Duration | false    |               |
| host              | string        | false    | ""            |
| path              | string        | true     | ""            |
| volume            | string        | true     | ""            |


## Details

### codec
* Value type is string
* Default value is `""`



### create_if_deleted
* Value type is bool
* Default value is `false`



### dir_mode
* Value type is os.FileMode
* Default value is `?`



### file_mode
* Value type is os.FileMode
* Default value is `?`



### flush_interval
* Value type is time.Duration
* Default value is ``



### host
* Value type is string
* Default value is `""`



### path
* This is a required setting.
* Value type is string
* Default value is `""`



### volume
* This is a required setting.
* Value type is string
* Default value is `""`





## Configuration blueprint

```
glusterfsoutput{
	codec => ""
	create_if_deleted => bool
	dir_mode => os.FileMode
	file_mode => os.FileMode
	flush_interval => 30
	host => ""
	path => ""
	volume => ""
}
```
