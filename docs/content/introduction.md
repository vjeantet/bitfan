+++
title = "Introduction"
description = ""

[menu.main]
parent = ""
name = "Introduction"
identifier = "introduction"
weight = 2
+++

Bitfan is an event and data processing system, which allows you to perform automated tasks from/with a multitude of sources.

> Ingest or query data, detect event, transform and enrich them to finally take any actions on them.

Bitfan get a lot of inspiration from logstash and huggin.

You describe each of your usecase as a pipeline with a serie of processors organised as "inputs", "filters" and "outputs".
Pipeline specification format is a extension of the logstash.


Bitfan runs your pipelines independently, its execution model allows to gracefully stop, start them without affecting other ones.

It works as-is on every plateform, and does not require any runtime. 

Just download, write your first pipeline spec and execute it !

See pipeline library if you want to save the "write pipeline" part :-)

## Use cases

### Email Report
Every Monday, I receive a mail with the KPI of my team.
Bitfan execute several queries to multiples sql databases, results are sent to a digest processor which waits for 10:00 AM to compute received data with a HTML template.

### IOC alert
When a device on my network requests a domain know to be used by wannacry, alert !

Send a http request to a specific endpoint and send a mail to security teams when an event with a blacklisted word appears in the dns server log.

### ETL like
Search for entries in a LDAP, and insert each of them in a mysql database



Its configuration format is a extension of the logstash's one, with ability to reference other configuration file and load configuration from a remote repository



Bitfan works on every platform and does not require a runtime or any library to run it.




{{%children style="h2" description="true"%}}