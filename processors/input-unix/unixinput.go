//go:generate bitfanDoc
package unixinput

import (
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"bitfan/processors"
	"github.com/clbanning/mxj"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The read timeout in seconds. If a particular connection is idle for more than this timeout period, we will assume it is dead and close it.
	// If you never want to timeout, use 0.
	// Default value is 0
	DataTimeout time.Duration `mapstructure:"data_timeout"`

	// Remove socket file in case of EADDRINUSE failure
	// Default value is false
	ForceUnlink bool `mapstructure:"force_unlink"`

	// Mode to operate in. server listens for client connections, client connects to a server.
	// Value can be any of: "server", "client"
	// Default value is "server"
	Mode string `mapstructure:"mode"`

	// When mode is server, the path to listen on. When mode is client, the path to connect to.
	Path string `mapstructure:"path" validate:"required"`

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	Codec string `mapstructure:"codec"`
}

type processor struct {
	processors.Base

	opt *options
	ln  *net.UnixListener
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		DataTimeout: 0,
		ForceUnlink: false,
		Mode:        "server",
		Codec:       "line",
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	switch p.opt.Codec {
	case "line", "json", "xml":
	default:
		return fmt.Errorf("invalid codec '%s'. Valid codecs are: 'line', 'json' and 'xml'", p.opt.Codec)
	}
	return err
}

func (p *processor) startServer() (err error) {
	p.ln, err = net.ListenUnix("unix", &net.UnixAddr{Name: p.opt.Path, Net: "unix"})

	if isAddrInUse(err) {
		if p.opt.ForceUnlink {
			os.Remove(p.opt.Path)
			p.ln, err = net.ListenUnix("unix", &net.UnixAddr{Name: p.opt.Path, Net: "unix"})
		} else {
			return fmt.Errorf("could not start server: %v", err)
		}
	}

	if err != nil {
		return fmt.Errorf("could not start server: %v", err)
	}

	return err
}

func (p *processor) Start(e processors.IPacket) error {

	switch p.opt.Mode {
	case "server":
		if err := p.startServer(); err != nil {
			return err
		}
		go func() {
			for {
				conn, err := p.ln.AcceptUnix()
				if err != nil {
					netErr, ok := err.(net.Error)
					//If this is a timeout, then continue to wait for new connections
					if ok && netErr.Timeout() && netErr.Temporary() {
						continue
					}
				}
				if p.opt.DataTimeout > 0 {
					conn.SetReadDeadline(time.Now().Add(p.opt.DataTimeout * time.Second))
				}
				go p.parse(conn)
			}
		}()
	case "client":
		go func() {
			for {
				conn, err := net.Dial("unix", p.opt.Path)
				if err != nil {
					continue
				}
				p.parse(conn)
			}
		}()
	default:
		return fmt.Errorf("Unrecognized mode: %s. Must be either 'server' or 'client'.", p.opt.Mode)
	}

	return nil
}

func isAddrInUse(err error) bool {
	if err, ok := err.(*net.OpError); ok {
		if err, ok := err.Err.(*os.SyscallError); ok {
			return err.Err == syscall.EADDRINUSE
		}
	}
	return false
}

func (p *processor) parse(conn net.Conn) {
	defer conn.Close()
	var event processors.IPacket

	switch p.opt.Codec {
	case "line":
		buf := make([]byte, 65536)
		buflen, err := conn.Read(buf)
		if err != nil {
			p.Logger.Errorf(err.Error())
		}
		message := strings.TrimSpace(string(buf[:buflen]))
		event := p.NewPacket(mxj.Map{"message": message})
		p.opt.ProcessCommonOptions(event.Fields())
		p.Send(event)

	case "json":
		json, raw, err := mxj.NewMapJsonReaderRaw(conn)
		if err != nil {
			p.Logger.Errorf(err.Error())
			event = p.NewPacket(mxj.Map{"message": string(raw)})
		} else {
			event = p.NewPacket(json)
		}
		p.opt.ProcessCommonOptions(event.Fields())
		p.Send(event)

	case "xml":
		xml, raw, err := mxj.NewMapXmlReaderRaw(conn)
		if err != nil {
			p.Logger.Errorf(err.Error())
			event = p.NewPacket(mxj.Map{"message": string(raw)})
		} else {
			event = p.NewPacket(xml)
		}
		p.opt.ProcessCommonOptions(event.Fields())
		p.Send(event)
	}
}

func (p *processor) Stop(e processors.IPacket) error {
	if p.opt.Mode == "server" {
		p.ln.Close()
	}
	return nil
}
