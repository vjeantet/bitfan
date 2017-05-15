+++
title = "Home"
description = ""
+++


> Bitfan is an open source data processing pipeline.


## Run 
Example with a remote configuration file which ingest data from stdin and output a tranformation to stdout.
```
$ bitfan run https://raw.githubusercontent.com/vjeantet/bitfan/master/examples.d/simple.conf
```
copy/paste this in your console

```
127.0.0.1 - - [11/Dec/2013:00:01:45 -0800] "GET /xampp/status.php HTTP/1.1" 200 3891 "http://cadenza/xampp/navi.php" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:25.0) Gecko/20100101 Firefox/25.0"
```

## Other commands
type `bitfan help` to display usage information

  
## TODO

- [x] parse logstash config file
- [x] support command line flags "à la logstash"
- [x] generic input support
- [x] generic filter support
- [x] generic output support
- [x] configuration condition (if else) support
- [x] dynamic %{field.key} support in config file
- [x] gracefully stop
- [x] gracefully start
- [x] name all contributors, imported packages, similar projects
- [x] use remote configuration file
- [x] include local and remote files from configuration files
- [ ] codec support
- [x] log to file
- [x] plugins autodocumentation
- [x] install bitfan as a system daemon / service
- [x] list currently runnnung pipelines
- [x] start new pipelines in a running instance
- [x] stop a pipeline without stopping other
- [x] import external configuration from configuration (use)
- [x] dispatch message to another configuration from configuration (fork)


# Supported inputs, filters and outputs in config file

type `bitfan doc` to list all available plugins

## INPUT `bitfan doc --type input`

|   PLUGIN   |          DESCRIPTION           |
|------------|--------------------------------|
| rabbitmq   |                                |
| syslog     |                                |
| sql        |                                |
| file       |                                |
| beats      |                                |
| exec       |                                |
| unix       |                                |
| http       |                                |
| gennumbers | generate a number of event     |
| stdin      | Reads events from stdin        |
| tail       |                                |
| readfile   |                                |
| twitter    |                                |
| udp        |                                |
| use        | Include a config file          |

type `bitfan doc pluginname` to get more information about plugin configuration and usage

## FILTER `bitfan doc --type filter`

| PLUGIN |          DESCRIPTION           |
|--------|--------------------------------|
| grok   |                                |
| split  | Splits multi-line messages into distinct events     |
| drop   | Drops all events               |
| html   |                                |
| digest | Digest events every x          |
| mutate |                                |
| json   | Parses JSON events             |
| uuid   | Adds a UUID to events          |
| date   | Parses dates from fields    |
| geoip  | Adds geographical information from IP |
| kv     | Parses key-value pairs         |
| use    | Include a config file          |
| route  | route message to other pipelines  |

type `bitfan doc pluginname` to get more information about plugin configuration and usage

## OUTPUT `bitfan doc --type output`

|     PLUGIN     |          DESCRIPTION           |
|----------------|--------------------------------|
| null           | Drops everything received      |
| elasticsearch  |                                |
| elasticsearch2 |                                |
| file           |                                |
| stdout         | Prints events to the stdout    |
| statsd         |                                |
| mongodb        |                                |
| glusterfs      |                                |
| rabbitmq       |                                |
| email          | Sends email      |
| use            | Include a config file          |

type `bitfan doc pluginname` to get more information about plugin configuration and usage