package cmd

import (
	"github.com/veino/processors/filter-date"
	"github.com/veino/processors/filter-drop"
	"github.com/veino/processors/filter-geoip"
	"github.com/veino/processors/filter-grok"
	"github.com/veino/processors/filter-json"
	"github.com/veino/processors/filter-mutate"
	"github.com/veino/processors/filter-split"
	"github.com/veino/processors/filter-uuid"
	"github.com/veino/processors/input-amqp"
	"github.com/veino/processors/input-beats"
	"github.com/veino/processors/input-exec"
	"github.com/veino/processors/input-file"
	"github.com/veino/processors/input-stdin"
	"github.com/veino/processors/input-twitter"
	"github.com/veino/processors/output-elasticsearch"
	"github.com/veino/processors/output-elasticsearch2"
	"github.com/veino/processors/output-file"
	"github.com/veino/processors/output-glusterfs"
	"github.com/veino/processors/output-mongodb"
	"github.com/veino/processors/output-null"
	"github.com/veino/processors/output-rabbitmq"
	"github.com/veino/processors/output-stdout"
	"github.com/veino/processors/when"
	"github.com/veino/runtime"
)

func init() {

	runtime.RegisterProcessor("input_stdin", stdin.New)
	runtime.RegisterProcessor("input_twitter", twitter.New)
	runtime.RegisterProcessor("input_file", fileinput.New)
	runtime.RegisterProcessor("input_exec", execinput.New)
	runtime.RegisterProcessor("input_beats", beatsinput.New)
	runtime.RegisterProcessor("input_rabbitmq", rabbitmqinput.New)

	runtime.RegisterProcessor("grok", grok.New)
	runtime.RegisterProcessor("mutate", mutate.New)
	runtime.RegisterProcessor("split", split.New)
	runtime.RegisterProcessor("date", date.New)
	runtime.RegisterProcessor("json", json.New)
	runtime.RegisterProcessor("uuid", uuid.New)
	runtime.RegisterProcessor("drop", drop.New)
	runtime.RegisterProcessor("geoip", geoip.New)

	runtime.RegisterProcessor("output_stdout", stdout.New)
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
