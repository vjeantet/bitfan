# KAFKAOUTPUT


## Synopsys


|    SETTING    |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------------|--------|----------|---------------|
| brokers       | array  | false    | []            |
| topic         | string | true     | ""            |
| balancer      | string | false    | ""            |
| max_attempts  | int    | false    |             0 |
| queue_size    | int    | false    |             0 |
| batch_size    | int    | false    |             0 |
| keepalive     | int    | false    |             0 |
| io_timeout    | int    | false    |             0 |
| required_acks | int    | false    |             0 |


## Details

### brokers
* Value type is array
* Default value is `[]`

Broker list

### topic
* This is a required setting.
* Value type is string
* Default value is `""`

Kafka topic

### balancer
* Value type is string
* Default value is `""`

Balancer ( roundrobin, hash or leastbytes )

### max_attempts
* Value type is int
* Default value is `0`

Max Attempts

### queue_size
* Value type is int
* Default value is `0`

Queue Size

### batch_size
* Value type is int
* Default value is `0`

Batch Size

### keepalive
* Value type is int
* Default value is `0`

Keep Alive ( in seconds )

### io_timeout
* Value type is int
* Default value is `0`

IO Timeout ( in seconds )

### required_acks
* Value type is int
* Default value is `0`

Required Acks ( number of replicas that must acknowledge write. -1 for all replicas )



## Configuration blueprint

```
kafkaoutput{
	brokers => []
	topic => ""
	balancer => ""
	max_attempts => 123
	queue_size => 123
	batch_size => 123
	keepalive => 123
	io_timeout => 123
	required_acks => 123
}
```
