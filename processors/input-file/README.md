# FILE
Read file on

* received event
* when new file discovered

this processor remember last files used, it stores references in sincedb, set it to "/dev/null" to not remember used files

## Synopsys


|      SETTING      |  TYPE  | REQUIRED |             DEFAULT VALUE              |
|-------------------|--------|----------|----------------------------------------|
| codec             | codec  | false    | "plain"                                |
| read_older        | int    | false    |                                      0 |
| discover_interval | int    | false    |                                     15 |
| exclude           | array  | false    | []                                     |
| ignore_older      | int    | false    |                                      0 |
| max_open_files    | int    | false    |                                      0 |
| path              | array  | true     | []                                     |
| sincedb_path      | string | false    | :                                      |
|                   |        |          | "$dataLocation/readfile/.sincedb.json" |
| target            | string | false    | ""                                     |


## Details

### codec
* Value type is codec
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
* Default value is `15`

How often (in seconds) we expand the filename patterns in the path option
to discover new files to watch. Default value is 15
When value is 0, processor will read file, one time, on start.

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
* Default value is `: "$dataLocation/readfile/.sincedb.json"`

Path of the sincedb database file
The sincedb database keeps track of the current position of monitored
log files that will be written to disk.
Set it to "/dev/null" to not use sincedb features

### target
* Value type is string
* Default value is `""`

When decoded data is an array it stores the resulting data into the given target field.



## Configuration blueprint

```
file{
	codec => "plain"
	read_older => 123
	discover_interval => 15
	exclude => []
	ignore_older => 123
	max_open_files => 123
	path => []
	: sincedb_path => "/dev/null"
	target => ""
}
```
