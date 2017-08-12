# BEATSINPUT


## Synopsys


|           SETTING           |  TYPE  | REQUIRED | DEFAULT VALUE |
|-----------------------------|--------|----------|---------------|
| Congestion_threshold        | int    | false    |             0 |
| Host                        | string | false    | ""            |
| Port                        | int    | false    |             0 |
| Ssl                         | bool   | false    | ?             |
| Ssl_certificate             | string | false    | ""            |
| Ssl_certificate_authorities | array  | false    | []            |
| Ssl_key                     | string | false    | ""            |
| Ssl_key_passphrase          | string | false    | ""            |
| Ssl_verify_mode             | string | false    | ""            |


## Details

### Congestion_threshold
* Value type is int
* Default value is `0`

The number of seconds before we raise a timeout,
this option is useful to control how much time to wait if something is blocking
the pipeline

### Host
* Value type is string
* Default value is `""`

The IP address to listen on

### Port
* Value type is int
* Default value is `0`

The port to listen on (default 5044)

### Ssl
* Value type is bool
* Default value is `?`

Events are by default send in plain text,
you can enable encryption by using ssl to true and
configuring the ssl_certificate and ssl_key options

### Ssl_certificate
* Value type is string
* Default value is `""`

SSL certificate to use (path)

### Ssl_certificate_authorities
* Value type is array
* Default value is `[]`

Validate client certificates against theses authorities
 You can defined multiples files or path, all the certificates will be read
 and added to the trust store.
 You need to configure the ssl_verify_mode to peer or force_peer to enable
 the verification.
This feature only support certificate directly signed by your root ca.
Intermediate CA are currently not supported.

### Ssl_key
* Value type is string
* Default value is `""`

SSL key to use (path)

### Ssl_key_passphrase
* Value type is string
* Default value is `""`

SSL key passphrase to use. (not yet implemented)

### Ssl_verify_mode
* Value type is string
* Default value is `""`

By default the server dont do any client verification,
peer will make the server ask the client to provide a certificate,
  if the client provide the certificate it will be validated.
force_peer will make the server ask the client for their certificate,
  if the clients doesn’t provide it the connection will be closed.
This option need to be used with ssl_certificate_authorities and a defined list of CA.
Value can be any of: none, peer, force_peer



## Configuration blueprint

```
beatsinput{
	congestion_threshold => 123
	host => ""
	port => 123
	ssl => bool
	ssl_certificate => ""
	ssl_certificate_authorities => []
	ssl_key => ""
	ssl_key_passphrase => ""
	ssl_verify_mode => ""
}
```
