# RABBITMQOUTPUT


## Synopsys


|        SETTING         |     TYPE      | REQUIRED | DEFAULT VALUE |
|------------------------|---------------|----------|---------------|
| add_field              | hash          | false    | {}            |
| arguments              | amqp.Table    | false    | ?             |
| connect_retry_interval | time.Duration | false    |               |
| connection_timeout     | time.Duration | false    |               |
| durable                | bool          | false    | ?             |
| exchange               | string        | true     | ""            |
| exchange_type          | string        | true     | ""            |
| heartbeat              | time.Duration | false    |               |
| host                   | string        | false    | ""            |
| key                    | string        | false    | ""            |
| passive                | bool          | false    | ?             |
| password               | string        | false    | ""            |
| persistent             | bool          | false    | ?             |
| port                   | int           | false    |             0 |
| ssl                    | bool          | false    | ?             |
| tags                   | array         | false    | []            |
| user                   | string        | false    | ""            |
| verify_ssl             | bool          | false    | ?             |
| vhost                  | string        | false    | ""            |


## Details

### add_field
* Value type is hash
* Default value is `{}`

Add a field to an event. Default value is {}

### arguments
* Value type is amqp.Table
* Default value is `?`

Extra exchange arguments. Default value is {}

### connect_retry_interval
* Value type is time.Duration
* Default value is ``

Time in seconds to wait before retrying a connection. Default value is 1

### connection_timeout
* Value type is time.Duration
* Default value is ``

Time in seconds to wait before timing-out. Default value is 0 (no timeout)

### durable
* Value type is bool
* Default value is `?`

Is this exchange durable - should it survive a broker restart? Default value is true

### exchange
* This is a required setting.
* Value type is string
* Default value is `""`

The name of the exchange to send message to. There is no default value for this setting.

### exchange_type
* This is a required setting.
* Value type is string
* Default value is `""`

The exchange type (fanout, topic, direct). There is no default value for this setting.

### heartbeat
* Value type is time.Duration
* Default value is ``

Interval (in second) to send heartbeat to rabbitmq. Default value is 0
If value if lower than 1, server's interval setting will be used.

### host
* Value type is string
* Default value is `""`

RabbitMQ server address. There is no default value for this setting.

### key
* Value type is string
* Default value is `""`

The routing key to use when binding a queue to the exchange. Default value is ""
This is only relevant for direct or topic exchanges (Routing keys are ignored on fanout exchanges).
This setting can be dynamic using the %{foo} syntax.

### passive
* Value type is bool
* Default value is `?`

Use queue passively declared, meaning it must already exist on the server. Default value is false
To have BitFan to create the queue if necessary leave this option as false.
If actively declaring a queue that already exists, the queue options for this plugin (durable, etc) must match those of the existing queue.

### password
* Value type is string
* Default value is `""`

RabbitMQ password. Default value is "guest"

### persistent
* Value type is bool
* Default value is `?`

Should RabbitMQ persist messages to disk? Default value is true

### port
* Value type is int
* Default value is `0`

RabbitMQ port to connect on. Default value is 5672

### ssl
* Value type is bool
* Default value is `?`

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
* Default value is `?`

Validate SSL certificate. Default value is false

### vhost
* Value type is string
* Default value is `""`

The vhost to use. Default value is "/"



## Configuration blueprint

```
rabbitmqoutput{
	add_field => {}
	arguments => amqp.Table
	connect_retry_interval => 30
	connection_timeout => 30
	durable => bool
	exchange => ""
	exchange_type => ""
	heartbeat => 30
	host => ""
	key => ""
	passive => bool
	password => ""
	persistent => bool
	port => 123
	ssl => bool
	tags => []
	user => ""
	verify_ssl => bool
	vhost => ""
}
```
