//go:generate bitfanDoc
package sysloginput

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mcuadros/go-syslog.v2"

	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	AddField map[string]interface{} `mapstructure:"add_field"`

	// Which format to use to decode syslog messages. Default value is "automatic"
	// Value can be "automatic","rfc3164", "rfc5424" or "rfc6587"
	//
	// Note on "automatic" format: if you don't know which format to select,
	// or have multiple incoming formats, this is the one to go for.
	// There is a theoretical performance penalty (it has to look at a few bytes
	// at the start of the frame), and a risk that you may parse things you don't want to parse
	// (rogue syslog clients using other formats), so if you can be absolutely sure of your syslog
	// format, it would be best to select it explicitly.
	Format string `mapstructure:"format"`

	// Port number to listen on
	Port int `mapstructure:"port"`

	// Protocol to use to listen to syslog messages
	// Value can be either "tcp" or "udp"
	Protocol string `mapstructure:"protocol"`

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	Tags []string `mapstructure:"tags"`

	// Add a type field to all events handled by this input
	Type string `mapstructure:"type"`
}

type processor struct {
	processors.Base

	opt *options

	s  *syslog.Server
	ch syslog.LogPartsChannel
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Format:   "automatic",
		Protocol: "udp",
	}
	p.opt = &defaults

	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	return nil
}

func (p *processor) Start(e processors.IPacket) error {
	p.ch = make(syslog.LogPartsChannel)
	p.s = syslog.NewServer()

	handler := syslog.NewChannelHandler(p.ch)

	switch p.opt.Format {
	case "automatic":
		p.s.SetFormat(syslog.Automatic)
	case "rfc3164":
		p.s.SetFormat(syslog.RFC3164)
	case "rfc5424":
		p.s.SetFormat(syslog.RFC5424)
	case "rfc6587":
		p.s.SetFormat(syslog.RFC6587)
	default:
		return fmt.Errorf("%s is not a valid format", p.opt.Format)
	}

	p.s.SetHandler(handler)

	switch strings.ToLower(p.opt.Protocol) {
	case "udp":
		p.s.ListenUDP(fmt.Sprintf(":%d", p.opt.Port))
	case "tcp":
		p.s.ListenTCP(fmt.Sprintf(":%d", p.opt.Port))
	default:
		return fmt.Errorf("%s is not a valid protocol", p.opt.Protocol)
	}

	p.s.Boot()

	go func(channel syslog.LogPartsChannel) {
		for message := range channel {
			// Use syslog timestamp as @timestamp field, with correct format
			message["@timestamp"] = message["timestamp"].(time.Time)
			delete(message, "timestamp")

			ne := p.NewPacket("", message)
			processors.ProcessCommonFields(ne.Fields(), p.opt.AddField, p.opt.Tags, p.opt.Type)
			p.Send(ne)
		}
	}(p.ch)

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	if p.s != nil {
		p.s.Kill()
	}
	return nil
}
