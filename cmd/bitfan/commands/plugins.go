package commands

import (
	"github.com/awillis/bitfan/core"
	blacklist "github.com/awillis/bitfan/processors/filter-blacklist"
	change "github.com/awillis/bitfan/processors/filter-change"
	date "github.com/awillis/bitfan/processors/filter-date"
	digest "github.com/awillis/bitfan/processors/filter-digest"
	drop "github.com/awillis/bitfan/processors/filter-drop"
	evalprocessor "github.com/awillis/bitfan/processors/filter-eval"
	exec "github.com/awillis/bitfan/processors/filter-exec"
	geoip "github.com/awillis/bitfan/processors/filter-geoip"
	grok "github.com/awillis/bitfan/processors/filter-grok"
	html "github.com/awillis/bitfan/processors/filter-html"
	json "github.com/awillis/bitfan/processors/filter-json"
	kv "github.com/awillis/bitfan/processors/filter-kv"
	mutate "github.com/awillis/bitfan/processors/filter-mutate"
	newterm "github.com/awillis/bitfan/processors/filter-newterm"
	split "github.com/awillis/bitfan/processors/filter-split"
	uuid "github.com/awillis/bitfan/processors/filter-uuid"
	whitelist "github.com/awillis/bitfan/processors/filter-whitelist"
	httppoller "github.com/awillis/bitfan/processors/httppoller"
	beatsinput "github.com/awillis/bitfan/processors/input-beats"
	inputeventprocessor "github.com/awillis/bitfan/processors/input-event"
	execinput "github.com/awillis/bitfan/processors/input-exec"
	file "github.com/awillis/bitfan/processors/input-file"
	"github.com/awillis/bitfan/processors/input-kafka"
	rabbitmqinput "github.com/awillis/bitfan/processors/input-rabbitmq"
	stdin "github.com/awillis/bitfan/processors/input-stdin"
	inputstdout "github.com/awillis/bitfan/processors/input-stdout"
	sysloginput "github.com/awillis/bitfan/processors/input-syslog"
	tail "github.com/awillis/bitfan/processors/input-tail"
	tcpinput "github.com/awillis/bitfan/processors/input-tcp"
	twitter "github.com/awillis/bitfan/processors/input-twitter"
	udpinput "github.com/awillis/bitfan/processors/input-udp"
	unixinput "github.com/awillis/bitfan/processors/input-unix"
	websocketinput "github.com/awillis/bitfan/processors/input-websocket"
	elasticsearch "github.com/awillis/bitfan/processors/output-elasticsearch"
	elasticsearch2 "github.com/awillis/bitfan/processors/output-elasticsearch2"
	email "github.com/awillis/bitfan/processors/output-email"
	fileoutput "github.com/awillis/bitfan/processors/output-file"
	glusterfsoutput "github.com/awillis/bitfan/processors/output-glusterfs"
	httpoutput "github.com/awillis/bitfan/processors/output-http"
	kafkaoutput "github.com/awillis/bitfan/processors/output-kafka"
	mongodb "github.com/awillis/bitfan/processors/output-mongodb"
	null "github.com/awillis/bitfan/processors/output-null"
	rabbitmqoutput "github.com/awillis/bitfan/processors/output-rabbitmq"
	statsd "github.com/awillis/bitfan/processors/output-statsd"
	tcpoutput "github.com/awillis/bitfan/processors/output-tcp"
	pop3processor "github.com/awillis/bitfan/processors/pop3"
	route "github.com/awillis/bitfan/processors/route"
	stdout "github.com/awillis/bitfan/processors/stdout"
	webfan "github.com/awillis/bitfan/processors/webfan"
	websocket "github.com/awillis/bitfan/processors/websocket"

	use "github.com/awillis/bitfan/processors/use"
	when "github.com/awillis/bitfan/processors/when"

	httpoutprocessor "github.com/awillis/bitfan/processors/httpout"
	inputhttpserverprocessor "github.com/awillis/bitfan/processors/input-httpserver"
	ldapprocessor "github.com/awillis/bitfan/processors/ldap"
	sleepprocessor "github.com/awillis/bitfan/processors/sleep"
	sqlprocessor "github.com/awillis/bitfan/processors/sql"
	stopprocessor "github.com/awillis/bitfan/processors/stop"
	templateprocessor "github.com/awillis/bitfan/processors/template"
)

func init() {
	initPlugin("input", "webhook", webfan.New)
	initPlugin("input", "stdout", inputstdout.New)
	initPlugin("input", "stdin", stdin.New)
	initPlugin("input", "twitter", twitter.New)
	initPlugin("input", "tail", tail.New) //
	initPlugin("input", "file", tail.New) // logstash's file plugin is a tail plugin
	initPlugin("input", "exec", execinput.New)
	initPlugin("input", "beats", beatsinput.New)
	initPlugin("input", "rabbitmq", rabbitmqinput.New)
	initPlugin("input", "udp", udpinput.New)
	initPlugin("input", "tcp", tcpinput.New)
	initPlugin("input", "syslog", sysloginput.New)
	initPlugin("input", "unix", unixinput.New)
	initPlugin("input", "readfile", file.New)
	initPlugin("input", "sql", sqlprocessor.New)
	initPlugin("input", "http", httppoller.New)
	initPlugin("input", "use", use.New)
	initPlugin("input", "ldap", ldapprocessor.New)
	initPlugin("input", "stop", stopprocessor.New)
	initPlugin("input", "httpserver", inputhttpserverprocessor.New)
	initPlugin("input", "event", inputeventprocessor.New)
	initPlugin("input", "websocket", websocketinput.New)
	initPlugin("input", "pop3", pop3processor.New)
	initPlugin("input", "kafka", kafkainput.New)

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
	initPlugin("filter", "pop3", pop3processor.New)

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
	initPlugin("output", "kafka", kafkaoutput.New)
	initPlugin("output", "tcp", tcpoutput.New)
	initPlugin("output", "sql", sqlprocessor.New)
	initPlugin("output", "template", templateprocessor.New)
	initPlugin("output", "httpout", httpoutprocessor.New)
	initPlugin("output", "websocket", websocket.New)

	initPlugin("output", "when", when.New)
	initPlugin("output", "use", use.New)
	initPlugin("output", "pass", webfan.NewPass)
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
