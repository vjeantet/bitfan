package cmd

import (
	"github.com/vjeantet/bitfan/core"
	blacklist "github.com/vjeantet/bitfan/processors/filter-blacklist"
	change "github.com/vjeantet/bitfan/processors/filter-change"
	date "github.com/vjeantet/bitfan/processors/filter-date"
	digest "github.com/vjeantet/bitfan/processors/filter-digest"
	drop "github.com/vjeantet/bitfan/processors/filter-drop"
	evalprocessor "github.com/vjeantet/bitfan/processors/filter-eval"
	exec "github.com/vjeantet/bitfan/processors/filter-exec"
	geoip "github.com/vjeantet/bitfan/processors/filter-geoip"
	grok "github.com/vjeantet/bitfan/processors/filter-grok"
	html "github.com/vjeantet/bitfan/processors/filter-html"
	json "github.com/vjeantet/bitfan/processors/filter-json"
	kv "github.com/vjeantet/bitfan/processors/filter-kv"
	mutate "github.com/vjeantet/bitfan/processors/filter-mutate"
	newterm "github.com/vjeantet/bitfan/processors/filter-newterm"
	split "github.com/vjeantet/bitfan/processors/filter-split"
	uuid "github.com/vjeantet/bitfan/processors/filter-uuid"
	whitelist "github.com/vjeantet/bitfan/processors/filter-whitelist"
	httppoller "github.com/vjeantet/bitfan/processors/httppoller"
	beatsinput "github.com/vjeantet/bitfan/processors/input-beats"
	inputeventprocessor "github.com/vjeantet/bitfan/processors/input-event"
	execinput "github.com/vjeantet/bitfan/processors/input-exec"
	file "github.com/vjeantet/bitfan/processors/input-file"
	gennumbers "github.com/vjeantet/bitfan/processors/input-gennumbers"
	rabbitmqinput "github.com/vjeantet/bitfan/processors/input-rabbitmq"
	stdin "github.com/vjeantet/bitfan/processors/input-stdin"
	inputstdout "github.com/vjeantet/bitfan/processors/input-stdout"
	sysloginput "github.com/vjeantet/bitfan/processors/input-syslog"
	tail "github.com/vjeantet/bitfan/processors/input-tail"
	twitter "github.com/vjeantet/bitfan/processors/input-twitter"
	udpinput "github.com/vjeantet/bitfan/processors/input-udp"
	unixinput "github.com/vjeantet/bitfan/processors/input-unix"
	elasticsearch "github.com/vjeantet/bitfan/processors/output-elasticsearch"
	elasticsearch2 "github.com/vjeantet/bitfan/processors/output-elasticsearch2"
	email "github.com/vjeantet/bitfan/processors/output-email"
	fileoutput "github.com/vjeantet/bitfan/processors/output-file"
	glusterfsoutput "github.com/vjeantet/bitfan/processors/output-glusterfs"
	httpoutput "github.com/vjeantet/bitfan/processors/output-http"
	mongodb "github.com/vjeantet/bitfan/processors/output-mongodb"
	null "github.com/vjeantet/bitfan/processors/output-null"
	rabbitmqoutput "github.com/vjeantet/bitfan/processors/output-rabbitmq"
	statsd "github.com/vjeantet/bitfan/processors/output-statsd"
	route "github.com/vjeantet/bitfan/processors/route"
	stdout "github.com/vjeantet/bitfan/processors/stdout"

	use "github.com/vjeantet/bitfan/processors/use"
	when "github.com/vjeantet/bitfan/processors/when"

	httpoutprocessor "github.com/vjeantet/bitfan/processors/httpout"
	inputhttpserverprocessor "github.com/vjeantet/bitfan/processors/input-httpserver"
	ldapprocessor "github.com/vjeantet/bitfan/processors/ldap"
	sleepprocessor "github.com/vjeantet/bitfan/processors/sleep"
	sqlprocessor "github.com/vjeantet/bitfan/processors/sql"
	stopprocessor "github.com/vjeantet/bitfan/processors/stop"
	templateprocessor "github.com/vjeantet/bitfan/processors/template"
)

