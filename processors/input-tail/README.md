# TAIL


## Synopsys


|        SETTING         |  TYPE  | REQUIRED |  DEFAULT VALUE  |
|------------------------|--------|----------|-----------------|
| add_field              | hash   | false    | {}              |
| close_older            | int    | false    |            3600 |
| codec                  | codec  | false    | ?               |
| delimiter              | string | false    | "\n"            |
| discover_interval      | int    | false    |              15 |
| exclude                | array  | false    | []              |
| ignore_older           | int    | false    |           86400 |
| max_open_files         | string | false    | ""              |
| path                   | array  | true     | []              |
| sincedb_path           | string | false    | ".sincedb.json" |
| sincedb_write_interval | int    | false    |              15 |
| start_position         | string | false    | "end"           |
| stat_interval          | int    | false    |               1 |
| tags                   | array  | false    | []              |
| type                   | string | false    | ""              |


## Details

### add_field
* Value type is hash
* Default value is `{}`

Add a field to an event. Default value is {}

### close_older
* Value type is int
* Default value is `3600`

Closes any files that were last read the specified timespan in seconds ago.
Default value is 3600 (i.e. 1 hour)
This has different implications depending on if a file is being tailed or read.
If tailing, and there is a large time gap in incoming data the file can be
closed (allowing other files to be opened) but will be queued for reopening
when new data is detected. If reading, the file will be closed after
close_older seconds from when the last bytes were read.

### codec
* Value type is codec
* Default value is `?`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

### delimiter
* Value type is string
* Default value is `"\n"`

Set the new line delimiter. Default value is "\n"

### discover_interval
* Value type is int
* Default value is `15`

How often (in seconds) we expand the filename patterns in the path option
to discover new files to watch. Default value is 15

### exclude
* Value type is array
* Default value is `[]`

Exclusions (matched against the filename, not full path).
Filename patterns are valid here, too.

### ignore_older
* Value type is int
* Default value is `86400`

When the file input discovers a file that was last modified before the
specified timespan in seconds, the file is ignored.
After it’s discovery, if an ignored file is modified it is no longer ignored
and any new data is read.
Default value is 86400 (i.e. 24 hours)

### max_open_files
* Value type is string
* Default value is `""`

What is the maximum number of file_handles that this input consumes at any one time.
Use close_older to close some files if you need to process more files than this number.

### path
* This is a required setting.
* Value type is array
* Default value is `[]`

The path(s) to the file(s) to use as an input.
You can use filename patterns here, such as /var/log/*.log.
If you use a pattern like /var/log/**/*.log, a recursive search of /var/log
will be done for all *.log files.
Paths must be absolute and cannot be relative.
You may also configure multiple paths.

### sincedb_path
* Value type is string
* Default value is `".sincedb.json"`

Path of the sincedb database file
The sincedb database keeps track of the current position of monitored
log files that will be written to disk.

### sincedb_write_interval
* Value type is int
* Default value is `15`

How often (in seconds) to write a since database with the current position of monitored log files.
Default value is 15

### start_position
* Value type is string
* Default value is `"end"`

Choose where BitFan starts initially reading files: at the beginning or at the end.
The default behavior treats files like live streams and thus starts at the end.
If you have old data you want to import, set this to beginning.
This option only modifies "first contact" situations where a file is new
and not seen before, i.e. files that don’t have a current position recorded in a sincedb file.
If a file has already been seen before, this option has no effect and the
position recorded in the sincedb file will be used.
Default value is "end"
Value can be any of: "beginning", "end"

### stat_interval
* Value type is int
* Default value is `1`

How often (in seconds) we stat files to see if they have been modified.
Increasing this interval will decrease the number of system calls we make,
but increase the time to detect new log lines.
Default value is 1

### tags
* Value type is array
* Default value is `[]`

Add any number of arbitrary tags to your event. There is no default value for this setting.
This can help with processing later. Tags can be dynamic and include parts of the event using the %{field} syntax.

### type
* Value type is string
* Default value is `""`

Add a type field to all events handled by this input.
Types are used mainly for filter activation.



## Configuration blueprint

```
tail{
	add_field => {}
	close_older => 3600
	codec => codec
	delimiter => "\n"
	discover_interval => 15
	exclude => []
	ignore_older => 86400
	max_open_files => ""
	path => []
	sincedb_path => ".sincedb.json"
	sincedb_write_interval => 15
	start_position => "end"
	stat_interval => 1
	tags => []
	type => ""
}
```
