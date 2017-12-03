+++
date = "2017-05-16T21:08:14+02:00"
description = ""
title = "Using templates"
name = "Using templates"
weight = 20
+++

* <a href="#ts">TS</a> - format event's @timestamp with [jodaTime layout](http://joda-time.sourceforge.net/apidocs/org/joda/time/format/DateTimeFormat.html)
* <a href="#dateformat">DateFormat</a> - format a time.Time with [jodaTime layout](http://joda-time.sourceforge.net/apidocs/org/joda/time/format/DateTimeFormat.html)
* <a href="#time">Time</a> - converts the textual representation of the datetime string into a time.Time
* <a href="#now">Now</a> - returns the current local time.
* <a href="#numfmt">NumFmt</a> - format a number
* <a href="#safehtml">SafeHTML</a> - use string as HTML
* <a href="#htmlunescape">HTMLUnescape</a> - unescape a html string
* <a href="#htmlescape">HTMLEscape</a> - escape a html string
* <a href="#lower">Lower</a> - lowercase a string
* <a href="#upper">Upper</a> - uppercase a string
* <a href="#trim">Trim</a> - trim a string
* <a href="#trimprefix">TrimPrefix</a> - trim a prefix from a string
* <a href="#ago">Ago</a> - format durations according to a format string
* <a href="#markdown">Markdown</a> - renders a markdown string value to html
* <a href="#int">Int</a> - converts a numeric string value to int
* <a href="#_">_</a> - returns a value by key for a given map, empty string on failure


## TS
Formats a time.Time with jodaTime layout

	{{TS "dd/MM/yyyy" . }}

will render an event's @timestamp like `2017-12-31 21:05:02 Local` as
	
	31/12/2017

## DateFormat
format a time.Time with [jodaTime layout](http://joda-time.sourceforge.net/apidocs/org/joda/time/format/DateTimeFormat.html)

	{{DateFormat "dd/MM/YYYY" .datetimefield)}}

will render an event's datetimefield field value like `2017-12-31 21:05:02 Local` as
	
	31/12/2017


## Time
converts the textual representation of the datetime string into a time.Time

	(Time .txtdatefield)

## Now
returns the current local time.

## NumFmt
format a number

## SafeHTML
use string as HTML

## HTMLUnescape
unescape a html string

## HTMLEscape
escape a html string

## Lower
lowercase a string

## Upper
uppercase a string

## Trim
trim a string

## TrimPrefix
trim a prefix from a string

	{{TrimPrefix .stringvaluefield "name "}}

will render an event's stringvaluefield field value like "name Valere" as

	Valere

## Int
converts a numeric string value to int

## Markdown
renders a markdown string value to html

	{{markdown .stringfield}}

## _

	city : {{_ "city" .location}}

renders 

	city : Paris


## Ago

	updated {{ago "%d days ago" (Time .CREATEDDATESTRING)}}

renders

	updated 34 days ago


### Duration Format

The % character signifies that the next character is a modifier that specifies a particular duration unit. The following is the full list of modifiers supported:

* `%y` - # of years
* `%w` - # of weeks
* `%d` - # of days
* `%h` - # of hours
* `%m` - # of minutes
* `%s` - # of seconds
* `%%` - print a percent sign