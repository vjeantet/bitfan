# TWITTER


## Synopsys


|      SETTING       |  TYPE  | REQUIRED | DEFAULT VALUE |
|--------------------|--------|----------|---------------|
| Add_field          | hash   | false    | {}            |
| Consumer_key       | string | false    | ""            |
| Consumer_secret    | string | false    | ""            |
| Oauth_token        | string | false    | ""            |
| Oauth_token_secret | string | false    | ""            |
| Keywords           | array  | false    | []            |
| Follows            | array  | false    | []            |
| Full_tweet         | bool   | false    | false         |
| Ignore_retweets    | bool   | false    | false         |
| Languages          | array  | false    | []            |
| Locations          | array  | false    | []            |
| Tags               | array  | false    | []            |


## Details

### Add_field
* Value type is hash
* Default value is `{}`

If this filter is successful, add any arbitrary fields to this event.

### Consumer_key
* Value type is string
* Default value is `""`



### Consumer_secret
* Value type is string
* Default value is `""`



### Oauth_token
* Value type is string
* Default value is `""`



### Oauth_token_secret
* Value type is string
* Default value is `""`



### Keywords
* Value type is array
* Default value is `[]`

Any keywords to track in the Twitter stream. For multiple keywords,
use the syntax ["foo", "bar"]. There’s a logical OR between each keyword
string listed and a logical AND between words separated by spaces per keyword string.
See https://dev.twitter.com/streaming/overview/request-parameters#track for more details.

### Follows
* Value type is array
* Default value is `[]`

A comma separated list of user IDs, indicating the users to return statuses for
in the Twitter stream.
See https://dev.twitter.com/streaming/overview/request-parameters#follow for more details.

### Full_tweet
* Value type is bool
* Default value is `false`

Record full tweet object as given to us by the Twitter Streaming API

### Ignore_retweets
* Value type is bool
* Default value is `false`

Lets you ingore the retweets coming out of the Twitter API. Default false

### Languages
* Value type is array
* Default value is `[]`

A list of BCP 47 language identifiers corresponding to any of the languages
listed on Twitter’s advanced search page will only return tweets that have been
detected as being written in the specified languages

### Locations
* Value type is array
* Default value is `[]`

A comma-separated list of longitude, latitude pairs specifying a set of bounding boxes
to filter tweets by. See
https://dev.twitter.com/streaming/overview/request-parameters#locations for more details

### Tags
* Value type is array
* Default value is `[]`

If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
and include parts of the event using the %{field} syntax.



## Configuration blueprint

```
twitter{
	add_field => {}
	consumer_key => ""
	consumer_secret => ""
	oauth_token => ""
	oauth_token_secret => ""
	keywords => []
	follows => []
	full_tweet => bool
	ignore_retweets => bool
	languages => []
	locations => []
	tags => []
}
```
