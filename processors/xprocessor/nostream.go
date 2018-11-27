package xprocessor

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"bitfan/codecs"
	"bitfan/processors"
)

// Reads events from standard input
type noStreamProcessor struct {
	processor
}

func (p *noStreamProcessor) Start(e processors.IPacket) error {
	p.wg = new(sync.WaitGroup)
	return nil
}

func (p *noStreamProcessor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func (p *noStreamProcessor) Receive(e processors.IPacket) error {
	// p.wg.Add(1)

	cmd, stdin, stdout, stderr, err := p.startCommand(e)

	// Encode received event
	var enc codecs.Encoder
	enc, err = p.opt.Codec.NewEncoder(stdin)
	if err != nil {
		p.Logger.Errorln("codec error : ", err.Error())
		return err
	}
	enc.Encode(e.Fields().Old())
	stdin.Close()

	// Decode resulting output
	var dec codecs.Decoder
	if dec, err = p.opt.Codec.NewDecoder(stdout); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		return err
	}
	p.wg.Add(2)
	go func(s io.ReadCloser) {
		defer p.wg.Done()
		fscanner := bufio.NewScanner(s)
		for fscanner.Scan() {
			if strings.HasPrefix(fscanner.Text(), "[DEBUG] ") {
				p.Logger.Debugf("%s", strings.TrimPrefix(fscanner.Text(), "[DEBUG] "))
			} else {
				p.Logger.Errorf("%s", fscanner.Text())
			}

		}
	}(stderr)
	p.readAndSendEventsFromProcess(dec)

	err = cmd.Wait()
	return err
}

func (p *noStreamProcessor) Stop(e processors.IPacket) error {
	p.wg.Wait()
	return nil
}
