//go:generate bitfanDoc
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ShowMax/go-fqdn"
	zglob "github.com/mattn/go-zglob"
	"github.com/vjeantet/bitfan/processors"
	"github.com/vjeantet/bitfan/processors/codec"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	Tags []string

	// Add a type field to all events handled by this input
	Type string

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// @Default "plain"
	// @Type Codec
	Codec codec.Codec `mapstructure:"codec"`

	// How many seconds a file should stay unmodified to be read
	// use this to prevent reading a file while another process is writing into.
	ReadOlder int `mapstructure:"read_older"`

	// How often (in seconds) we expand the filename patterns in the path option
	// to discover new files to watch. Default value is 15
	DiscoverInterval int `mapstructure:"discover_interval"`

	// Exclusions (matched against the filename, not full path).
	// Filename patterns are valid here, too.
	Exclude []string `mapstructure:"exclude"`

	// When the file input discovers a file that was last modified before the
	// specified timespan in seconds, the file is ignored.
	// After it’s discovery, if an ignored file is modified it is no longer ignored
	// and any new data is read.
	// Default value is 86400 (i.e. 24 hours)
	IgnoreOlder int `mapstructure:"ignore_older"`

	// What is the maximum number of file_handles that this input consumes at any one time.
	// Use close_older to close some files if you need to process more files than this number.
	MaxOpenFiles int `mapstructure:"max_open_files"`

	// The path(s) to the file(s) to use as an input.
	// You can use filename patterns here, such as /var/log/*.log.
	// If you use a pattern like /var/log/**/*.log, a recursive search of /var/log
	// will be done for all *.log files.
	// Paths must be absolute and cannot be relative.
	// You may also configure multiple paths.
	Path []string `mapstructure:"path" validate:"required"`

	// Path of the sincedb database file
	// The sincedb database keeps track of the current position of monitored
	// log files that will be written to disk.
	SincedbPath string `mapstructure:"sincedb_path"`
}

type processor struct {
	processors.Base

	opt  *options
	q    chan bool
	q2   chan bool
	host string

	filestoWatch chan string

	sinceDBInfos        *sinceDBInfo
	sinceDBLastSaveTime time.Time
	sinceDBInfosMutex   *sync.Mutex
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		MaxOpenFiles:     5,
		DiscoverInterval: 15,
		ReadOlder:        5,
		SincedbPath:      ".sincedb-readfile.json",
		Codec:            codec.New("plain"),
	}

	p.opt = &defaults
	p.host = fqdn.Get()

	var err error
	err = p.ConfigureAndValidate(ctx, conf, p.opt)

	if false == filepath.IsAbs(p.opt.SincedbPath) {
		p.opt.SincedbPath = filepath.Join(p.DataLocation, p.opt.SincedbPath)
	}
	p.Logger.Debugf("sincedb=", p.opt.SincedbPath)

	// Fix relative paths
	fixedPaths := []string{}
	for _, path := range p.opt.Path {
		if !filepath.IsAbs(path) {
			path = filepath.Join(p.ConfigWorkingLocation, path)
		}
		fixedPaths = append(fixedPaths, path)
	}
	p.opt.Path = fixedPaths

	return err

}

func (p *processor) MaxConcurent() int { return 1 }

