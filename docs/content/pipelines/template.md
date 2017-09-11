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


## TS
Formats a time.Time with jodaTime layout

	{{TS "dd/MM/yyyy" . }}

will render an event's @timestamp like `2017-12-31 21:05:02 Local` as
	
	31/12/2017

## DateFormat
format a time.Time with [jodaTime layout](http://joda-time.sourceforge.net/apidocs/org/joda/time/format/DateTimeFormat.html)

## Time
converts the textual representation of the datetime string into a time.Time

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
