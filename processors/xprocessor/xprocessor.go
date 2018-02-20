package xprocessor

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/mitchellh/mapstructure"
	"github.com/vjeantet/bitfan/api/models"
	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

func NewWithSpec(spec *models.XProcessor) processors.Processor {
	p := &processor{opt: &options{
		Behavior: spec.Behavior,
		Stream:   spec.Stream,
		Command:  spec.Command,
		Args:     spec.Args,
		Code:     spec.Code,
		StdinAs:  spec.StdinAs,
		StdoutAs: spec.StdoutAs,
	}}

	if p.opt.Command == "php" && len(p.opt.Args) == 0 {
		if strings.HasPrefix(p.opt.Code, "<?php") {
			p.opt.Code = p.opt.Code[5:]
		}
		p.opt.Args = []string{"-d", "display_errors=stderr", "-r", p.opt.Code, "--"}
		//var_dump(getopt("",array("count:")));' -- --count=1

	}
	return p
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

	// What is the value's format of stdinputed value
	// @Default "json","none"
	StdinAs string `mapstructure:"stdin_as" validate:"required"`

	// What is the value's format of stdoutputed value
	// @Enum "json","string"
	// @Default "json"
	StdoutAs string `mapstructure:"stdout_as" validate:"required"`

	// Flags for delegated processors (will be passed as args)
	// @default {"Content-Type" => "application/json"}
	Flags map[string]string
}

// Reads events from standard input
type processor struct {
	processors.Base

	opt *options

	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

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
	delete(conf, "interval")

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

	return err
}

func (p *processor) Start(e processors.IPacket) error {
	p.wg = new(sync.WaitGroup)
	if p.opt.Stream == false {
		return nil
	}
	p.Logger.Infof("Start %s %s", p.opt.Behavior, p.opt.Stream)
	var err error
	p.cmd, p.stdin, p.stdout, p.stderr, err = p.startCommand(nil)

	go func(s io.ReadCloser) {
		scanner := bufio.NewScanner(s)
		for scanner.Scan() {
			p.Logger.Errorf("stderr : %s", scanner.Text())
		}
	}(p.stderr)

	if p.opt.Behavior == PRODUCER || p.opt.Behavior == TRANSFORMER {
		var dec codecs.Decoder
		if dec, err = p.opt.Codec.NewDecoder(p.stdout); err != nil {
			p.Logger.Errorln("decoder error : ", err.Error())
			return err
		}
		// READ FROM PROC OUTPUT AND SEND EVENTS
		go p.readAndSendEventsFromProcess(dec)
	}

	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	if p.opt.Stream == false {
		return p.Receive(e)
	}
	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	// if p.opt.Behavior == PRODUCER {
	// 	return nil
	// }

	if p.opt.Stream == true {
		return p.writeEventToProcess(e)
	} else {
		p.wg.Add(1)
		defer p.wg.Done()
		return p.runProcess(e)
	}

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	if p.opt.Stream == true {
		return p.stopStream(e)
	}
	p.wg.Wait()
	return nil
}

func (p *processor) buildCommandArgs(e processors.IPacket) []string {
	args := []string{}
	for _, v := range p.opt.Args {
		args = append(args, v)
	}
	for k, v := range p.opt.Flags {
		if k == "_" {
			continue
		}
		if v == "" {
			args = append(args, "--"+k)
		} else {
			if e != nil {
				processors.Dynamic(&v, e.Fields())
			}
			args = append(args, "--"+k+"="+v)
		}
	}
	if v, ok := p.opt.Flags["_"]; ok {
		args = append(args, v)
	}
	return args
}

func (p *processor) buildEnv() []string {
	env := os.Environ()
	env = append(env,
		fmt.Sprintf("BF_PIPELINE_UUID=%s", p.PipelineUUID),
		fmt.Sprintf("BF_PIPELINE_WORKING_PATH=%s", p.ConfigWorkingLocation),
		fmt.Sprintf("BF_PROCESSOR_DATA_PATH=%s", p.DataLocation),
		fmt.Sprintf("BF_PROCESSOR_NAME=%s", p.B().Name),
		fmt.Sprintf("BF_PROCESSOR_LABEL=%s", p.B().Label),
	)
	return env
}

func (p *processor) startCommand(e processors.IPacket) (*exec.Cmd, io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	args := p.buildCommandArgs(e)

	var cmd *exec.Cmd

	p.Logger.Debugf("command '%s', args=%s", p.opt.Command, args)
	cmd = exec.Command(p.opt.Command, args...)
	cmd.Dir = p.ConfigWorkingLocation
	cmd.Env = p.buildEnv()

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return cmd, stdin, stdout, stderr, err
	}

	p.Logger.Infof("delegator %s started PId=%d", p.B().Name, cmd.Process.Pid)
	return cmd, stdin, stdout, stderr, nil
}

func (p *processor) runProcess(e processors.IPacket) error {

	cmd, stdin, stdout, stderr, err := p.startCommand(e)

	// Encode received event
	var enc codecs.Encoder
	enc, err = p.opt.Codec.NewEncoder(stdin)
	if err != nil {
		p.Logger.Errorln("codec error : ", err.Error())
		return err
	}
	enc.Encode(e.Fields().Old())
	stdin.Close()

	// Decode resulting output
	var dec codecs.Decoder
	if dec, err = p.opt.Codec.NewDecoder(stdout); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		return err
	}
	p.readAndSendEventsFromProcess(dec)

	b, _ := ioutil.ReadAll(stderr)
	if len(b) > 0 {
		p.Logger.Errorf("stderr : %s", b)
	}
	err = cmd.Wait()

	return nil
}

func (p *processor) readAndSendEventsFromProcess(dec codecs.Decoder) error {
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

func (p *processor) writeEventToProcess(e processors.IPacket) error {
	var err error
	// Encode received event
	var enc codecs.Encoder
	enc, err = p.opt.Codec.NewEncoder(p.stdin)
	if err != nil {
		p.Logger.Errorln("codec error : ", err.Error())
		return err
	}
	enc.Encode(e.Fields().Old())

	return nil
}

func (p *processor) stopStream(e processors.IPacket) error {
	if p.cmd != nil {
		p.cmd.Process.Signal(syscall.SIGQUIT)
		p.cmd.Wait()
	}
	return nil
}
