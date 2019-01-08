//go:generate bitfanDoc
package execinput

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/awillis/bitfan/codecs"
	"github.com/awillis/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	Command  string
	Args     []string
	Interval string

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @Type Codec
	// @Default "plain"
	Codec codecs.CodecCollection `mapstructure:"codec"`
}

type processor struct {
	processors.Base

	opt  *options
	q    chan bool
	host string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec: codecs.CodecCollection{
			Dec: codecs.New("plain", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
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
	defer func() {
		pr.Close()
		pw.Close()
	}()

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
			var record interface{}
			if err := dec.Decode(&record); err != nil {
				if err == io.EOF {
					p.Logger.Debugln("error while exec docoding : ", err)
					return nil
				} else {
					p.Logger.Errorln("error while exec docoding : ", err)
					return err
				}
			} else {
				ne := p.NewPacket(map[string]interface{}{
					"message": data,
					"host":    p.host,
				})
				ne.Fields().SetValueForPath(record, "stdout")
				ne.Fields().SetValueForPath(p.opt.Command, "command")
				ne.Fields().SetValueForPath(strings.Join(p.opt.Args, ", "), "args")

				p.opt.ProcessCommonOptions(ne.Fields())
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
	// -----

	if err != nil {
		return fmt.Errorf("Error while executing command '%s' (%v)", p.opt.Command, err)
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
