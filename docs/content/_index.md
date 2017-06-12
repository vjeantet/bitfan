+++
title = "Home"
description = ""
+++

<span id="sidebar-toggle-span">
<a href="#" id="sidebar-toggle" data-sidebar-toggle=""><i class="fa fa-bars"></i></a>
</span>


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
- [x] codec support
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

* [More on processors]({{< relref "processors.md" >}})
