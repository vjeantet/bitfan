//go:generate bitfanDoc
package execinput

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	Command   string
	Args      []string
	Add_field map[string]interface{}
	Interval  string

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @Type Codec
	// @Default "plain"
	Codec codecs.Codec `mapstructure:"codec"`

	Tags []string
	Type string
}

type processor struct {
	processors.Base

	opt  *options
	q    chan bool
	host string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec: codecs.New("plain", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
	}
	p.opt = &defaults

	p.host = fqdn.Get()

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {
	p.q = make(chan bool)
	return nil
}
func (p *processor) Stop(e processors.IPacket) error {
	close(p.q)
	return nil
}

func (p *processor) Tick(e processors.IPacket) error {

	var (
		err  error
		data string
	)

	// data, err = p.doExec()

	// ----- EXEC

	var dec codecs.Decoder
	pr, pw := io.Pipe()

	if dec, err = p.opt.Codec.NewDecoder(pr); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		return err
	}

	cmd := exec.Command(p.opt.Command, p.opt.Args...)
	cmd.Stdout = pw
	cmd.Start()

	go func() error {
		defer p.Logger.Debugln("exiting loop")
		for dec.More() {
			if record, err := dec.Decode(); err != nil {
				return err
				break
			} else if record == nil {
				p.Logger.Debugln("waiting for more content...")
				continue
			} else {
				ne := p.NewPacket(data, map[string]interface{}{
					"host": p.host,
				})
				ne.Fields().SetValueForPath(record, "stdout")
				ne.Fields().SetValueForPath(p.opt.Command, "command")
				ne.Fields().SetValueForPath(strings.Join(p.opt.Args, ", "), "args")

				processors.ProcessCommonFields(ne.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
				p.Send(ne)
			}

			select {
			case <-p.q:
				return nil
			default:
			}
		}
		return nil
	}()

	err = cmd.Wait()
	pw.Close()

	// -----

	if err != nil {
		return fmt.Errorf("Error while executing command '%s' (%s)", p.opt.Command, err.Error())
	}

	return nil
}

// func (p *processor) doExec() (data string, err error) {
// 	var (
// 		buferr bytes.Buffer
// 		raw    []byte
// 		cmd    *exec.Cmd
// 	)
// 	cmd = exec.Command(p.opt.Command, p.opt.Args...)
// 	cmd.Stderr = &buferr
// 	if raw, err = cmd.Output(); err != nil {
// 		return
// 	}
// 	data = string(raw)
// 	if buferr.Len() > 0 {
// 		err = errors.New(buferr.String())
// 	}
// 	return
// }
