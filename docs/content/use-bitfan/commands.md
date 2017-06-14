+++
title = "Command-Line Flags"
description = ""
name = "Commands"
weight = 300
+++

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

## Hidden flags
flags "à la logstash" works as well `-f`, `-e`, `--configtest`, ...



## Commands
{{%children description="true"%}}