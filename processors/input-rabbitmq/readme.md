# RABBITMQINPUT


## Synopsys


|        SETTING         |    TYPE    | REQUIRED | DEFAULT VALUE |
|------------------------|------------|----------|---------------|
| ack                    | bool       | false    | false         |
| ack_batch_size         | uint64     | false    | ?             |
| add_field              | hash       | false    | {}            |
| arguments              | amqp.Table | false    | ?             |
| auto_delete            | bool       | false    | false         |
| codec                  | string     | false    | ""            |
| connect_retry_interval | uint64     | false    | ?             |
| durable                | bool       | false    | false         |
| exchange               | string     | false    | ""            |
| exclusive              | bool       | false    | false         |
| heartbeat              | int        | false    |             0 |
| host                   | string     | false    | ""            |
| key                    | string     | false    | ""            |
| metadata_enabled       | bool       | false    | false         |
| passive                | bool       | false    | false         |
| password               | string     | false    | ""            |
| port                   | int        | false    |             0 |
| prefetch_count         | int        | false    |             0 |
| queue                  | string     | false    | ""            |
| ssl                    | bool       | false    | false         |
| tags                   | array      | false    | []            |
| user                   | string     | false    | ""            |
| verify_ssl             | bool       | false    | false         |
| vhost                  | string     | false    | ""            |


## Details

### ack
* Value type is bool
* Default value is `false`

Enable message acknowledgements. Default value is true

With acknowledgements messages fetched but not yet sent into the pipeline will be requeued by the server if BitFan shuts down.
Acknowledgements will however hurt the message throughput.
This will only send an ack back every prefetch_count messages. Working in batches provides a performance boost.

### ack_batch_size
* Value type is uint64
* Default value is `?`

Acknowledge messages in batch of value.
Default value is 1 (acknowledge each message individually)

### add_field
* Value type is hash
* Default value is `{}`

Add a field to an event. Default value is {}

### arguments
* Value type is amqp.Table
* Default value is `?`

Extra queue arguments as an array. Default value is {}

E.g. to make a RabbitMQ queue mirrored, use: {"x-ha-policy" => "all"}

### auto_delete
* Value type is bool
* Default value is `false`

Should the queue be deleted on the broker when the last consumer disconnects? Default value is false

Set this option to false if you want the queue to remain on the broker, queueing up messages until a consumer comes along to consume them.

### codec
* Value type is string
* Default value is `""`

The codec used for input data. Default value is "json"

Input codecs are a convenient method for decoding your data before it enters the input, without needing a separate filter in your BitFan pipeline.

### connect_retry_interval
* Value type is uint64
* Default value is `?`

Time in seconds to wait before retrying a connection. Default value is 1

### durable
* Value type is bool
* Default value is `false`

Is this queue durable (a.k.a "Should it survive a broker restart?"")?  Default value is false

### exchange
* Value type is string
* Default value is `""`

The name of the exchange to bind the queue to. There is no default value for this setting.

### exclusive
* Value type is bool
* Default value is `false`

Is the queue exclusive? Default value is false

Exclusive queues can only be used by the connection that declared them and will be deleted when it is closed (e.g. due to a BitFan restart).

### heartbeat
* Value type is int
* Default value is `0`

Heartbeat delay in seconds. If unspecified no heartbeats will be sent

### host
* Value type is string
* Default value is `""`

RabbitMQ server address. There is no default value for this setting.

### key
* Value type is string
* Default value is `""`

The routing key to use when binding a queue to the exchange. Default value is ""

This is only relevant for direct or topic exchanges.

### metadata_enabled
* Value type is bool
* Default value is `false`

Not implemented! Enable the storage of message headers and properties in @metadata. Default value is false

This may impact performance

### passive
* Value type is bool
* Default value is `false`

Use queue passively declared, meaning it must already exist on the server. Default value is false

To have BitFan create the queue if necessary leave this option as false.
If actively declaring a queue that already exists, the queue options for this plugin (durable etc) must match those of the existing queue.

### password
* Value type is string
* Default value is `""`

RabbitMQ password. Default value is "guest"

### port
* Value type is int
* Default value is `0`

RabbitMQ port to connect on. Default value is 5672

### prefetch_count
* Value type is int
* Default value is `0`

Prefetch count. Default value is 256

If acknowledgements are enabled with the ack option, specifies the number of outstanding unacknowledged

### queue
* Value type is string
* Default value is `""`

The name of the queue BitFan will consume events from. If left empty, a transient queue with an randomly chosen name will be created.

### ssl
* Value type is bool
* Default value is `false`

Enable or disable SSL. Default value is false

### tags
* Value type is array
* Default value is `[]`

Add any number of arbitrary tags to your event. There is no default value for this setting.

This can help with processing later. Tags can be dynamic and include parts of the event using the %{field} syntax.

### user
* Value type is string
* Default value is `""`

RabbitMQ username. Default value is "guest"

### verify_ssl
* Value type is bool
* Default value is `false`

Validate SSL certificate. Default value is false

### vhost
* Value type is string
* Default value is `""`

The vhost to use. Default value is "/"



## Configuration blueprint

```
rabbitmqinput{
	ack => bool
	ack_batch_size => uint64
	add_field => {}
	arguments => amqp.Table
	auto_delete => bool
	codec => ""
	connect_retry_interval => uint64
	durable => bool
	exchange => ""
	exclusive => bool
	heartbeat => 123
	host => ""
	key => ""
	metadata_enabled => bool
	passive => bool
	password => ""
	port => 123
	prefetch_count => 123
	queue => ""
	ssl => bool
	tags => []
	user => ""
	verify_ssl => bool
	vhost => ""
}
```
