# Bitfan

Bitfan is an open source data processing pipeline really inspired by Logstash.

Bitfan is written in Go and should build on all platforms.

![Bitfan logo](docs/static/noun_307496_cc.png "Bitfan")

## Get bitfan, usage and configuration documentation and a availables processors 

 * Bitfan documentation [https://bitfan.io](https://bitfan.io)

[![GoDoc](https://godoc.org/github.com/vjeantet/bitfan?status.svg)](https://godoc.org/github.com/vjeantet/bitfan)
[![Go Report Card](https://goreportcard.com/badge/github.com/vjeantet/bitfan)](https://goreportcard.com/report/github.com/vjeantet/bitfan)
[![Build Status](https://travis-ci.org/vjeantet/bitfan.svg?branch=master)](https://travis-ci.org/vjeantet/bitfan)


## Features
- [x] configuration file compatible with logstash config file format.
- [x] support conditionals, env, sprintf variables in configuration  : %{[field][key]} ${ENVVAR}
- [x] supports input, filters, output and codecs
- [x] consume local and remote (http) configuration files
- [x] build complex pipelines with the `use` keyword to import, connect, fork to other pipelines/configuration files
- [x] gracefully stop and start each pipelines
- [x] install bitfan as a system daemon / service
- [x] manage running pipelines (list / stop / start a pipeline in a running bitfan)
- [x] monitor pipeline processors and events with prometheus
- [ ] REST API to manage Bitfan (WIP)

# Similar projets in go

* tsaikd/gogstash - Logstash like, written in golang
* packetzoom/logzoom - A lightweight replacement for logstash indexer in Go
* hailocab/logslam - A lightweight lumberjack protocol compliant logstash indexer
* spartanlogs/spartan - Spartan is a data process much like Logstash

# Credits
logo "hand fan" by lastspark from the Noun Project

# Contributors
* @vjeantet - Valere JEANTET
* @mirdhyn - Merlin Gaillard
* @AlexAkulov - Alexander AKULOV
