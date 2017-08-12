# GEOIP


## Synopsys


|     SETTING      |     TYPE      | REQUIRED | DEFAULT VALUE |
|------------------|---------------|----------|---------------|
| database         | string        | false    | ""            |
| database_type    | string        | false    | ""            |
| refresh_interval | time.Duration | false    |               |
| fields           | array         | false    | []            |
| lru_cache_size   | int64         | false    |             0 |
| source           | string        | true     | ""            |
| tag_on_failure   | array         | false    | []            |
| target           | string        | false    | ""            |
| language         | string        | false    | ""            |


## Details

### database
* Value type is string
* Default value is `""`

Path or URL to the MaxMind GeoIP2 database.
Default value is "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"
Note that URL can point to gzipped database (*.mmdb.gz) but local path must point to an unzipped file.

### database_type
* Value type is string
* Default value is `""`

Type of GeoIP database. Default value is "city"
Must be one of "city", "asn", "isp" or "organization".

### refresh_interval
* Value type is time.Duration
* Default value is ``

GeoIP database refresh interval, in minutes. Default value is 60
If `database` field is a path, file will be reloaded from disk.
If it is an URL, database will be fetched (if ETAG differs) and reloaded.

### fields
* Value type is array
* Default value is `[]`

An array of geoip fields to be included in the event.
Possible fields depend on the database type. By default, all geoip fields are included in the event.

### lru_cache_size
* Value type is int64
* Default value is `0`

Cache size. Default value is 1000

### source
* This is a required setting.
* Value type is string
* Default value is `""`

The field containing the IP address or hostname to map via geoip.

### tag_on_failure
* Value type is array
* Default value is `[]`

Append values to the tags field when there has been no successful match
Default value is ["_geoipparsefailure"]

### target
* Value type is string
* Default value is `""`

Define the target field for placing the parsed data. If this setting is omitted,
the geoip data will be stored at the root (top level) of the event

### language
* Value type is string
* Default value is `""`

Language to use for city/region/continent names



## Configuration blueprint

```
geoip{
	database => ""
	database_type => ""
	refresh_interval => 30
	fields => []
	lru_cache_size => 123
	source => ""
	tag_on_failure => []
	target => ""
	language => ""
}
```
