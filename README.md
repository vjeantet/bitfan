# Logfan

Logstash like implementation with GO.

![Logfan logo](docs/static/noun_307496_cc.png "Logfan")

## Install
2 ways to get logfan : download a released version or compile it from source.

### Download binary
linux, windows, osx available here : https://github.com/veino/logfan/releases

### Get source and compile
```
$ go get -u github.com/veino/logfan
```

## Run 
```
$ logfan run https://raw.githubusercontent.com/veino/logfan/master/examples.d/simple.conf
```

	logstash flags works as well `-f`, `-e`, `--configtest`, ...

copy/paste this in your console

```
127.0.0.1 - - [11/Dec/2013:00:01:45 -0800] "GET /xampp/status.php HTTP/1.1" 200 3891 "http://cadenza/xampp/navi.php" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:25.0) Gecko/20100101 Firefox/25.0"
```

## Other command
type `logfan help` to display usage information

	
```
Usage:
  logfan [flags]
  logfan [command]

Available Commands:
  doc         Display documentation about plugins
  list        List running pipelines
  run         Run logfan
  service     Install and manage logfan service
  start       Start a new pipeline in a running logfan
  stop        Stop a running pipeline
  test        Test configurations (files, url, directories)
  version     Display version informations

Flags:
  -f, --config string       Load the Logstash config from a file a directory or a url
  -t, --configtest          Test config file or directory
      --debug               Increase verbosity to the last level (trace), more verbose.
  -e, --eval string         Use the given string as the configuration data.
  -w, --filterworkers int   number of workers (default 4)
  -h, --help                help for logfan
  -l, --log string          Log to a given path. Default is to log to stdout.
      --settings string     Set the directory containing the logfan.toml settings (default "current dir, then ~/.logfan/ then /etc/logfan/")
      --verbose             Increase verbosity to the first level (info), less verbose.
  -V, --version             Display version info.

Use "logfan [command] --help" for more information about a command.
```

## include configuration from other configuration
### include configuration from an URL
```
$ logfan run "input{stdin{}} filter{use{url=>'https://raw.githubusercontent.com/veino/logfan/master/examples.d/use/lol/test.conf'}} output{stdout{codec=>rubydebug}}"
```

## use configuration file on local filesystem
```
$ logfan run "input{stdin{}} filter{use{path=>'apachelogs.conf'}} output{stdout{codec=>rubydebug}}"
```

See examples in examples.d/use/ folder

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
- [x] install logfan as a system daemon / service
- [x] list currently runnnung pipelines
- [x] start new pipelines in a running instance
- [x] stop a pipeline without stopping other



# Supported inputs, filters and outputs in config file
can be found here : https://github.com/veino/veino/tree/master/processors

type `logfan doc` to list all available plugins

## INPUT

|  PLUGIN  |          DESCRIPTION           |
|----------|--------------------------------|
| twitter  |                                |
| exec     |                                |
| unix     |                                |
| stdin    | Reads events from standard  input  |
| file     |                                |
| beats    |                                |
| rabbitmq |                                |
| udp      |                                |
| syslog   |                                |
| readfile |                                |

type `logfan doc pluginname` to get more information about plugin configuration and usage

## FILTER

| PLUGIN |          DESCRIPTION           |
|--------|--------------------------------|
| date   | Parses dates from fields to use as the Logfan timestamp  for an event |
| grok   |                                |
| split  | Splits multi-line messages into distinct events |
| json   | Parses JSON events             |
| uuid   | Adds a UUID to events          |
| drop   | Drops all events               |
| geoip  | Adds geographical information about an IP address |
| kv     | Parses key-value pairs         |
| html   |                                |
| mutate |                                |

type `logfan doc pluginname` to get more information about plugin configuration and usage

## OUTPUT

|     PLUGIN     |          DESCRIPTION           |
|----------------|--------------------------------|
| stdout         | Prints events to the standard output |
| null           | Drops everything received      |
| file           |                                |
| glusterfs      |                                |
| statsd         |                                |
| mongodb        |                                |
| elasticsearch  |                                |
| elasticsearch2 |                                |
| rabbitmq       |                                |

type `logfan doc pluginname` to get more information about plugin configuration and usage

## SPECIAL for all sections
|     PLUGIN     |          DESCRIPTION           |
|----------------|--------------------------------|
| use         | reference another configuration file (URL or local path) to include (copy/paste) in your current configuration  |


# Used package
* kardianos/govendor Go vendor tool that works with the standard vendor file
* spf13/cobra - A Commander for modern Go CLI interactions
* bbuck/go-lexer (a forked version) - Lexer based on Rob Pike's talk on YouTube
* veino/veino - all plugins and runtime used by logfan 


# Similar projets in go

* tsaikd/gogstash - Logstash like, written in golang
* packetzoom/logzoom - A lightweight replacement for logstash indexer in Go
* hailocab/logslam - A lightweight lumberjack protocol compliant logstash indexer


# Credits
logo "hand fan" by lastspark from the Noun Project

# Contributors
* Valere JEANTET
* Merlin Gaillard
* Alexander AKULOV
