// +build extra

package glusterfsoutput

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"bitfan/processors"
	"github.com/jehiah/go-strftime"
	"github.com/kshlm/gogfapi/gfapi"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt    *options
	buffer bytes.Buffer
	file   *gfapi.File
	volume *gfapi.Volume
}

type options struct {
	Codec           string        `mapstructure:"codec"`
	CreateIfDeleted bool          `mapstructure:"create_if_deleted"`
	DirMode         os.FileMode   `mapstructure:"dir_mode"`
	FileMode        os.FileMode   `mapstructure:"file_mode"`
	FlushInterval   time.Duration `mapstructure:"flush_interval"`
	Host            string        `mapstructure:"host"`
	Path            string        `mapstructure:"path" validate:"required"`
	Volume          string        `mapstructure:"volume" validate:"required"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Codec:           "json_lines",
		FileMode:        0640,
		DirMode:         0750,
		CreateIfDeleted: true,
		Host:            "localhost",
	}
	p.opt = &defaults
	p.volume = new(gfapi.Volume)
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	var (
		eventBytes []byte
		err        error
	)
	switch p.opt.Codec {
	default:
		p.Logger.Errorf("GlusterFS output: invalid codec '%s', using 'json_lines instead'", p.opt.Codec)
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
	if err := p.volume.Init(p.opt.Host, p.opt.Volume); err != 0 {
		return fmt.Errorf("GlusterFS output: Failed to initialize volume '%s'", p.opt.Volume)
	}
	if err := p.volume.Mount(); err != 0 {
		return fmt.Errorf("GlusterFS output: Failed to mount volume '%s'", p.opt.Volume)
	}

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
		return nil
	}

	filepath := strftime.Format(p.opt.Path, time.Now())
	err = p.volume.MkdirAll(path.Dir(filepath), p.opt.DirMode)
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
		if _, statErr := p.volume.Stat(filepath); os.IsNotExist(statErr) {
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
	p.file, err = p.volume.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, p.opt.FileMode)
	//p.file.Chmod(p.opt.FileMode)
	return err
}

func (p *processor) Stop(e processors.IPacket) error {
	p.writeToFile()
	p.file.Close()
	p.volume.Unmount()
	return nil
}
