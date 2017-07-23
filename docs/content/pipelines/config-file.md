+++
date = "2017-05-16T21:06:27+02:00"
description = ""
title = "Structure of a config File"
weight = 10
+++

A Bitfan config file has a separate section for each type of processor you want to add to the event processing pipeline. For example:

```
# This is a comment. You should use comments to describe
# parts of your configuration.
input {
  ...
}

filter {
  ...
}

output {
  ...
}
```

Each section contains the configuration options for one or more processors. If you specify multiple filters, they are applied in the order of their appearance in the configuration file.

## Processor configuration
The configuration of a processor consists of the processor name followed by a block of settings for that processor. For example, this input section configures two file inputs:

```js
input {
  file {
    path => "/var/log/messages"
    type => "syslog"
  }

  file {
    path => "/var/log/apache/access.log"
    type => "apache"
  }
}
```

In this example, two settings are configured for each of the file inputs: path and type.

The settings you can configure vary according to the processor type. For information about each processor, see [Input processors]({{%relref "inputs/_index.md"%}}), [Output processors]({{%relref "outputs/_index.md"%}}) and [Filter]({{%relref "filters/_index.md"%}}) processors.

two specials processors exists to extends your pipeline configuration : 

* [**use** processor]({{% relref "routers/use-processor.md" %}}) which allows you to include another configuration from another one
* [**route** processor]({{% relref "routers/route-processors.md" %}}) which allows you to route or fork events to other configuration files.

## Named processor
You can insert a name between the processor type and its configuration section, it can be usefull for debugging or when used with API.

example :
```js
input {
  file "messages-reader" {
    path => "/var/log/messages"
    type => "syslog"
  }

  file "apache-reader" {
    path => "/var/log/apache/access.log"
    type => "apache"
  }
}
```

## Value Types

A processor can require that the value for a setting be a
certain type, such as boolean, array, or hash. The following value
types are supported.

[{{%icon fa-arrow-circle-o-right%}} List of value types]({{%relref "config-value-types.md"%}})


## Comments

Comments are the same as in perl, ruby, and python. A comment starts with a '#' character, and does not need to be at the beginning of a line. For example:

```
# this is a comment

input { # comments can appear at the end of a line, too
  # ...
}
```
