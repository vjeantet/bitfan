//go:generate bitfanDoc
//Execute a command and use its stdout as event data
package exec

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

const (
	PORT_SUCCESS = 0
)

// no concurency ! only one worker
func (p *processor) MaxConcurent() int { return 0 }

// drop event when field value is the same in the last event
type processor struct {
	processors.Base
	opt *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	Command string `mapstructure:"command"  validate:"required"`

	Args []string `mapstructure:"args"`

	// Pass the complete event to stdin as a json string
	// @Default false
	Stdin bool `mapstructure:"stdin"`

	// Where do the output should be stored
	// Set "." when output is json formated and want to replace current event fields with output
	// response. (useful)
	// @Default "stdout"
	Target string `mapstructure:"target"`

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @Type Codec
	// @Default "plain"
	Codec codecs.Codec `mapstructure:"codec"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		Target: "stdout",
		Stdin:  false,
		Codec:  codecs.New("plain", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
	}
	p.opt = &defaults

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	var err error
	jsonBytes := []byte{}
	if p.opt.Stdin {
		jsonBytes, err = e.Fields().Json()
		if err != nil {
			return err
		}
	}

	d, err := p.doExec(jsonBytes, e)
	if err != nil {
		return err
	}

	if p.opt.Target == "." {
		var dat map[string]interface{}
		if err := json.Unmarshal(d, &dat); err != nil {
			processors.AddTags([]string{"_execJsonFailure"}, e.Fields())
			p.Send(e, PORT_SUCCESS)
			return nil
		}

		// recover @timestamp
		dat["@timestamp"], _ = e.Fields().ValueForPath("@timestamp")
		e = p.NewPacket("", dat)
	} else {
		value := strings.TrimSpace(string(d))
		err := e.Fields().SetValueForPath(value, p.opt.Target)
		if err != nil {
			return err
		}
	}

	p.opt.ProcessCommonOptions(e.Fields())
	p.Send(e, PORT_SUCCESS)
	return nil
}

func (p *processor) doExec(inData []byte, e processors.IPacket) (data []byte, err error) {
	var (
		buferr bytes.Buffer
		cmd    *exec.Cmd
	)

	args := []string{}
	for _, a := range p.opt.Args {
		processors.Dynamic(&a, e.Fields())
		args = append(args, a)
	}
	p.Logger.Debugf("command '%s', args=%s", p.opt.Command, args)
	cmd = exec.Command(p.opt.Command, args...)
	cmd.Stderr = &buferr
	stdin, err := cmd.StdinPipe()
	stdout, err := cmd.StdoutPipe()
	cmd.Start()
	if len(inData) > 0 {
		stdin.Write(inData)
		stdin.Close()
	}
	data, _ = ioutil.ReadAll(stdout)
	err = cmd.Wait()

	if buferr.Len() > 0 {
		err = errors.New(buferr.String())
	}

	return
}
