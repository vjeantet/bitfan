package xprocessor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
	"bitfan/api/models"
	"bitfan/codecs"
	"bitfan/processors"
	"bitfan/processors/doc"
)

func NewWithSpec(spec *models.XProcessor) processors.Processor {
	opt := &options{
		Behavior:              spec.Behavior,
		Stream:                spec.Stream,
		Command:               spec.Command,
		Args:                  spec.Args,
		Kind:                  spec.Kind,
		Code:                  spec.Code,
		StdinAs:               spec.StdinAs,
		StdoutAs:              spec.StdoutAs,
		OptionsCompositionTpl: spec.OptionsCompositionTpl,
	}

	if spec.Stream == true {
		p := &streamProcessor{}
		p.opt = opt
		p.hasDoc = spec.HasDoc
		return p
	} else {
		p := &noStreamProcessor{}
		p.opt = opt
		p.hasDoc = spec.HasDoc
		return p
	}
}

const (
	TRANSFORMER string = "transformer"
	CONSUMER    string = "consumer"
	PRODUCER    string = "producer"
)

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	Codec codecs.CodecCollection

	// Producer ? Consumer ? Transformer ?
	// @Enum "producer","transformer","consumer"
	// @Default "transformer"
	Behavior string `mapstructure:"behavior" validate:"required"`

	// Delegated processor is started one time and receives events through its stdin.
	// When it should be started for each received event set value to "false"
	// @Default false
	Stream bool `mapstructure:"stream" `

	// Path to the bin used as delegated processor
	Command string   `mapstructure:"command" validate:"required"`
	Args    []string `mapstructure:"args" `
	Code    string   `mapstructure:"code"`

	Kind string `mapstructure:"kind"`

	// What is the value's format of stdinputed value
	// @Default "json","none"
	StdinAs string `mapstructure:"stdin_as" validate:"required"`

	// What is the value's format of stdoutputed value
	// @Enum "json","string"
	// @Default "json"
	StdoutAs string `mapstructure:"stdout_as" validate:"required"`

	OptionsCompositionTpl string `mapstructure:"options_composition_tpl"`

	// Flags for delegated processors (will be passed as args)
	// @default {"Content-Type" => "application/json"}
	Flags map[string]interface{}
}

