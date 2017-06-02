+++
date = "2017-05-16T21:08:31+02:00"
description = ""
title = "Conditionals expressions"

[menu.main]
identifier = "pplconf-conditions"
parent = "pplconf"
weight = 40

+++

Sometimes you only want to filter or output an event under certain conditions. For that, you can use a conditional.

Conditionals in Bitfan look and act the same way they do in programming languages. Conditionals support if, else if and else statements and can be nested.

The conditional syntax is:
```
if EXPRESSION {
  ...
} else if EXPRESSION {
  ...
} else {
  ...
}
```

What’s an expression? Comparison tests, boolean logic, and so on!

You can use the following comparison operators:

* equality: ==, !=, <, >, <=, >=
* regexp: =~, !~ (checks a pattern on the right against a string value on the left)
* inclusion: in, not in

The supported boolean operators are:

* and, or

The supported unary operators are:

* !

Expressions can be long and complex. Expressions can contain other expressions, you can negate expressions with !, and you can group them with parentheses (...).

For example, the following conditional uses the mutate filter to remove the field secret if the field action has a value of login:

```js
filter {
  if [action] == "login" {
    mutate { remove_field => "secret" }
  }
}
```

You can specify multiple expressions in a single condition:

```js
output {
  # Display production errors to console
  if [loglevel] == "ERROR" and [deployment] == "production" {
    stdout {
    ...
    }
  }
}
```

You can use the in operator to test whether a field contains a specific string, key, or (for lists) element:

```js
filter {
  if [foo] in [foobar] {
    mutate { add_tag => "field in field" }
  }
  if [foo] in "foo" {
    mutate { add_tag => "field in string" }
  }
  if "hello" in [greeting] {
    mutate { add_tag => "string in field" }
  }
  if [foo] in ["hello", "world", "foo"] {
    mutate { add_tag => "field in list" }
  }
  if [missing] in [alsomissing] {
    mutate { add_tag => "shouldnotexist" }
  }
  if !("foo" in ["hello", "world"]) {
    mutate { add_tag => "shouldexist" }
  }
}
```

You use the not in conditional the same way. For example, you could use not in to only route events to Elasticsearch when grok is successful:

```js
output {
  if "_grokparsefailure" not in [tags] {
    elasticsearch { ... }
  }
}
```

You can check for the existence of a specific field, but there’s currently no way to differentiate between a field that doesn’t exist versus a field that’s simply false. The expression if [foo] returns false when:

* [foo] doesn’t exist in the event,
* [foo] exists in the event, but is false, or
* [foo] exists in the event, but is null



