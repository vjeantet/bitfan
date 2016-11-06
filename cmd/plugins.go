package cmd

import (
	"github.com/veino/veino"
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
	readfile "github.com/veino/veino/processors/input-readfile"
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
	initPlugin("input", "stdin", stdin.New)
	initPlugin("input", "twitter", twitter.New)
	initPlugin("input", "file", fileinput.New)
	initPlugin("input", "exec", execinput.New)
	initPlugin("input", "beats", beatsinput.New)
	initPlugin("input", "rabbitmq", rabbitmqinput.New)
	initPlugin("input", "udp", udpinput.New)
	initPlugin("input", "syslog", sysloginput.New)
	initPlugin("input", "unix", unixinput.New)
	initPlugin("input", "readfile", readfile.New)

	initPlugin("filter", "grok", grok.New)
	initPlugin("filter", "mutate", mutate.New)
	initPlugin("filter", "split", split.New)
	initPlugin("filter", "date", date.New)
	initPlugin("filter", "json", json.New)
	initPlugin("filter", "uuid", uuid.New)
	initPlugin("filter", "drop", drop.New)
	initPlugin("filter", "geoip", geoip.New)
	initPlugin("filter", "kv", kv.New)
	initPlugin("filter", "html", html.New)
	initPlugin("filter", "when", when.New)

	initPlugin("output", "stdout", stdout.New)
	initPlugin("output", "statsd", statsd.New)
	initPlugin("output", "mongodb", mongodb.New)
	initPlugin("output", "null", null.New)
	initPlugin("output", "elasticsearch", elasticsearch.New)
	initPlugin("output", "elasticsearch2", elasticsearch2.New)
	initPlugin("output", "file", fileoutput.New)
	initPlugin("output", "glusterfs", glusterfsoutput.New)
	initPlugin("output", "rabbitmq", rabbitmqoutput.New)

	initPlugin("output", "when", when.New)
	// plugins = map[string]map[string]*veino.ProcessorFactory{}

}

func initPluginsMap() map[string]map[string]veino.ProcessorFactory {
	return map[string]map[string]veino.ProcessorFactory{
		"input":  map[string]veino.ProcessorFactory{},
		"filter": map[string]veino.ProcessorFactory{},
		"output": map[string]veino.ProcessorFactory{},
	}
}

var plugins = initPluginsMap()

func initPlugin(kind string, name string, proc veino.ProcessorFactory) {
	pl := plugins[kind]
	pl[name] = proc
	plugins[kind] = pl

	prefix := kind + "_"
	if kind == "filter" {
		prefix = ""
	}
	runtime.RegisterProcessor(prefix+name, proc)
}
