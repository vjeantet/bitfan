package xprocessor

import (
	"encoding/json"
	"log"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var Envs map[string]string

var Logger = log.New(os.Stderr, "", 0)

func init() {
	Envs = map[string]string{
		"PIPELINE_UUID":         os.Getenv("BF_PIPELINE_UUID"),
		"PIPELINE_WORKING_PATH": os.Getenv("BF_PIPELINE_WORKING_PATH"),
		"PROCESSOR_DATA_PATH":   os.Getenv("BF_PROCESSOR_DATA_PATH"),
		"PROCESSOR_NAME":        os.Getenv("BF_PROCESSOR_NAME"),
		"PROCESSOR_LABEL":       os.Getenv("BF_PROCESSOR_LABEL"),
	}
}

type Runner struct {
	opt Options

	Description      string
	ShortDescription string

	configure func(Options) error
	start     func() error
	receive   func(interface{}) error
	stop      func() error

	enc *json.Encoder // TODO use interface
	//dec *decoder // TODO @see Run
}

func New(
	configure func(Options) error,
	start func() error,
	receive func(interface{}) error,
	stop func() error,
) *Runner {
	return &Runner{
		opt:       Options{},
		configure: configure,
		start:     start,
		receive:   receive,
		stop:      stop,
	}
}

func (r *Runner) option(name string, required bool, doc string, defaultValue interface{}) *Option {
	return &Option{
		Name:         name,
		Alias:        name,
		Doc:          doc,
		Required:     required,
		DefaultValue: defaultValue,
	}
}
func (r *Runner) OptionBool(name string, required bool, doc string, defaultValue bool) {
	r.opt[name] = r.option(name, required, doc, defaultValue)
	r.opt[name].Type = "bool"
}
func (r *Runner) OptionString(name string, required bool, doc string, defaultValue string) {
	r.opt[name] = r.option(name, required, doc, defaultValue)
	r.opt[name].Type = "string"
}
func (r *Runner) OptionInt(name string, required bool, doc string, defaultValue int) {
	r.opt[name] = r.option(name, required, doc, defaultValue)
	r.opt[name].Type = "int"
}
func (r *Runner) OptionStringSlice(name string, required bool, doc string, defaultValue []string) {
	r.opt[name] = r.option(name, required, doc, defaultValue)
	r.opt[name].Type = "[]string"
}
func (r *Runner) OptionIntSlice(name string, required bool, doc string, defaultValue []int) {
	r.opt[name] = r.option(name, required, doc, defaultValue)
	r.opt[name].Type = "[]int"
}

func (r *Runner) Logf(format string, args ...interface{}) {
	Logger.Printf(format, args...)
}

func (t *Runner) Run(maxConcurrent int) {

	// Spec Options
	for _, spec := range t.opt {
		varName := spec.Name
		f := kingpin.Flag(varName, spec.Doc)
		if spec.Required {
			f.Required()
		}

		switch spec.Type {
		case "string":
			spec.Value = f.String()
		case "bool":
			spec.Value = f.Bool()
		case "int":
			spec.Value = f.Int()
		case "[]string":
			spec.Value = f.Strings()
		case "[]int":
			spec.Value = f.Ints()
		}
	}

	kingpin.Flag("help-json", "output spec as json").Bool()

	action := func(pc *kingpin.ParseContext) error {
		for _, pe := range pc.Elements {
			switch v := pe.Clause.(type) {
			case *kingpin.FlagClause:
				if v.Model().Name == "help-json" {
					mapB, _ := json.Marshal(
						struct {
							Description      string
							ShortDescription string
							Options          Options
						}{
							Description:      t.Description,
							ShortDescription: t.ShortDescription,
							Options:          t.opt,
						},
					)
					os.Stdout.Write(mapB)
					os.Exit(0)
				}
			}

		}
		return nil
	}

	kingpin.CommandLine.PreAction(action)
	kingpin.Parse()

	// pp.Println("t-->", t.opt)

	// Configure Processor with flags values
	if t.configure != nil {
		err := t.configure(t.opt)
		if err != nil {
			Logger.Println(err)
			return
		}
	}

	// Start processor
	if t.start != nil {
		t.start()
	}

	t.enc = json.NewEncoder(os.Stdout)
	// TODO : pass decoder to processStdinDataWith
	processStdinDataWith(t.receive, maxConcurrent).Wait()

	// Stop
	if t.stop != nil {
		t.stop()
	}
}

func (t *Runner) Send(data map[string]interface{}) error {
	return t.enc.Encode(data)
}