// Reads events from standard input
type processor struct {
	processors.Base

	hasDoc bool

	opt *options

	flagsTpl *template.Template

	wg *sync.WaitGroup
	q  chan bool
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	// remove common option from conf
	if err := mapstructure.WeakDecode(conf, &p.opt.CommonOptions); err != nil {
		return err
	}
	delete(conf, "add_field")
	delete(conf, "type")
	delete(conf, "remove_tag")
	delete(conf, "remove_field")
	delete(conf, "add_tag")
	delete(conf, "trace")
	delete(conf, "interval") // todo remove only when producer noStream
	delete(conf, "workers")

	// Set processor's user options
	if err := mapstructure.WeakDecode(conf, &p.opt.Flags); err != nil {
		return err
	}

	p.opt.Codec = codecs.CodecCollection{}
	if p.opt.StdinAs == "json" {
		p.opt.Codec.Enc = codecs.New("json", nil, ctx.Log(), ctx.ConfigWorkingLocation())
	}

	if p.opt.StdinAs == "line" {
		p.opt.Codec.Enc = codecs.New("line", map[string]interface{}{"format": "{{.message}}"}, ctx.Log(), ctx.ConfigWorkingLocation())
	}

	if p.opt.StdoutAs == "json" {
		p.opt.Codec.Dec = codecs.New("json", nil, ctx.Log(), ctx.ConfigWorkingLocation())
	}
	if p.opt.StdoutAs == "line" {
		p.opt.Codec.Dec = codecs.New("line", nil, ctx.Log(), ctx.ConfigWorkingLocation())
	}

	if p.opt.Behavior != TRANSFORMER && p.opt.Behavior != CONSUMER && p.opt.Behavior != PRODUCER {
		return fmt.Errorf("unknow behavior '%s'", p.opt.Behavior)
	}

	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.opt.OptionsCompositionTpl != "" {
		//Prepare templates

		funcMap := template.FuncMap{
			"isSlice": func(a interface{}) (bool, error) {
				av := reflect.ValueOf(a).Kind()
				return av == reflect.Slice, nil
			},
		}
		// p.opt.OptionsCompositionTpl = `{{ range $name, $value := .Options }}{{if ne $name "_" }}{{if isSlice $value}}{{ range $value }}--{{$name}}="{{.}}" {{end}}{{else}}--{{$name}}="{{$value}}" {{end}}{{ end }}{{ end }}{{ index .Options "_"}}`

		tpl, err := template.New("").Option("missingkey=zero").Funcs(funcMap).Parse(p.opt.OptionsCompositionTpl)
		if err != nil {
			return err
		}
		p.flagsTpl = tpl
	}

	switch p.opt.Kind {
	case "php":
		if strings.HasPrefix(p.opt.Code, "<?php") {
			p.opt.Code = p.opt.Code[5:]
		}
		p.opt.Command = "php"
		p.opt.Args = []string{"-d", "display_errors=stderr", "-r", p.opt.Code, "--"}
	case "python":
		p.opt.Command = "python"
		p.opt.Args = []string{"-u", "-c", p.opt.Code}
	case "golang":
		tmpGoFile := filepath.Join(p.DataLocation, "main.go")

		err := ioutil.WriteFile(tmpGoFile, []byte(p.opt.Code), 0644)
		if err != nil {
			return err
		}

		p.opt.Command = "go"
		p.opt.Args = []string{"run", tmpGoFile}
	}

	return err
}

