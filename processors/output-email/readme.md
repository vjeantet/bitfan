# EMAIL
Send email when an output is received. Alternatively, you may include or exclude the email output execution using conditionals.

## Synopsys


|     SETTING      |   TYPE   | REQUIRED |    DEFAULT VALUE     |
|------------------|----------|----------|----------------------|
| address          | string   | false    | "localhost"          |
| port             | int      | false    |                   25 |
| username         | string   | false    | ""                   |
| password         | string   | false    | ""                   |
| from             | string   | false    | "bitfan@nowhere.com" |
| replyto          | string   | false    | ""                   |
| to               | string   | true     | ""                   |
| cc               | string   | false    | ""                   |
| bcc              | string   | false    | ""                   |
| subject          | string   | false    | ""                   |
| subjectfile      | string   | false    | ""                   |
| htmlbody         | location | false    | ?                    |
| body             | location | false    | ?                    |
| attachments      | array    | false    | []                   |
| images           | array    | false    | []                   |
| embed_b64_images | bool     | false    | false                |


## Details

### address
* Value type is string
* Default value is `"localhost"`

The address used to connect to the mail server

### port
* Value type is int
* Default value is `25`

Port used to communicate with the mail server

### username
* Value type is string
* Default value is `""`

Username to authenticate with the server

### password
* Value type is string
* Default value is `""`

Password to authenticate with the server

### from
* Value type is string
* Default value is `"bitfan@nowhere.com"`

The fully-qualified email address for the From: field in the email

### replyto
* Value type is string
* Default value is `""`

The fully qualified email address for the Reply-To: field

### to
* This is a required setting.
* Value type is string
* Default value is `""`

The fully-qualified email address to send the email to.

This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`

You can also use dynamic fields from the event with the %{fieldname} syntax

### cc
* Value type is string
* Default value is `""`

The fully-qualified email address(es) to include as cc: address(es).

This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`

### bcc
* Value type is string
* Default value is `""`

The fully-qualified email address(es) to include as bcc: address(es).

This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`

### subject
* Value type is string
* Default value is `""`

Subject: for the email

You can use template

### subjectfile
* Value type is string
* Default value is `""`

Path to Subject template file for the email

### htmlbody
* Value type is location
* Default value is `?`

HTML Body for the email, which may contain HTML markup

### body
* Value type is location
* Default value is `?`

Body for the email - plain text only.

### attachments
* Value type is array
* Default value is `[]`

Attachments - specify the name(s) and location(s) of the files

### images
* Value type is array
* Default value is `[]`

Images - specify the name(s) and location(s) of the images

### embed_b64_images
* Value type is bool
* Default value is `false`

Search for img:data in HTML body, and replace them to a reference to inline attachment



## Configuration blueprint

```
email{
	address => "localhost"
	port => 25
	username => ""
	password => ""
	from => "bitfan@nowhere.com"
	replyto => "test@nowhere.com"
	to => "me@host.com, you@host.com"
	cc => "me@host.com, you@host.com"
	bcc => "me@host.com, you@host.com"
	subject => "message from {{.host}}"
	subjectfile => ""
	htmlBody => "<h1>Hello</h1> message received : {{.message}}"
	body => "message : {{.message}}. from {{.host}}."
	attachments => []
	images => []
	embed_b64_images => false
}
```
