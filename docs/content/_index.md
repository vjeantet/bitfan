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



# Supported inputs, filters and outputs in config file

type `bitfan doc` to list all available plugins

* [More on processors]({{< relref "processors/_index.md" >}})
