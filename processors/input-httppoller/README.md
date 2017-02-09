# HTTPPOLLER
HTTPPoller allows you to call an HTTP Endpoint, decode the output of it into an event

## Synopsys


| SETTING |  TYPE  | REQUIRED | DEFAULT VALUE |
|---------|--------|----------|---------------|
| Method  | string | false    | ""            |
| Url     | string | false    | ""            |


## Details

### Method
* Value type is string
* Default value is `""`



### Url
* Value type is string
* Default value is `""`





## Configuration blueprint

```
httppoller{
	method => ""
	url => ""
}
```
