+++
date = "2017-05-16T21:06:27+02:00"
description = ""
title = "Value Types"
weight = 10
+++





## Array
Example:
```
  path => [ "/var/log/messages", "/var/log/*.log" ]
  uris => [ "http://elastic.co", "http://example.net" ]
```

This example configures `path`, which is a `string` to be a array that contains an element for each of the two strings.

## Bool
A bool must be either `true` or `false`. Note that the `true` and `false` keywords are not enclosed in quotes.

Example:
```
  ssl_enable => true
```

## Hash

A hash is a collection of key value pairs specified in the format `"field1" => "value1"`.
Note that multiple key value entries are separated by spaces rather than commas.

Example:
```
match => {
  "field1" => "value1"
  "field2" => "value2"
  ...
}
```

## Int

Int must be valid numeric values (floating point or integer).

Example:
```
  port => 33
```

## String

A string must be a single character sequence. Note that string values are
enclosed in quotes, either double or single. 

Literal quotes in the string
need to be escaped with a backslash if they are of the same kind as the string delimiter, i.e. single quotes within a single-quoted string need to be escaped as well as double quotes within a double-quoted string.

Example:
```
  name => "Hello world"
  name => 'It\'s a beautiful day'
```


## Path

A path is a string that represents a valid operating system path.

Example:
```
  my_path => "/tmp/logstash"
```

## Location

Location is a "smart string", when its value is a :

* string --> well.. it will be used as a string
* system path --> the file's content will be used
* web url --> the url's raw body will be used

{{%info%}}**system path and web url** can be relative to their configuration file, even if the configuration file was used from a remote URL{{%/info%}}

Value will be parsed as a go template.

Theses functions are available to go templates


| name | description | params | examples |
| ----------   | ---------- | --------------  | -------------------------- |
| TS | format event's @timestamp with [jodaTime layout](http://joda-time.sourceforge.net/apidocs/org/joda/time/format/DateTimeFormat.html) | params | examples |
| DateFormat | format a time.Time with jodaTime layout | params | examples |
| Time | converts the textual representation of the datetime string into a time.Time | params | examples |
| Now | returns the current local time. | params | examples |
| NumFmt | format a number | params | examples |
| SafeHTML | use string as HTML | params | examples |
| HTMLUnescape | unescape a html string | params | examples |
| HTMLEscape | escape a html string | params | examples |
| Lower | lowercase a string | params | examples |
| Upper | uppercase a string | params | examples |
| Trim | trim a string | params | examples |



## Interval

Express interval and schedule processor run with a **cron spec** a **predefined schedule** or a **every <duration>** format.


### CRON Expression Format
A cron expression represents a set of times, using 6 space-separated fields.

{{%alert success%}}**tips** : have a look at [https://contrab.guru](https://crontab.guru/) {{%/alert%}}

| Field name   | Mandatory? | Allowed values  | Allowed special characters |
| ----------   | ---------- | --------------  | -------------------------- |
| Seconds      | Yes        | 0-59            | * / , - |
| Minutes      | Yes        | 0-59            | * / , - |
| Hours        | Yes        | 0-23            | * / , - |
| Day of month | Yes        | 1-31            | * / , - ? |
| Month        | Yes        | 1-12 or JAN-DEC | * / , - |
| Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ? |

Example:
```
  #interval => "0 5 4 * * *" # At 04:05
  #interval => "0 15 14 1 * *" # At 14:15 on day-of-month 1
  #interval => "0 0 22 * * 1-5" # At 22:00 on every day-of-week from Monday through Friday.
```

### Predefined schedules
You may use one of several pre-defined schedules in place of a cron expression.

| Entry                  | Description                                | Equivalent To |
| -----                  | -----------                                | ------------- |
| @yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *|
| @monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *|
| @weekly                | Run once a week, midnight on Sunday        | 0 0 0 * * 0|
| @daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *|
| @hourly                | Run once an hour, beginning of hour        | 0 0 * * * *|

Example:
```
  interval => "@hourly"
```

### Every X
You may also schedule a processor to execute at fixed intervals. This is supported by formatting the cron spec like this:

	@every <duration>

where "duration" is a string accepted by [time.ParseDuration](http://golang.org/pkg/time/#ParseDuration).

For example, "@every 1h30m10s" would indicate a schedule that activates every 1 hour, 30 minutes, 10 seconds.

Note: The interval does not take the processor runtime into account. For example, if a processor takes 3 minutes to run, and it is scheduled to run every 5 minutes, it will have only 2 minutes of idle time between each run.

Example:
```
  #interval => "@every 10s" # Every 10 seconds
```



