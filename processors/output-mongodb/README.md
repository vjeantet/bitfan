# MONGODB


## Synopsys


|   SETTING   |  TYPE  | REQUIRED | DEFAULT VALUE |
|-------------|--------|----------|---------------|
| Collection  | string | false    | ""            |
| Database    | string | false    | ""            |
| GenerateId  | bool   | false    | false         |
| Isodate     | bool   | false    | false         |
| Retry_delay | int    | false    |             0 |
| Uri         | string | false    | ""            |


## Details

### Collection
* Value type is string
* Default value is `""`

The collection to use. This value can use %{foo} values to dynamically
select a collection based on data in the event

### Database
* Value type is string
* Default value is `""`

The database to use

### GenerateId
* Value type is bool
* Default value is `false`

If true, an "_id" field will be added to the document before insertion.
The "_id" field will use the timestamp of the event and overwrite an
existing "_id" field in the event

### Isodate
* Value type is bool
* Default value is `false`

If true, store the @timestamp field in mongodb as an ISODate type
instead of an ISO8601 string. For more information about this,
see http://www.mongodb.org/display/DOCS/Dates

### Retry_delay
* Value type is int
* Default value is `0`

Number of seconds to wait after failure before retrying

### Uri
* Value type is string
* Default value is `""`

a MongoDB URI to connect to See http://docs.mongodb.org/manual/reference/connection-string/



## Configuration blueprint

```
mongodb{
	collection => ""
	database => ""
	generateid => bool
	isodate => bool
	retry_delay => 123
	uri => ""
}
```
