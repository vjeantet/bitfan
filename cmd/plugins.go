package cmd

import (
	"github.com/veino/processor-filter-date"
	"github.com/veino/processor-filter-drop"
	"github.com/veino/processor-filter-grok"
	"github.com/veino/processor-filter-json"
	"github.com/veino/processor-filter-mutate"
	"github.com/veino/processor-filter-split"
	"github.com/veino/processor-filter-uuid"
	"github.com/veino/processor-input-exec"
	"github.com/veino/processor-input-file"
	"github.com/veino/processor-input-stdin"
	"github.com/veino/processor-input-twitter"
	"github.com/veino/processor-output-elasticsearch"
	"github.com/veino/processor-output-mongodb"
	"github.com/veino/processor-output-null"
	"github.com/veino/processor-output-stdout"
	"github.com/veino/processor-when"
	"github.com/veino/runtime"
)

func init() {

	runtime.RegisterProcessor("input_stdin", stdin.New)
	runtime.RegisterProcessor("input_twitter", twitter.New)
	runtime.RegisterProcessor("input_file", fileinput.New)
	runtime.RegisterProcessor("input_exec", execinput.New)

	runtime.RegisterProcessor("grok", grok.New)
	runtime.RegisterProcessor("mutate", mutate.New)
	runtime.RegisterProcessor("split", split.New)
	runtime.RegisterProcessor("date", date.New)
	runtime.RegisterProcessor("json", json.New)
	runtime.RegisterProcessor("uuid", uuid.New)
	runtime.RegisterProcessor("drop", drop.New)

	runtime.RegisterProcessor("output_stdout", stdout.New)
	runtime.RegisterProcessor("output_mongodb", mongodb.New)
	runtime.RegisterProcessor("output_null", null.New)
	runtime.RegisterProcessor("output_elasticsearch", elasticsearch.New)

	runtime.RegisterProcessor("when", when.New)
	runtime.RegisterProcessor("output_when", when.New)

	// veino.RegisterProcessor("httppoller", httppoller.New)
	// veino.RegisterProcessor("fileinput", fileinput.New)
	// veino.RegisterProcessor("imap_input", imap_input.New)
	// veino.RegisterProcessor("file-output", fileoutput.New)
}
