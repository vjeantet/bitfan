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
	// OS socket buffer size
	SockBufferSize int `mapstructure:"sock_buffer_size"`
	// application buffer size, used to read from OS
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
		SockBufferSize: 65536,
		ReadBufferSize: 131072,
	}
	p.opt = &defaults

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {

	var err error

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p.opt.Port))
	if err != nil {
		p.Logger.Errorf("could not resolve tcp socket address: %v", err)
		return err
	}

	p.sock, err = net.ListenTCP("tcp", addr)
	if err != nil {
		p.Logger.Errorf("could not start TCP input: %v", err)
		return err
	}

	err = p.sock.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		p.Logger.Errorf("could not set socket accept deadline: %v", err)
	}

	go func(p *processor) {

		p.wg.Add(1)
		defer close(p.start)
		defer p.wg.Done()

		for {
			conn, err := p.sock.AcceptTCP()

			if err != nil {
				if strings.Contains(err.Error(), "accept tcp") {
					if err = p.sock.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
						p.Logger.Info("shutting down tcp acceptor")
						break
					} else {
						continue
					}
				}
				p.Logger.Errorf("socket error: %v", err)
			}

			if err := conn.SetReadBuffer(p.opt.SockBufferSize); err != nil {
				p.Logger.Errorf("error setting socket buffer size: %v", err)
			}
			p.conntable.Store(conn.RemoteAddr().String(), *conn)
			p.start <- conn
		}
	}(p)

	go func() {

		p.wg.Add(1)
		var shutdown = false
		var waiter = new(sync.WaitGroup)
		defer close(p.end)
		defer p.wg.Done()

		for {
			select {
			case conn, ok := <-p.start:
				if ok {
					go func(p *processor) {
						waiter.Add(1)
						defer waiter.Done()

						scanner := bufio.NewScanner(bufio.NewReader(conn))
						scanner.Buffer(make([]byte, 0, p.opt.ReadBufferSize), p.opt.ReadBufferSize)
						hostname, port, err := net.SplitHostPort(conn.RemoteAddr().String())
						if err != nil {
							p.Logger.Errorf("error getting remote host address")
						}

						for scanner.Scan() {
							ne := p.NewPacket(map[string]interface{}{
								"message":  scanner.Text(),
								"hostname": hostname,
								"port":     port,
							})
							p.opt.ProcessCommonOptions(ne.Fields())
							p.Send(ne)
						}
						if err := scanner.Err(); err != nil {
							p.Logger.Errorf("error while reading from client: %v", err)
						}

						p.end <- conn
					}(p)
				} else {
					shutdown = true
				}
			}
			if shutdown {
				waiter.Wait()
				break
			}
		}
	}()

	go func(p *processor) {

		p.wg.Add(1)
		var shutdown = false
		defer p.wg.Done()

		for {
			select {
			case conn, ok := <-p.end:
				if ok {
					if err := conn.Close(); err != nil {
						p.Logger.Errorf("error while closing connection: %v", err)
					} else {
						p.conntable.Delete(conn.RemoteAddr().String())
					}
				} else {
					shutdown = true
				}
			}
			if shutdown {
				break
			}
		}
	}(p)

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
