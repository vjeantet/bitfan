# LDAPPROCESSOR
Performs a search for a specified filter on the directory and fire events with results

## Synopsys


|      SETTING      |  TYPE  | REQUIRED |   DEFAULT VALUE   |
|-------------------|--------|----------|-------------------|
| Add_field         | hash   | false    | {}                |
| Tags              | array  | false    | []                |
| Type              | string | false    | ""                |
| host              | string | true     | ""                |
| port              | int    | true     |               389 |
| bind_dn           | string | false    | ""                |
| bind_password     | string | false    | ""                |
| base_dn           | string | true     | ""                |
| search_base       | string | false    | ""                |
| search_filter     | string | true     | "(objectClass=*)" |
| search_attributes | array  | false    | []                |
| search_scope      | string | false    | "subtree"         |
| size_limit        | int    | false    |                 0 |
| paging_size       | int    | false    |              1000 |
| event_by          | string | false    | "entry"           |
| interval          | string | false    | ""                |
| var               | hash   | false    | {}                |
| target            | string | false    | "data"            |


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

### host
* This is a required setting.
* Value type is string
* Default value is `""`

ldap hostname

### port
* This is a required setting.
* Value type is int
* Default value is `389`

ldap port

### bind_dn
* Value type is string
* Default value is `""`

Bind dn

### bind_password
* Value type is string
* Default value is `""`

Bind password

### base_dn
* This is a required setting.
* Value type is string
* Default value is `""`

Base DN
If bind_dn is not specified or is empty, an anonymous bind is attempted.
This is defined in https://tools.ietf.org/html/rfc2251#section-4.2.2

### search_base
* Value type is string
* Default value is `""`

A search base (the distinguished name of the search base object) defines the
location in the directory from which the LDAP search begins.

### search_filter
* This is a required setting.
* Value type is string
* Default value is `"(objectClass=*)"`

The search filter can be simple or advanced, using boolean operators in the format
described in the LDAP documentation (see [RFC4515](http://www.faqs.org/rfcs/rfc4515) for full information on filters).

### search_attributes
* Value type is array
* Default value is `[]`

An array of the required attributes, e.g. ["mail", "sn", "cn"].

Note that the "dn" is always returned irrespective of which attributes types are requested.

Using this parameter is much more efficient than the default action (which is to return all attributes and their associated values).

The use of this parameter should therefore be considered good practice.

### search_scope
* Value type is string
* Default value is `"subtree"`

The SCOPE setting is the starting point of an LDAP search and the depth from the
base DN to which the search should occur.

There are three options (values) that can be assigned to the SCOPE parameter:

* **base** : indicate searching only the entry at the base DN, resulting in only that entry being returned
* **one** : indicate searching all entries one level under the base DN - but not including the base DN and not including any entries under that one level under the base DN.
* **subtree** : indicate searching of all entries at all levels under and including the specified base DN

![scope](../ldapscope.gif)

### size_limit
* Value type is int
* Default value is `0`

Maximum entries to return (leave empty to let the server decide)

### paging_size
* Value type is int
* Default value is `1000`

Desired page size in order to execute LDAP queries to fulfill the
search request.

Set 0 to not use Paging

### event_by
* Value type is string
* Default value is `"entry"`

Send an event row by row or one event with all results
possible values "entry", "result"

### interval
* Value type is string
* Default value is `""`

Set an interval when this processor is used as a input

### var
* Value type is hash
* Default value is `{}`

You can set variable to be used in Search Query by using ${var}.
each reference will be replaced by the value of the variable found in search query content
The replacement is case-sensitive.

### target
* Value type is string
* Default value is `"data"`

Define the target field for placing the retrieved data. If this setting is omitted,
the data will be stored in the "data" field
Set the value to "." to store value to the root (top level) of the event



## Configuration blueprint

```
ldapprocessor{
	add_field => {}
	tags => []
	type => ""
	host => "ldap.forumsys.com"
	port => 389
	bind_dn => "cn=read-only-admin,dc=example,dc=com"
	bind_password => "password"
	base_dn => "dc=example,dc=com"
	search_base => ""
	search_filter => "(objectClass=*)"
	search_attributes => ["mail", "sn", "cn"]
	search_scope => "subtree"
	size_limit => 0
	paging_size => 1000
	event_by => "entry"
	interval => "10"
	var => {"hostname"=>"myhost","varname"=>"varvalue"}
	target => "data"
}
```
