//go:generate bitfanDoc
package imap_input

import (
	"github.com/etrepat/postman/watch"
	"bitfan/processors"
)

func New() processors.Processor {
	return &processor{}
}

type processor struct {
	processors.Base

	config  *watch.Flags
	watcher *watch.Watch
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.config = watch.NewFlags()
	p.config.Host = conf["host"].(string)
	p.config.Port = uint(conf["port"].(float64))
	p.config.Ssl = conf["ssl"].(bool)
	p.config.Mailbox = conf["mailbox"].(string)
	p.config.Password = conf["password"].(string)
	p.config.Username = conf["username"].(string)
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.Logger.Debug("imap input - closing connection...")
	p.watcher.Stop()

	return nil
}

func (p *processor) Start(e processors.IPacket) error {
	p.watcher = watch.New(p.config, newToJsonHandler(p.NewPacket, e, p.Send))
	go p.watcher.Start()
	return nil
}
