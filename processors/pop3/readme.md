# POP3PROCESSOR
Periodically scan an POP3 mailbox for new emails.

## Synopsys


|      SETTING      |   TYPE   | REQUIRED |  DEFAULT VALUE  |
|-------------------|----------|----------|-----------------|
| interval          | interval | false    | ?               |
| host              | string   | true     | ""              |
| port              | int      | false    |             110 |
| secure            | bool     | false    | false           |
| username          | string   | true     | ""              |
| password          | string   | true     | ""              |
| dial_timeout      | int      | false    |              30 |
| delete            | bool     | false    | true            |
| strip_attachments | bool     | false    | false           |
| verify_cert       | bool     | false    | true            |
| sincedb_path      | string   | false    | : Host@Username |
| add_raw_message   | bool     | false    | false           |
| add_all_headers   | bool     | false    | false           |


## Details

### interval
* Value type is interval
* Default value is `?`

When new mail should be retreived from POP3 server ?
Nothing by default, as this processor can be used in filter

### host
* This is a required setting.
* Value type is string
* Default value is `""`

POP3 host name

### port
* Value type is int
* Default value is `110`

POP3 server's port.

When empty and secure is true (pop3s) the default port number is 995

### secure
* Value type is bool
* Default value is `false`

Use TLS POP3S connexion with server.
The default pop3s port is 995 in this case

### username
* This is a required setting.
* Value type is string
* Default value is `""`

POP3 mailbox Username

### password
* This is a required setting.
* Value type is string
* Default value is `""`

POP3 mailbox Password
you may use an env variable to pass value, like password => "${BITFAN_POP3_PASSWORD}"

### dial_timeout
* Value type is int
* Default value is `30`

How long to wait for the server to respond ?
(in second)

### delete
* Value type is bool
* Default value is `true`

Should delete message after retreiving it ?

When false, this processor will use sinceDB to not retreive an already seen message

### strip_attachments
* Value type is bool
* Default value is `false`

Add Attachements, Inlines, in the produced event ?
When false Parts are added like
```
 "parts": {
  {
    "Size":        336303,
    "Content":     $$ContentAsBytes$$,
    "Type":        "inline",
    "ContentType": "image/png",
    "Disposition": "inline",
    "FileName":    "Capture d’écran 2018-01-12 à 12.11.52.png",
  },
  {
    "Content":     $$ContentAsBytes$$,
    "Type":        "attachement",
    "ContentType": "application/pdf",
    "Disposition": "attachment",
    "FileName":    "58831639.pdf",
    "Size":        14962,
  },
},
```

### verify_cert
* Value type is bool
* Default value is `true`

When using a secure pop connexion (POP3S) should server'cert be verified ?

### sincedb_path
* Value type is string
* Default value is `: Host@Username`

Path of the sincedb database file

The sincedb database keeps track of the last seen message

Set it to `"/dev/null"` to not persist sincedb features

Tracks are done by host and username combination, you can customize this if needed giving a specific path

### add_raw_message
* Value type is bool
* Default value is `false`

Add a field to event with the raw message data ?

### add_all_headers
* Value type is bool
* Default value is `false`

Add a field to event with all headers as hash ?



## Configuration blueprint

```
pop3processor{
	interval => interval
	host => ""
	port => 110
	secure => false
	username => ""
	password => ""
	dial_timeout => 30
	delete => true
	strip_attachments => false
	verify_cert => true
	: sincedb_path => "/dev/null"
	add_raw_message => false
	add_all_headers => false
}
```
