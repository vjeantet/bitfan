package commands

import (
	"bitfan/core"
	blacklist "bitfan/processors/filter-blacklist"
	change "bitfan/processors/filter-change"
	date "bitfan/processors/filter-date"
	digest "bitfan/processors/filter-digest"
	drop "bitfan/processors/filter-drop"
	evalprocessor "bitfan/processors/filter-eval"
	exec "bitfan/processors/filter-exec"
	geoip "bitfan/processors/filter-geoip"
	grok "bitfan/processors/filter-grok"
	html "bitfan/processors/filter-html"
	json "bitfan/processors/filter-json"
	kv "bitfan/processors/filter-kv"
	mutate "bitfan/processors/filter-mutate"
	newterm "bitfan/processors/filter-newterm"
	split "bitfan/processors/filter-split"
	uuid "bitfan/processors/filter-uuid"
	whitelist "bitfan/processors/filter-whitelist"
	httppoller "bitfan/processors/httppoller"
	beatsinput "bitfan/processors/input-beats"
	inputeventprocessor "bitfan/processors/input-event"
	execinput "bitfan/processors/input-exec"
	file "bitfan/processors/input-file"
	rabbitmqinput "bitfan/processors/input-rabbitmq"
	stdin "bitfan/processors/input-stdin"
	inputstdout "bitfan/processors/input-stdout"
	sysloginput "bitfan/processors/input-syslog"
	tail "bitfan/processors/input-tail"
	tcpinput "bitfan/processors/input-tcp"
	twitter "bitfan/processors/input-twitter"
	udpinput "bitfan/processors/input-udp"
	unixinput "bitfan/processors/input-unix"
	websocketinput "bitfan/processors/input-websocket"
	elasticsearch "bitfan/processors/output-elasticsearch"
	elasticsearch2 "bitfan/processors/output-elasticsearch2"
	email "bitfan/processors/output-email"
	fileoutput "bitfan/processors/output-file"
	glusterfsoutput "bitfan/processors/output-glusterfs"
	httpoutput "bitfan/processors/output-http"
	kafkaoutput "bitfan/processors/output-kafka"
	mongodb "bitfan/processors/output-mongodb"
	null "bitfan/processors/output-null"
	rabbitmqoutput "bitfan/processors/output-rabbitmq"
	statsd "bitfan/processors/output-statsd"
	tcpoutput "bitfan/processors/output-tcp"
	pop3processor "bitfan/processors/pop3"
	route "bitfan/processors/route"
	stdout "bitfan/processors/stdout"
	webfan "bitfan/processors/webfan"
	websocket "bitfan/processors/websocket"

	use "bitfan/processors/use"
	when "bitfan/processors/when"

	httpoutprocessor "bitfan/processors/httpout"
	inputhttpserverprocessor "bitfan/processors/input-httpserver"
	ldapprocessor "bitfan/processors/ldap"
	sleepprocessor "bitfan/processors/sleep"
	sqlprocessor "bitfan/processors/sql"
	stopprocessor "bitfan/processors/stop"
	templateprocessor "bitfan/processors/template"
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