func (p *processor) buildCommandArgs(e processors.IPacket) []string {
	finalArgs := []string{}
	for _, v := range p.opt.Args {
		v := strings.TrimSpace(v)
		if v == "" {
			continue
		}
		finalArgs = append(finalArgs, v)
	}
	if p.flagsTpl == nil {
		for k, v := range p.opt.Flags {
			if k == "_" {
				continue
			}
			switch vt := v.(type) {
			case string:
				if v == "" {
					finalArgs = append(finalArgs, "--"+k)
				} else {
					if e != nil {
						processors.Dynamic(&vt, e.Fields())
					}
					finalArgs = append(finalArgs, "--"+k+"="+vt)
				}
			case int64:
				finalArgs = append(finalArgs, fmt.Sprintf("--%s=%d", k, vt))
			case bool:
				if vt == true {
					finalArgs = append(finalArgs, fmt.Sprintf("--%s", k))
				}
			case []interface{}:
				for _, l := range vt {
					switch lt := l.(type) {
					case string:
						if e != nil {
							processors.Dynamic(&lt, e.Fields())
						}
						finalArgs = append(finalArgs, fmt.Sprintf("--%s=%s", k, lt))
					case int64:
						finalArgs = append(finalArgs, fmt.Sprintf("--%s=%d", k, lt))
					default:
						p.Logger.Errorf("not handled slice type : %s=>%v", k, v)
					}
				}
			case map[string]interface{}:
				for key, value := range vt {
					switch valuet := value.(type) {
					case string:
						finalArgs = append(finalArgs, fmt.Sprintf("--%s=%s:%s", k, key, valuet))
					default:
						p.Logger.Errorf("not handled map value type : %s=>%s", key, value)
					}
				}

			default:
				p.Logger.Errorf("not handled type : %s=>%v", k, v)
			}

		}
		if v, ok := p.opt.Flags["_"]; ok {
			switch vt := v.(type) {
			case string:
				finalArgs = append(finalArgs, vt)
			default:
				p.Logger.Errorf("not handled type : _ =>%v", v)
			}

		}
	} else { // use template
		buff := bytes.NewBufferString("")
		err := p.flagsTpl.Execute(buff, struct{ Options map[string]interface{} }{p.opt.Flags})
		if err != nil {
			p.Logger.Errorf("template error : %v", err)
			return finalArgs
		}
		finalArgs = append(finalArgs, strings.TrimSpace(buff.String()))
	}

	return finalArgs
}
func (p *processor) Doc() *doc.Processor {
	d := &doc.Processor{
		Behavior: p.opt.Behavior,
		Name:     p.Name,
		Doc:      "",
		DocShort: "",
		Options: &doc.ProcessorOptions{
			Options: []*doc.ProcessorOption{},
		},
	}

	if p.hasDoc == false {
		return d
	}

	args := p.buildCommandArgs(nil)
	args = append(args, "--help-json")
	out, err := exec.Command(p.opt.Command, args...).Output()
	if err != nil {
		return d
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(out, &dat); err != nil {
		return d
	}

	d.Doc = dat["Description"].(string)
	d.DocShort = dat["ShortDescription"].(string)

	for _, opti := range dat["Options"].(map[string]interface{}) {
		opt := opti.(map[string]interface{})

		d.Options.Options = append(d.Options.Options, &doc.ProcessorOption{
			Name:         opt["name"].(string),
			Alias:        opt["name"].(string),
			Doc:          opt["doc"].(string),
			Required:     opt["required"].(bool),
			DefaultValue: opt["default_value"],
			Type:         mapType(opt["type"].(string)),
			// PossibleValues: []string{},
			// ExampleLS:      "",
		})
	}

	return d
}

func mapType(s string) string {
	switch s {
	case "[]string":
		return "array"
	case "[]int":
		return "array"
	case "map[string]string":
		return "hash"
	}
	return s
}

func (p *processor) startCommand(e processors.IPacket) (*exec.Cmd, io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	args := p.buildCommandArgs(e)

	var cmd *exec.Cmd

	cmd = exec.Command(p.opt.Command, args...)
	cmd.Dir = p.ConfigWorkingLocation
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("BF_PIPELINE_UUID=%s", p.PipelineUUID),
		fmt.Sprintf("BF_PIPELINE_WORKING_PATH=%s", p.ConfigWorkingLocation),
		fmt.Sprintf("BF_PROCESSOR_DATA_PATH=%s", p.DataLocation),
		fmt.Sprintf("BF_PROCESSOR_NAME=%s", p.B().Name),
		fmt.Sprintf("BF_PROCESSOR_LABEL=%s", p.B().Label),
	)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		p.Logger.Errorf("processor start error [%s] with command='%s' %s", err.Error(), p.opt.Command, args)
		return cmd, stdin, stdout, stderr, err
	}
	p.Logger.Infof("processor started process PID=%d with command= '%s' %s", cmd.Process.Pid, p.opt.Command, args)
	return cmd, stdin, stdout, stderr, nil
}

func (p *processor) readAndSendEventsFromProcess(dec codecs.Decoder) error {
	defer p.wg.Done()
	for {
		var record interface{}
		if err := dec.Decode(&record); err != nil {
			if err == io.EOF {
				p.Logger.Debugf("codec end of file : %s", err.Error())
				break
			} else {
				p.Logger.Errorln("codec error : ", err.Error())
				return err
			}
		}

		var ne processors.IPacket
		switch v := record.(type) {
		case string:
			ne = p.NewPacket(map[string]interface{}{
				"message": v,
			})
		case map[string]interface{}:
			ne = p.NewPacket(v)
		case []interface{}:
			ne = p.NewPacket(map[string]interface{}{
				"data": v,
			})
		default:
			p.Logger.Errorf("Unknow structure %#v", v)
			continue
		}

		p.opt.ProcessCommonOptions(ne.Fields())
		p.Send(ne)
	}
	return nil
}