func init() {
	initPlugin("input", "stdout", inputstdout.New)
	initPlugin("input", "stdin", stdin.New)
	initPlugin("input", "twitter", twitter.New)
	initPlugin("input", "tail", tail.New) //
	initPlugin("input", "file", tail.New) // logstash's file plugin is a tail plugin
	initPlugin("input", "exec", execinput.New)
	initPlugin("input", "beats", beatsinput.New)
	initPlugin("input", "rabbitmq", rabbitmqinput.New)
	initPlugin("input", "udp", udpinput.New)
	initPlugin("input", "syslog", sysloginput.New)
	initPlugin("input", "unix", unixinput.New)
	initPlugin("input", "readfile", file.New)
	initPlugin("input", "sql", sqlprocessor.New)
	initPlugin("input", "http", httppoller.New)
	initPlugin("input", "use", use.New)
	initPlugin("input", "gennumbers", gennumbers.New)
	initPlugin("input", "ldap", ldapprocessor.New)
	initPlugin("input", "stop", stopprocessor.New)
	initPlugin("input", "httpserver", inputhttpserverprocessor.New)
	initPlugin("input", "event", inputeventprocessor.New)

	initPlugin("filter", "eval", evalprocessor.New)
	initPlugin("filter", "readfile", file.New)
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
	initPlugin("filter", "digest", digest.New)
	initPlugin("filter", "blacklist", blacklist.New)
	initPlugin("filter", "whitelist", whitelist.New)
	initPlugin("filter", "change", change.New)
	initPlugin("filter", "newterm", newterm.New)
	initPlugin("filter", "exec", exec.New)
	initPlugin("filter", "sql", sqlprocessor.New)
	initPlugin("filter", "template", templateprocessor.New)
	initPlugin("filter", "ldap", ldapprocessor.New)
	initPlugin("filter", "sleep", sleepprocessor.New)
	initPlugin("filter", "stdout", stdout.New)
	initPlugin("filter", "http", httppoller.New)
	initPlugin("filter", "httpout", httpoutprocessor.New)

	initPlugin("filter", "use", use.New)
	initPlugin("filter", "route", route.New)

	initPlugin("output", "stdout", stdout.New)
	initPlugin("output", "statsd", statsd.New)
	initPlugin("output", "mongodb", mongodb.New)
	initPlugin("output", "null", null.New)
	initPlugin("output", "elasticsearch", elasticsearch.New)
	initPlugin("output", "elasticsearch2", elasticsearch2.New)
	initPlugin("output", "file", fileoutput.New)
	initPlugin("output", "glusterfs", glusterfsoutput.New)
	initPlugin("output", "rabbitmq", rabbitmqoutput.New)
	initPlugin("output", "email", email.New)
	initPlugin("output", "http", httpoutput.New)
	initPlugin("output", "sql", sqlprocessor.New)
	initPlugin("output", "template", templateprocessor.New)
	initPlugin("output", "httpout", httpoutprocessor.New)

	initPlugin("output", "when", when.New)
	initPlugin("output", "use", use.New)
	// plugins = map[string]map[string]*core.ProcessorFactory{}

}

func initPluginsMap() map[string]map[string]core.ProcessorFactory {
	return map[string]map[string]core.ProcessorFactory{
		"input":  {},
		"filter": {},
		"output": {},
	}
}

var plugins = initPluginsMap()

func initPlugin(kind string, name string, proc core.ProcessorFactory) {
	pl := plugins[kind]
	pl[name] = proc
	plugins[kind] = pl

	prefix := kind + "_"
	if kind == "filter" {
		prefix = ""
	}
	core.RegisterProcessor(prefix+name, proc)
}