func (p *processor) Start(e processors.IPacket) error {

	p.loadSinceDBInfos()
	p.q2 = make(chan bool)
	p.q = make(chan bool)

	p.filestoWatch = make(chan string, p.opt.MaxOpenFiles)
	p.Logger.Debug("Start discovering file looper -> towatch ")
	// Start discovering file looper -> towatch
	if p.opt.DiscoverInterval > 0 {
		go func() {
			err := p.discoverFilesToRead()
			if err != nil {
				p.Logger.Debugf("discover files to read : %s", err)
			}
		}()
	}

	p.Logger.Debug("Start file reader <- towatch")
	// Start file reader <- towatch
	go func() {

		for {
			select {
			case <-p.q2:
				close(p.q)
				return
			case filepath := <-p.filestoWatch:
				p.Logger.Debugf("reading %s", filepath)
				p.readfile(filepath)
			}
		}

	}()

	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	// read files
	var matches []string

	// find files
	for _, currentPath := range p.opt.Path {
		if currentMatches, err := zglob.Glob(currentPath); err == nil {
			// if currentMatches, err := filepath.Glob(currentPath); err == nil {
			matches = append(matches, currentMatches...)
			continue
		}
		return fmt.Errorf("glob(%q) failed", currentPath)
	}

	// ignore excluded
	if len(p.opt.Exclude) > 0 {
		var matches_tmp []string
		for _, pattern := range p.opt.Exclude {
			for _, name := range matches {
				if match, _ := filepath.Match(pattern, name); match == false {
					matches_tmp = append(matches_tmp, name)
				}
			}
		}
		matches = matches_tmp
	}

	// ignore already seen files
	var matches_tmp []string
	for _, name := range matches {
		if !p.sinceDBInfos.has(name) {
			matches_tmp = append(matches_tmp, name)
		}
	}
	matches = matches_tmp

	matches_tmp = []string{}
	for _, name := range matches {
		info, err := os.Stat(name)
		if err != nil {
			p.Logger.Warnf("Error while stating " + name)
			break
		}
		duration := time.Since(info.ModTime()).Seconds()
		// ignore modified to soon
		if duration > float64(p.opt.ReadOlder) {
			// ignore  too old file
			if p.opt.IgnoreOlder > 0 && duration < float64(p.opt.IgnoreOlder) {
			} else {
				matches_tmp = append(matches_tmp, name)
			}
		}
	}
	matches = matches_tmp
	// send to watchChan
	for _, name := range matches {
		p.readfile(name)
	}

	return nil
}

func (p *processor) discoverFilesToRead() error {
	for {
		var matches []string

		// find files
		for _, currentPath := range p.opt.Path {
			if currentMatches, err := zglob.Glob(currentPath); err == nil {
				// if currentMatches, err := filepath.Glob(currentPath); err == nil {
				matches = append(matches, currentMatches...)
				continue
			}
			return fmt.Errorf("glob(%q) failed", currentPath)
		}

		// ignore excluded
		if len(p.opt.Exclude) > 0 {
			var matches_tmp []string
			for _, pattern := range p.opt.Exclude {
				for _, name := range matches {
					if match, _ := filepath.Match(pattern, name); match == false {
						matches_tmp = append(matches_tmp, name)
					}
				}
			}
			matches = matches_tmp
		}

		// ignore already seen files
		var matches_tmp []string
		for _, name := range matches {
			if !p.sinceDBInfos.has(name) {
				matches_tmp = append(matches_tmp, name)
			}
		}
		matches = matches_tmp

		matches_tmp = []string{}
		for _, name := range matches {
			info, err := os.Stat(name)
			if err != nil {
				p.Logger.Warnf("Error while stating " + name)
				break
			}
			duration := time.Since(info.ModTime()).Seconds()
			// ignore modified to soon
			if duration > float64(p.opt.ReadOlder) {
				// ignore  too old file
				if p.opt.IgnoreOlder > 0 && duration < float64(p.opt.IgnoreOlder) {
				} else {
					matches_tmp = append(matches_tmp, name)
				}
			}
		}
		matches = matches_tmp

		// send to watchChan
		for _, name := range matches {
			p.filestoWatch <- name
		}

		select {
		case <-time.NewTicker(time.Second * time.Duration(p.opt.DiscoverInterval)).C:
			p.saveSinceDBInfos()
			continue
		case <-p.q2:
			return nil
		}
	}
	return nil
}

func (p *processor) readfile(pathfile string) error {

	p.Logger.Debug("reading " + pathfile)

	// Create a reader
	f, err := os.Open(pathfile)
	if err != nil {
		p.Logger.Errorf("Error while opening %s : %s", pathfile, err)
		return err
	}
	defer f.Close()

	var dec Decoder

	if dec, err = p.opt.Codec.Decoder(f); err != nil {
		p.Logger.Errorln("decoder error : ", err.Error())
		return err
	}

	for dec.More() {
		if record, err := dec.Decode(); err != nil {
			return err
			break
		} else if record == nil {
			p.Logger.Debugln("waiting for more content...")
			continue
		} else {

			record["file"] = map[string]interface{}{
				"basename": filepath.Base(pathfile),
				"path":     pathfile,
			}
			record["host"] = p.host
			e := p.NewPacket("", record)
			processors.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
			p.Send(e)
		}

		select {
		case <-p.q2:
			return nil // file will not be marked as read :(
		default:
		}
	}

	// mark file read on sincedb
	p.markFileReaded(pathfile)
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	close(p.q2)
	<-p.q
	p.saveSinceDBInfos()
	return nil
}
