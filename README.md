<h1 align="center">
  <a href="https://bitfan.io"><img src="docs/static/open-fan-black-medium.png" height="150px" alt="Bitfan"></a>
  <br>
  Bitfan
</h1>

<h4 align="center">Data and Event processing pipeliner really inspired by Logstash.</h4>

<p align="center">
	<a href="https://godoc.org/github.com/vjeantet/bitfan">
		<img src="https://godoc.org/github.com/vjeantet/bitfan?status.svg" alt="GoDoc" style="max-width:100%;">
	</a>
    <a href="https://goreportcard.com/report/github.com/vjeantet/bitfan">
    	<img src="https://goreportcard.com/badge/github.com/vjeantet/bitfan" alt="Go Report Card" style="max-width:100%;">
    </a>
    <a href="https://travis-ci.org/vjeantet/bitfan">
    	<img src="https://travis-ci.org/vjeantet/bitfan.svg?branch=master" alt="Build Status" style="max-width:100%;">
    </a>
    <a href="https://codecov.io/gh/vjeantet/bitfan">
        <img src="https://codecov.io/gh/vjeantet/bitfan/branch/master/graph/badge.svg" alt="Codecov" />
    </a>
</p>

---

# Get bitfan, usage and configuration documentation and a availables processors 

 * Bitfan documentation [https://bitfan.io](https://bitfan.io)


# Features
- [x] configuration file compatible with logstash config file format.
- [x] support conditionals, env, sprintf variables in configuration  : %{[field][key]} ${ENVVAR}
- [x] supports input, filters, output and codecs
- [x] consume local and remote (http) configuration files
- [x] build complex pipelines with the `use` keyword to import, connect, fork to other pipelines/configuration files
- [x] gracefully stop and start each pipelines
- [x] install bitfan as a system daemon / service
- [x] manage running pipelines (list / stop / start a pipeline in a running bitfan)
- [x] monitor pipeline processors and events with prometheus
- [x] REST API to manage Bitfan
- [x] WebUI




# Similar projets in go

* tsaikd/gogstash - Logstash like, written in golang
* packetzoom/logzoom - A lightweight replacement for logstash indexer in Go
* hailocab/logslam - A lightweight lumberjack protocol compliant logstash indexer
* spartanlogs/spartan - Spartan is a data process much like Logstash

# Credits
Icon made by [Freepik](http://www.freepik.com) from [www.flaticon.com](https://www.flaticon.com/) is licensed by [CC 3.0 BY](http://creativecommons.org/licenses/by/3.0/)

# Contributors
* @vjeantet - Valere JEANTET
* @mirdhyn - Merlin Gaillard
* @AlexAkulov - Alexander AKULOV
* @lor00x - Thomas GUILLIER