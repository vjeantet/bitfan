//go:generate bitfanDoc
package tcpinput

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{
		opt:       &options{},
		start:     make(chan *net.TCPConn),
		end:       make(chan *net.TCPConn),
		conntable: new(sync.Map),
	}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// TCP port number to listen on
	Port int `mapstructure:"port"`
	// Message buffer size
	ReadBufferSize int `mapstructure:"read_buffer_size"`
}

type processor struct {
	processors.Base

	opt       *options
	sock      *net.TCPListener
	start     chan *net.TCPConn
	end       chan *net.TCPConn
	conntable *sync.Map
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Port:           5151,
		ReadBufferSize: 65536,
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
					if err = p.sock.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
						p.Logger.Error(err)
					}
				} else {
					p.Logger.Error(err)
				}
				continue
			}

			if err := conn.SetReadBuffer(p.opt.ReadBufferSize); err != nil {
				p.Logger.Error(err)
			}
			p.conntable.Store(conn.RemoteAddr().String(), *conn)
			p.start <- conn

		}
	}(p)

	go func(p *processor) {
		for {
			conn := <-p.end
			p.conntable.Delete(conn.RemoteAddr().String())
			if err := conn.Close(); err != nil {
				p.Logger.Error(err)
			}
		}
	}(p)

	go func() {
		for {
			select {
			case conn := <-p.start:
				go func(p *processor) {

					buf := bufio.NewReader(conn)
					scanner := bufio.NewScanner(buf)

					for scanner.Scan() {
						ne := p.NewPacket(map[string]interface{}{
							"message": scanner.Text(),
							"host":    conn.LocalAddr().String(),
						})
						p.opt.ProcessCommonOptions(ne.Fields())
						p.Send(ne)
					}
					if err := scanner.Err(); err != nil {
						p.Logger.Error(err)
					}
					p.end <- conn
				}(p)

			}
		}
	}()

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {

	var err error

	p.conntable.Range(func(key, value interface{}) bool {
		if err = value.(*net.TCPConn).Close(); err != nil {
			p.Logger.Error(err)
			return false
		} else {
			return true
		}
	})

	if p.sock != nil {
		err = p.sock.Close()
	}
	return err
}
