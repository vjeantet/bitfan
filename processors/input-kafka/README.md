# KAFKAINPUT


## Synopsys


|      SETTING      |  TYPE  | REQUIRED | DEFAULT VALUE |
|-------------------|--------|----------|---------------|
| bootstrap_server  | string | false    | ""            |
| brokers           | array  | false    | []            |
| topic_id          | string | true     | ""            |
| group_id          | string | false    | ""            |
| client_id         | string | false    | ""            |
| balancer          | string | false    | ""            |
| compression       | string | false    | ""            |
| max_attempts      | int    | false    |             0 |
| queue_size        | int    | false    |             0 |
| request_bytes_min | int    | false    |             0 |
| request_bytes_max | int    | false    |             0 |
| keepalive         | int    | false    |             0 |
| max_wait          | int    | false    |             0 |
| read_lag_interval | int    | false    |             0 |


## Details

### bootstrap_server
* Value type is string
* Default value is `""`

Bootstrap Server ( "host:port" )

### brokers
* Value type is array
* Default value is `[]`

Broker list

### topic_id
* This is a required setting.
* Value type is string
* Default value is `""`

Kafka topic

### group_id
* Value type is string
* Default value is `""`

Kafka group id

### client_id
* Value type is string
* Default value is `""`

Kafka client id

### balancer
* Value type is string
* Default value is `""`

Balancer ( roundrobin, hash or leastbytes )

### compression
* Value type is string
* Default value is `""`

Compression algorithm ( 'gzip', 'snappy', or 'lz4' )

### max_attempts
* Value type is int
* Default value is `0`

Max Attempts

### queue_size
* Value type is int
* Default value is `0`

Queue Size

### request_bytes_min
* Value type is int
* Default value is `0`

Minimum amount of bytes to fetch per request

### request_bytes_max
* Value type is int
* Default value is `0`

Maximum amount of bytes to fetch per request

### keepalive
* Value type is int
* Default value is `0`

Keep Alive ( in seconds )

### max_wait
* Value type is int
* Default value is `0`

Max time to wait for new data when fetching batches ( in seconds )

### read_lag_interval
* Value type is int
* Default value is `0`

Frequency at which the reader lag is updated. Negative value disables lag reporting.



## Configuration blueprint

```
kafkainput{
	bootstrap_server => ""
	brokers => []
	topic_id => ""
	group_id => ""
	client_id => ""
	balancer => ""
	compression => ""
	max_attempts => 123
	queue_size => 123
	request_bytes_min => 123
	request_bytes_max => 123
	keepalive => 123
	max_wait => 123
	read_lag_interval => 123
}
```
