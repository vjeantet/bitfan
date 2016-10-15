package cmd

import (
	date "github.com/veino/veino/processors/filter-date"
	drop "github.com/veino/veino/processors/filter-drop"
	geoip "github.com/veino/veino/processors/filter-geoip"
	grok "github.com/veino/veino/processors/filter-grok"
	html "github.com/veino/veino/processors/filter-html"
	json "github.com/veino/veino/processors/filter-json"
	kv "github.com/veino/veino/processors/filter-kv"
	mutate "github.com/veino/veino/processors/filter-mutate"
	split "github.com/veino/veino/processors/filter-split"
	uuid "github.com/veino/veino/processors/filter-uuid"
	beatsinput "github.com/veino/veino/processors/input-beats"
	execinput "github.com/veino/veino/processors/input-exec"
	fileinput "github.com/veino/veino/processors/input-file"
	rabbitmqinput "github.com/veino/veino/processors/input-rabbitmq"
	stdin "github.com/veino/veino/processors/input-stdin"
	sysloginput "github.com/veino/veino/processors/input-syslog"
	twitter "github.com/veino/veino/processors/input-twitter"
	udpinput "github.com/veino/veino/processors/input-udp"
	unixinput "github.com/veino/veino/processors/input-unix"
	elasticsearch "github.com/veino/veino/processors/output-elasticsearch"
	elasticsearch2 "github.com/veino/veino/processors/output-elasticsearch2"
	fileoutput "github.com/veino/veino/processors/output-file"
	glusterfsoutput "github.com/veino/veino/processors/output-glusterfs"
	mongodb "github.com/veino/veino/processors/output-mongodb"
	null "github.com/veino/veino/processors/output-null"
	rabbitmqoutput "github.com/veino/veino/processors/output-rabbitmq"
	statsd "github.com/veino/veino/processors/output-statsd"
	stdout "github.com/veino/veino/processors/output-stdout"
	when "github.com/veino/veino/processors/when"
	"github.com/veino/veino/runtime"
)

func init() {

	runtime.RegisterProcessor("input_stdin", stdin.New)
	runtime.RegisterProcessor("input_twitter", twitter.New)
	runtime.RegisterProcessor("input_file", fileinput.New)
	runtime.RegisterProcessor("input_exec", execinput.New)
	runtime.RegisterProcessor("input_beats", beatsinput.New)
	runtime.RegisterProcessor("input_rabbitmq", rabbitmqinput.New)
	runtime.RegisterProcessor("input_udp", udpinput.New)
	runtime.RegisterProcessor("input_syslog", sysloginput.New)
	runtime.RegisterProcessor("input_unix", unixinput.New)

	runtime.RegisterProcessor("grok", grok.New)
	runtime.RegisterProcessor("mutate", mutate.New)
	runtime.RegisterProcessor("split", split.New)
	runtime.RegisterProcessor("date", date.New)
	runtime.RegisterProcessor("json", json.New)
	runtime.RegisterProcessor("uuid", uuid.New)
	runtime.RegisterProcessor("drop", drop.New)
	runtime.RegisterProcessor("geoip", geoip.New)
	runtime.RegisterProcessor("kv", kv.New)
	runtime.RegisterProcessor("html", html.New)

	runtime.RegisterProcessor("output_stdout", stdout.New)
	runtime.RegisterProcessor("output_statsd", statsd.New)
	runtime.RegisterProcessor("output_mongodb", mongodb.New)
	runtime.RegisterProcessor("output_null", null.New)
	runtime.RegisterProcessor("output_elasticsearch", elasticsearch.New)
	runtime.RegisterProcessor("output_elasticsearch2", elasticsearch2.New)
	runtime.RegisterProcessor("output_file", fileoutput.New)
	runtime.RegisterProcessor("output_glusterfs", glusterfsoutput.New)
	runtime.RegisterProcessor("output_rabbitmq", rabbitmqoutput.New)

	runtime.RegisterProcessor("when", when.New)
	runtime.RegisterProcessor("output_when", when.New)
}
