+++
title = "Command-Line Flags"
description = ""
name = "Commands"
weight = 300
+++

# bitfan
```
Usage:
  bitfan [flags]
  bitfan [command]

Available Commands:
  conf        Retrieve configuration file and its related files of a running pipeline
  doc         Display documentation about plugins
  list        List running pipelines
  run         Run bitfan
  service     Install and manage bitfan service
  start       Start a pipeline to the running bitfan
  stop        Stop a running pipeline
  test        Test configurations (files, url, directories)
  version     Display version informations

Flags:
      --debug             Increase verbosity to the last level (trace)
  -l, --log string        Log to a given path. Default is to log to stdout.
      --settings string   Set the directory containing the bitfan.toml settings (default "current dir, then ~/.bitfan/ then /etc/bitfan/")
      --verbose           Increase verbosity of logs (default true)

Use "bitfan [command] --help" for more information about a command.
```

## Hidden flags
flags "à la logstash" works as well `-f`, `-e`, `--configtest`, ...

## Commands
{{%children description="true"%}}


# bitfanUI
```
Usage:
  bitfan-ui [flags]
  bitfan-ui [command]

Available Commands:
  service     Install and manage bitfan service
  version     Display version informations

Flags:
  -a, --api string      Bitfan API to connect to (default "127.0.0.1:5123")
      --config string   config file (default : /etc/bitfan/bitfan-ui.toml, $HOME/.bitfan/bitfan-ui.toml, ./bitfan-ui.toml)
      --dev             dev mode (serve asset and templates from disk
  -H, --host string     Serve UI on Host (default "127.0.0.1:8081")

Use "bitfan-ui [command] --help" for more information about a command.

```

## Commands
{{%children description="true"%}}
