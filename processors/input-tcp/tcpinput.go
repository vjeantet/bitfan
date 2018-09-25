//go:generate bitfanDoc
package tcpinput

import (
	"bytes"
	"fmt"
	"github.com/vjeantet/bitfan/processors"
	"net"
	"strings"
	"time"
)

func New() processors.Processor {
	return &processor{
		opt:       &options{},
		start:     make(chan *net.TCPConn),
		end:       make(chan *net.TCPConn),
		conntable: make(map[*net.TCPConn]bool),
	}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// TCP port number to listen on
	Port int `mapstructure:"port"`
	// Message buffer size
	BufferSize int `mapstructure:"buffer_size"`
}

type processor struct {
	processors.Base

	opt       *options
	sock      *net.TCPListener
	start     chan *net.TCPConn
	end       chan *net.TCPConn
	conntable map[*net.TCPConn]bool
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Port:       5151,
		BufferSize: 65536,
	}
	p.opt = &defaults

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p.opt.Port))
	if err != nil {
		p.Logger.Errorf("Could not resolve tcp socket address: %s", err.Error())
		return err
	}

	p.sock, err = net.ListenTCP("tcp", addr)
	if err != nil {
		p.Logger.Errorf("Could not start TCP input: %s", err.Error())
		return err
	}

	err = p.sock.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		p.Logger.Error(err)
	}

	go func(p *processor) {
		for {
			conn, err := p.sock.AcceptTCP()

			if err != nil {
				if strings.Contains(err.Error(), "accept tcp") {
					p.sock.SetDeadline(time.Now().Add(3 * time.Second))
				} else {
					p.Logger.Error(err)
				}
				continue
			}
			p.conntable[conn] = true
			p.start <- conn

		}
	}(p)

	go func(p *processor) {
		for {
			conn := <-p.end
			conn.Close()
		}
	}(p)

	go func() {
		for {
			select {
			case conn := <-p.start:
				go func(p *processor) {

					buf := bytes.NewBuffer(make([]byte, p.opt.BufferSize))
					length, err := conn.ReadFrom(buf)

					if err != nil {
						p.Logger.Errorln(err)
					}

					p.end <- conn

					if length == 0 {
						p.Logger.Print("No data read from socket")
						return
					}

					ne := p.NewPacket(map[string]interface{}{
						"message": buf.String(),
						"host":    addr.String(),
					})

					p.opt.ProcessCommonOptions(ne.Fields())
					p.Send(ne)
				}(p)

			}
		}
	}()

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {

	if p.sock != nil {
		err := p.sock.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
