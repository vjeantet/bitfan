//go:generate bitfanDoc
package tcpinput

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"bitfan/processors"
)

func New() processors.Processor {
	return &processor{
		opt:       &options{},
		wg:        new(sync.WaitGroup),
		start:     make(chan *net.TCPConn, 512),
		end:       make(chan *net.TCPConn, 512),
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
	wg        *sync.WaitGroup
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

	var err error

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p.opt.Port))
	if err != nil {
		p.Logger.Errorf("could not resolve tcp socket address: %s", err.Error())
		return err
	}

	p.sock, err = net.ListenTCP("tcp", addr)
	if err != nil {
		p.Logger.Errorf("could not start TCP input: %s", err.Error())
		return err
	}

	err = p.sock.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		p.Logger.Errorf("could not set socket accept deadline", err)
	}

	go func(p *processor) {
		p.wg.Add(1)
		defer p.wg.Done()
		for {
			conn, err := p.sock.AcceptTCP()

			if err != nil {
				if strings.Contains(err.Error(), "accept tcp") {
					if err = p.sock.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
						p.Logger.Error(err)
					}
					continue
				}
				p.Logger.Infof("shutting down tcp acceptor: %v", err)
				break
			}

			if err := conn.SetReadBuffer(p.opt.ReadBufferSize); err != nil {
				p.Logger.Error(err)
			}
			p.conntable.Store(conn.RemoteAddr().String(), *conn)
			p.start <- conn
		}
		close(p.start)
	}(p)

	go func(p *processor) {
		p.wg.Add(1)
		defer p.wg.Done()
		for {
			select {
			case conn, ok := <-p.end:
				if ok {
					if err := conn.Close(); err != nil {
						p.Logger.Error(err)
					} else {
						p.conntable.Delete(conn.RemoteAddr().String())
					}
				} else {
					break
				}
			}
		}
	}(p)

	go func() {
		p.wg.Add(1)
		defer p.wg.Done()
		for {
			select {
			case conn, ok := <-p.start:
				if ok {
					go func(p *processor) {
						p.wg.Add(1)
						defer p.wg.Done()

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
				} else {
					close(p.end)
				}
			}
		}
	}()

	return err
}

func (p *processor) Stop(e processors.IPacket) error {

	var err error

	if p.sock != nil {
		err = p.sock.Close()
	}

	p.wg.Wait()
	return err
}
