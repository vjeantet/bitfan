+++
title = "Introduction"
description = ""
weight = 20
+++

Bitfan get a lot of inspiration from huggin and logstash.

Its configuration file format and is compatible and comes from logstash.

* You describe each of your usecase as a pipeline with a serie of processors organised as "inputs" > "filters" > "outputs".
* Pipeline specification format is a extension of the logstash's one. Some added features :
  * importing configuration file from a conf file.
  * routing event, on condition, to other configuration file.
  * using remote configuration file.
  * naming processors in your configuration (usefull to debug, API)

Bitfan runs **each** pipeline independently, its execution model allows to gracefully stop, start them without affecting other ones. See bitfan commands to operate a running bitfan.

Writtent with GoLang, Bitfan works on every platform and does not require a runtime or any library to run it.

{{%children style="h2" description="true"%}}