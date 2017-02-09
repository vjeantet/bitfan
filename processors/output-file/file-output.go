//go:generate bitfanDoc
package fileoutput

import (
	"bytes"
	"os"
	"path"
	"time"

	"github.com/jehiah/go-strftime"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt    *options
	buffer bytes.Buffer
	file   *os.File
}

type options struct {
	Codec           string        `mapstructure:"codec"`
	CreateIfDeleted bool          `mapstructure:"create_if_deleted"`
	DirMode         os.FileMode   `mapstructure:"dir_mode"`
	FileMode        os.FileMode   `mapstructure:"file_mode"`
	FlushInterval   time.Duration `mapstructure:"flush_interval"`
	Path            string        `mapstructure:"path" validate:"required"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec:           "json_lines",
		CreateIfDeleted: true,
		DirMode:         0750,
		FileMode:        0640,
		FlushInterval:   2,
	}
	p.opt = &defaults
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	var (
		eventBytes []byte
		err        error
	)
	switch p.opt.Codec {
	default:
		p.Logger.Errorf("File output: invalid codec '%s', using 'json_lines instead'", p.opt.Codec)
		fallthrough
	case "json_lines":
		eventBytes, err = e.Fields().Json(true)
	case "json":
		eventBytes, err = e.Fields().JsonIndent("", "  ", true)
	case "xml_lines":
		eventBytes, err = e.Fields().Xml()
	case "xml":
		eventBytes, err = e.Fields().XmlIndent("", "  ")

	}
	p.buffer.Write(eventBytes)
	p.buffer.WriteRune('\n')
	return err
}

func (p *processor) Start(e processors.IPacket) error {
	ticker := time.NewTicker(p.opt.FlushInterval * time.Second)
	go func() {
		for range ticker.C {
			p.writeToFile()
		}
	}()
	return nil
}

func (p *processor) writeToFile() (err error) {
	if p.buffer.Len() == 0 {
		return
	}

	filepath := strftime.Format(p.opt.Path, time.Now())
	err = os.MkdirAll(path.Dir(filepath), p.opt.DirMode)
	if err != nil {
		return err
	}

	switch {
	// file is not yet opened
	case p.file == nil:
		err = p.openFile(filepath)
	// filename changed
	case p.file.Name() != filepath:
		p.file.Close()
		err = p.openFile(filepath)
	// filename didn't change but create_if_deleted is true
	case p.opt.CreateIfDeleted && p.file.Name() == filepath:
		if _, statErr := os.Stat(filepath); os.IsNotExist(statErr) {
			p.file.Close()
			err = p.openFile(filepath)
		}
	}

	if err != nil {
		return err
	}

	bufferBytes := p.buffer.Bytes()
	p.buffer.Read(bufferBytes)
	_, err = p.file.Write(bufferBytes)
	return err
}

func (p *processor) openFile(filepath string) (err error) {
	p.file, err = os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, p.opt.FileMode)
	return err
}

func (p *processor) Stop(e processors.IPacket) error {
	p.writeToFile()
	p.file.Close()
	return nil
}
