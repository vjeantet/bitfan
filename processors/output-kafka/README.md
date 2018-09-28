# KAFKAOUTPUT


## Synopsys


|      SETTING      |  TYPE  | REQUIRED | DEFAULT VALUE |
|-------------------|--------|----------|---------------|
| bootstrap_servers | string | false    | ""            |
| brokers           | array  | false    | []            |
| topic_id          | string | true     | ""            |
| client_id         | string | false    | ""            |
| balancer          | string | false    | ""            |
| max_attempts      | int    | false    |             0 |
| queue_size        | int    | false    |             0 |
| batch_size        | int    | false    |             0 |
| keepalive         | int    | false    |             0 |
| io_timeout        | int    | false    |             0 |
| acks              | int    | false    |             0 |


## Details

### bootstrap_servers
* Value type is string
* Default value is `""`

Bootstrap Servers ( "host:port" )

### brokers
* Value type is array
* Default value is `[]`

Broker list

### topic_id
* This is a required setting.
* Value type is string
* Default value is `""`

Kafka topic

### client_id
* Value type is string
* Default value is `""`

Kafka client id

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

### acks
* Value type is int
* Default value is `0`

Required Acks ( number of replicas that must acknowledge write. -1 for all replicas )



## Configuration blueprint

```
kafkaoutput{
	bootstrap_servers => ""
	brokers => []
	topic_id => ""
	client_id => ""
	balancer => ""
	max_attempts => 123
	queue_size => 123
	batch_size => 123
	keepalive => 123
	io_timeout => 123
	acks => 123
}
```
