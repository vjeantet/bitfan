# ELASTICSEARCH2


## Synopsys


|     SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------------|--------|----------|---------------|
| document_type   | string | false    | ""            |
| flush_count     | int    | false    |             0 |
| flush_size      | int    | false    |             0 |
| host            | string | false    | ""            |
| idle_flush_time | int    | false    |             0 |
| index           | string | false    | ""            |
| password        | string | false    | ""            |
| path            | string | false    | ""            |
| port            | int    | false    |             0 |
| user            | string | false    | ""            |
| ssl             | bool   | false    | ?             |


## Details

### document_type
* Value type is string
* Default value is `""`

The document type to write events to. There is no default value for this setting.

Generally you should try to write only similar events to the same type.
String expansion %{foo} works here. Unless you set document_type, the event type will
be used if it exists otherwise the document type will be assigned the value of logs

### flush_count
* Value type is int
* Default value is `0`

The number of requests that can be enqueued before flushing them. Default value is 1000

### flush_size
* Value type is int
* Default value is `0`

The number of bytes that the bulk requests can take up before the bulk processor decides to flush. Default value is 5242880 (5MB).

### host
* Value type is string
* Default value is `""`

Host of the remote instance. Default value is "localhost"

### idle_flush_time
* Value type is int
* Default value is `0`

The amount of seconds since last flush before a flush is forced. Default value is 1

This setting helps ensure slow event rates don’t get stuck.
For example, if your flush_size is 100, and you have received 10 events,
and it has been more than idle_flush_time seconds since the last flush,
those 10 events will be flushed automatically.
This helps keep both fast and slow log streams moving along in near-real-time.

### index
* Value type is string
* Default value is `""`

The index to write events to. Default value is "logstash-%Y.%m.%d"

This can be dynamic using the %{foo} syntax and strftime syntax (see http://strftime.org/).
The default value will partition your indices by day.

### password
* Value type is string
* Default value is `""`

Password to authenticate to a secure Elasticsearch cluster. There is no default value for this setting.

### path
* Value type is string
* Default value is `""`

HTTP Path at which the Elasticsearch server lives. Default value is "/"

Use this if you must run Elasticsearch behind a proxy that remaps the root path for the Elasticsearch HTTP API lives.

### port
* Value type is int
* Default value is `0`

ElasticSearch port to connect on. Default value is 9200

### user
* Value type is string
* Default value is `""`

Username to authenticate to a secure Elasticsearch cluster. There is no default value for this setting.

### ssl
* Value type is bool
* Default value is `?`

Enable SSL/TLS secured communication to Elasticsearch cluster. Default value is false



## Configuration blueprint

```
elasticsearch2{
	document_type => ""
	flush_count => 123
	flush_size => 123
	host => ""
	idle_flush_time => 123
	index => ""
	password => ""
	path => ""
	port => 123
	user => ""
	ssl => bool
}
```
