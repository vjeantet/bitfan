//go:generate bitfanDoc
package udpinput

import (
	"fmt"
	"net"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// UDP port number to listen on
	Port int `mapstructure:"port"`

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	Codec string `mapstructure:"codec"`
}

type processor struct {
	processors.Base

	opt *options
	uc  *net.UDPConn
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Port: 514,
	}
	p.opt = &defaults

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func makeUdpServer(port int) (*net.UDPConn, error) {
	var err error

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return sock, nil
}

func (p *processor) Start(e processors.IPacket) error {
	udpSock, err := makeUdpServer(p.opt.Port)
	if err != nil {
		p.Logger.Errorf("Could not start UDP input")
		return err
	}
	p.uc = udpSock

	go func(p *processor, conn *net.UDPConn) {
		buf := make([]byte, 65536)

		for {
			buflen, saddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				p.Logger.Errorf("ReadFromUDP: %s input-udp goroutine exiting", err)
				return
			}
			ne := p.NewPacket(string(buf[:buflen]), map[string]interface{}{
				"host": saddr.IP.String(),
			})

			p.opt.ProcessCommonOptions(ne.Fields())
			p.Send(ne)
		}
	}(p, p.uc)

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	if p.uc != nil {
		p.uc.Close()
	}
	return nil
}
