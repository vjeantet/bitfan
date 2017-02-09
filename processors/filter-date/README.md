# DATE
The date filter is used for parsing dates from fields, and then using that date or timestamp as the logstash timestamp for the event.
For example, syslog events usually have timestamps like this:
`"Apr 17 09:32:01"`
You would use the date format MMM dd HH:mm:ss to parse this.
The date filter is especially important for sorting events and for backfilling old data. If you don’t get the date correct in your event, then searching for them later will likely sort out of order.
In the absence of this filter, logstash will choose a timestamp based on the first time it sees the event (at input time), if the timestamp is not already set in the event. For example, with file input, the timestamp is set to the time of each read.

## Synopsys


|    SETTING     |  TYPE  | REQUIRED | DEFAULT VALUE |
|----------------|--------|----------|---------------|
| add_field      | hash   | false    | {}            |
| match          | array  | false    | []            |
| add_tag        | array  | false    | []            |
| remove_field   | array  | false    | []            |
| remove_tag     | array  | false    | []            |
| tag_on_failure | array  | false    | []            |
| target         | string | false    | ""            |
| timezone       | string | false    | ""            |


## Details

### add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### match
* Value type is array
* Default value is `[]`

The date formats allowed are anything allowed by Golang time format.
You can see the docs for this format https://golang.org/src/time/format.go#L20
An array with field name first, and format patterns following, [ field, formats... ]

### add_tag
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.

### remove_field
* Value type is array
* Default value is `[]`

If this filter is successful, remove arbitrary fields from this event.

### remove_tag
* Value type is array
* Default value is `[]`



### tag_on_failure
* Value type is array
* Default value is `[]`

Append values to the tags field when there has been no successful match
Default value is ["_dateparsefailure"]

### target
* Value type is string
* Default value is `""`

Store the matching timestamp into the given target field. If not provided,
default to updating the @timestamp field of the event

### timezone
* Value type is string
* Default value is `""`

Specify a time zone canonical ID to be used for date parsing.
The valid IDs are listed on IANA Time Zone database, such as "America/New_York".
This is useful in case the time zone cannot be extracted from the value,
and is not the platform default. If this is not specified the platform default
 will be used. Canonical ID is good as it takes care of daylight saving time
for you For example, America/Los_Angeles or Europe/Paris are valid IDs.
This field can be dynamic and include parts of the event using the %{field} syntax



## Configuration blueprint

```
date{
	add_field => {}
	match => []
	add_tag => []
	remove_field => []
	remove_tag => []
	tag_on_failure => []
	target => ""
	timezone => ""
}
```
