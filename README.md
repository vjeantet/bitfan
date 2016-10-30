# Logfan

Logstash implementation in GO.

![Logfan logo](docs/static/noun_307496_cc.png "Logfan")

## Install
2 ways to get logfan : download a released version or compile it from source.

### Download binary
linux, windows, osx available here : https://github.com/veino/logfan/releases

### Get source and compile
```
$ go get -u github.com/veino/logfan
```

## Usage
```
$ logfan -f https://raw.githubusercontent.com/veino/logfan/master/examples.d/simple.conf
```

copy/paste this in your console

```
127.0.0.1 - - [11/Dec/2013:00:01:45 -0800] "GET /xampp/status.php HTTP/1.1" 200 3891 "http://cadenza/xampp/navi.php" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:25.0) Gecko/20100101 Firefox/25.0"
```

## Use configuration from another location 
### include configuration from an URL
```
$ logfan -e "input{stdin{}} filter{use{url=>'https://raw.githubusercontent.com/veino/logfan/master/examples.d/use/lol/test.conf'}} output{stdout{codec=>rubydebug}}"
```

## use configuration file on local filesystem
```
$ logfan -e "input{stdin{}} filter{use{path=>'apachelogs.conf'}} output{stdout{codec=>rubydebug}}"
```

See examples in examples.d/use/ folder

### TODO

- [x] parse logstash config file
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
- [ ] log to file


# Supported inputs, filters and outputs in config file
can be found here : https://github.com/veino/processors

## input
* beats
* elasticsearch
* exec
* file
* httppoller
* imap
* rabbitmq
* stdin
* syslog
* twitter
* udp
* unix
* readfile

## filter
* date
* drop
* geoip
* grok
* json
* mutate
* split
* uuid
* html
* use - copy paste a configuration file here (URL or localpath)

## output
* elasticsearch v1, v2
* file
* glusterfs
* mongodb
* null
* rabbitmq
* statsd
* stdout

# Used package
* kardianos/govendor Go vendor tool that works with the standard vendor file
* spf13/cobra - A Commander for modern Go CLI interactions
* bbuck/go-lexer (a forked version) - Lexer based on Rob Pike's talk on YouTube
* veino/processors - all plugins used in logfan 

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
