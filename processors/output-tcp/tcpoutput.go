//go:generate bitfanDoc

package tcpoutput

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/vjeantet/bitfan/codecs"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	conn net.Conn
	processors.Base
	enc codecs.Encoder
	opt *options
}

type options struct {
	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @Default "line"
	// @Enum "json","line","pp","rubydebug"
	// @Type codec
	Codec codecs.CodecCollection `mapstructure:"codec"`

	Host string `mapstructure:"host" validate:"required"`

	Port uint `mapstructure:"port" validate:"required"`

	// Turn this on to enable HTTP keepalive support. Default value is true
	// @Default true
	KeepAlive bool `mapstructure:"keepalive"`

	// Timeout (in seconds) for the entire request. Default value is 60
	// @Default 30
	RequestTimeout uint `mapstructure:"request_timeout"`

	// @Default 10
	RetryInterval uint `mapstructure:"retry_interval"`

	// Add any number of arbitrary tags to your event. There is no default value for this setting.
	// This can help with processing later. Tags can be dynamic and include parts of the event using the %{field} syntax.
	// Tags []string `mapstructure:"tags"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		KeepAlive:      true,
		RequestTimeout: 30,
		RetryInterval:  10,
		Codec: codecs.CodecCollection{
			Enc: codecs.New("line", nil, ctx.Log(), ctx.ConfigWorkingLocation()),
		},
	}
	p.opt = &defaults

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	if err := p.connect(); err != nil {
		time.Sleep(time.Duration(p.opt.RetryInterval) * time.Second)
		return err
	}

	var body bytes.Buffer
	writer := bufio.NewWriter(&body)
	enc, err := p.opt.Codec.NewEncoder(writer)
	if err != nil {
		return fmt.Errorf("Codec failed with: %v", err)
	}
	if err := enc.Encode(e.Fields().Old()); err != nil {
		return fmt.Errorf("Can't encode item with error: %v", err)
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	p.conn.SetDeadline(time.Now().Add(time.Duration(p.opt.RequestTimeout) * time.Second))
	if _, err := p.conn.Write(body.Bytes()); err != nil {
		p.conn.Close()
		p.conn = nil
		return err
	}
	return nil
}

func (p *processor) Start(e processors.IPacket) error {
	return p.connect()
}

func (p *processor) Stop(e processors.IPacket) error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *processor) connect() error {
	var (
		addr *net.TCPAddr
		err  error
	)
	if p.conn == nil {
		if addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.opt.Host, p.opt.Port)); err != nil {
			return err
		}
		if p.conn, err = net.DialTCP("tcp", nil, addr); err != nil {
			return err
		}
		p.conn.(*net.TCPConn).SetNoDelay(false)
		p.conn.(*net.TCPConn).SetKeepAlive(p.opt.KeepAlive)
		p.conn.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	}
	return nil
}
