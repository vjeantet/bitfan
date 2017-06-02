# FILE


## Synopsys


|      SETTING      |  TYPE  | REQUIRED | DEFAULT VALUE |
|-------------------|--------|----------|---------------|
| Add_field         | hash   | false    | {}            |
| Tags              | array  | false    | []            |
| Type              | string | false    | ""            |
| Codec             | string | false    | "plain"       |
| read_older        | int    | false    |             0 |
| discover_interval | int    | false    |             0 |
| exclude           | array  | false    | []            |
| ignore_older      | int    | false    |             0 |
| max_open_files    | int    | false    |             0 |
| path              | array  | true     | []            |
| sincedb_path      | string | false    | ""            |


## Details

### Add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### Tags
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.

### Type
* Value type is string
* Default value is `""`

Add a type field to all events handled by this input

### Codec
* Value type is string
* Default value is `"plain"`

The codec used for input data. Input codecs are a convenient method for decoding
your data before it enters the input, without needing a separate filter in your bitfan pipeline

### read_older
* Value type is int
* Default value is `0`

How many seconds a file should stay unmodified to be read
use this to prevent reading a file while another process is writing into.

### discover_interval
* Value type is int
* Default value is `0`

How often (in seconds) we expand the filename patterns in the path option
to discover new files to watch. Default value is 15

### exclude
* Value type is array
* Default value is `[]`

Exclusions (matched against the filename, not full path).
Filename patterns are valid here, too.

### ignore_older
* Value type is int
* Default value is `0`

When the file input discovers a file that was last modified before the
specified timespan in seconds, the file is ignored.
After it’s discovery, if an ignored file is modified it is no longer ignored
and any new data is read.
Default value is 86400 (i.e. 24 hours)

### max_open_files
* Value type is int
* Default value is `0`

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
* Default value is `""`

Path of the sincedb database file
The sincedb database keeps track of the current position of monitored
log files that will be written to disk.



## Configuration blueprint

```
file{
	add_field => {}
	tags => []
	type => ""
	codec => "plain"
	read_older => 123
	discover_interval => 123
	exclude => []
	ignore_older => 123
	max_open_files => 123
	path => []
	sincedb_path => ""
}
```
