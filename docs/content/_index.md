+++
title = "Home"
description = ""
+++

# BitFan

Bitfan is an open source data processing software system compatible with logstash configuration format.

> Ingest or query data, detect event, transform and enrich them to finally take any actions on them.


I use it as a “Swiss Army Knife” to complete a wide variety of different tasks, such as:


* Loading and **parsing** log files from a file system.
* Performing real time **anomaly detection** on any data flowing through a pipeline.
* **Shipping data** from one location to another with transformation.
* Sending weekly **email reports** computed from multiples sources datastores.
* Launching external processes to gather operational data from the local system.


![scope](../images/screenshot.png)


## Processors / Plugins availables

<table>
	<thead>
		<tr>
			<th><a href="/inputs/">INPUTS</a></th>
			<th><a href="/filters/">FILTERS</a></th>
			<th><a href="/outputs/">OUTPUTS</a></th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td style="vertical-align: top">{{%children page="inputs"%}}</td>
			<td style="vertical-align: top">{{%children page="filters"%}}</td>
			<td style="vertical-align: top">{{%children page="outputs"%}}</td>
		</tr>
	</tbody>
</table>

{{%info%}}Type `$ bitfan doc` to list all available plugins and get usage doc about them.{{%/info%}}

## Very QuickStart
start a pipeline from a configuration file hosted on github.com.

> this pipeline configuration ingests data from stdin and output a tranformation to stdout. have a look here : [see configuration file](https://raw.githubusercontent.com/vjeantet/bitfan/master/cmd/bitfan/examples.d/simple.conf)

```
$ bitfan run https://raw.githubusercontent.com/vjeantet/bitfan/master/cmd/bitfan/examples.d/simple.conf
```

Feed the pipeline with a copy/paste of following lines in your console :

```
127.0.0.1 - - [11/Dec/2013:00:01:45 -0800] "GET /xampp/status.php HTTP/1.1" 200 3891 "http://cadenza/xampp/navi.php" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:25.0) Gecko/20100101 Firefox/25.0"
```
{{%info%}}Type `bitfan help` to bitfan display usage information.{{%/info%}}


## More on bitfan


{{%children style="h2" description="true" page="home"%}}

