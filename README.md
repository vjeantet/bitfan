# Bitfan

Bitfan is an open source data processing pipeline.
[![GoDoc](https://godoc.org/github.com/vjeantet/bitfan?status.svg)](https://godoc.org/github.com/vjeantet/bitfan)
[![Go Report Card](https://goreportcard.com/badge/github.com/vjeantet/bitfan)](https://goreportcard.com/report/github.com/vjeantet/bitfan)
[![Build Status](https://travis-ci.org/vjeantet/bitfan.svg?branch=master)](https://travis-ci.org/vjeantet/bitfan)

![Bitfan logo](docs/static/noun_307496_cc.png "Bitfan")

## Install

### Download binary
linux, windows, osx available here : https://github.com/vjeantet/bitfan/releases

### Or compile from sources
```
$ go get -u github.com/vjeantet/bitfan
```

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

  
```
Usage:
  bitfan [flags]
  bitfan [command]

Available Commands:
  doc         Display documentation about plugins
  list        List running pipelines
  run         Run bitfan
  service     Install and manage bitfan service
  start       Start a new pipeline in a running bitfan
  stop        Stop a running pipeline
  test        Test configurations (files, url, directories)
  version     Display version informations

Flags:
  -f, --config string       Load the Logstash config from a file a directory or a url
  -t, --configtest          Test config file or directory
      --debug               Increase verbosity to the last level (trace), more verbose.
  -e, --eval string         Use the given string as the configuration data.
  -w, --filterworkers int   number of workers (default 4)
  -h, --help                help for bitfan
  -l, --log string          Log to a given path. Default is to log to stdout.
      --settings string     Set the directory containing the bitfan.toml settings (default "current dir, then ~/.bitfan/ then /etc/bitfan/")
      --verbose             Increase verbosity to the first level (info), less verbose.
  -V, --version             Display version info.

Use "bitfan [command] --help" for more information about a command.
```

  logstash flags works as well `-f`, `-e`, `--configtest`, ...


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



# Similar projets in go

* tsaikd/gogstash - Logstash like, written in golang
* packetzoom/logzoom - A lightweight replacement for logstash indexer in Go
* hailocab/logslam - A lightweight lumberjack protocol compliant logstash indexer


# Credits
logo "hand fan" by lastspark from the Noun Project

# Contributors
* @vjeantet - Valere JEANTET
* @mirdhyn - Merlin Gaillard
* @AlexAkulov - Alexander AKULOV


# Used packages
* mitchellh/mapstructure
* ChimeraCoder/anaconda
* etrepat/postman/watch
* go-fsnotify/fsnotify
* hpcloud/tail
* nu7hatch/gouuid
* parnurzeal/gorequest
* vjeantet/go.enmime
* Knetic/govaluate
* vjeantet/grok
* vjeantet/jodaTime
* streadway/amqp
* ShowMax/go-fqdn
* oschwald/geoip2-golang
* gopkg.in/fsnotify.v1
* gopkg.in/go-playground/validator.v8
* gopkg.in/mgo.v2
* gopkg.in/olivere/elastic.v2
* gopkg.in/olivere/elastic.v3
* gopkg.in/alexcesaro/statsd.v2
* kardianos/govendor 
* spf13/cobra
* bbuck/go-lexer

* k0kubun/pp (debug)
