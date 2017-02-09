# STATSD


## Synopsys


|   SETTING   |  TYPE   | REQUIRED | DEFAULT VALUE |
|-------------|---------|----------|---------------|
| host        | string  | false    | ""            |
| port        | int     | false    |             0 |
| protocol    | string  | false    | ""            |
| sender      | string  | false    | ""            |
| count       | hash    | false    | {}            |
| decrement   | array   | false    | []            |
| gauge       | hash    | false    | {}            |
| increment   | array   | false    | []            |
| namespace   | string  | false    | ""            |
| sample_rate | float32 | false    | ?             |
| set         | hash    | false    | {}            |
| timing      | hash    | false    | {}            |


## Details

### host
* Value type is string
* Default value is `""`

The address of the statsd server.

### port
* Value type is int
* Default value is `0`

The port to connect to on your statsd server.

### protocol
* Value type is string
* Default value is `""`



### sender
* Value type is string
* Default value is `""`

The name of the sender. Dots will be replaced with underscores.

### count
* Value type is hash
* Default value is `{}`

A count metric. metric_name => count as hash

### decrement
* Value type is array
* Default value is `[]`

A decrement metric. Metric names as array.

### gauge
* Value type is hash
* Default value is `{}`

A gauge metric. metric_name => gauge as hash.

### increment
* Value type is array
* Default value is `[]`

An increment metric. Metric names as array.

### namespace
* Value type is string
* Default value is `""`

The statsd namespace to use for this metric.

### sample_rate
* Value type is float32
* Default value is `?`

The sample rate for the metric.

### set
* Value type is hash
* Default value is `{}`

A set metric. metric_name => "string" to append as hash

### timing
* Value type is hash
* Default value is `{}`

A timing metric. metric_name => duration as hash



## Configuration blueprint

```
statsd{
	host => ""
	port => 123
	protocol => ""
	sender => ""
	count => {}
	decrement => []
	gauge => {}
	increment => []
	namespace => ""
	sample_rate => float32
	set => {}
	timing => {}
}
```
