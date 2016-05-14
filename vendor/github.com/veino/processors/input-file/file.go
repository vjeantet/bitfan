package fileinput

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-fsnotify/fsnotify"
	"github.com/veino/veino"
)

func (p *processor) Start2(e veino.IPacket) (err error) {
	defer func() {
		if err != nil {
			p.Logger.Println(err)
		}
	}()

	var (
		matches []string
		fi      os.FileInfo
	)

	if err = p.LoadSinceDBInfos(); err != nil {
		return
	}
	current_path := p.opt.Path[0]
	if matches, err = filepath.Glob(current_path); err != nil {
		return fmt.Errorf("glob(%q) failed", current_path)
	}

	go p.CheckSaveSinceDBInfosLoop()

	for _, fpath := range matches {
		if fpath, err = filepath.EvalSymlinks(fpath); err != nil {
			p.Logger.Printf("Get symlinks failed: %q\n%v", fpath, err)
			continue
		}

		if fi, err = os.Stat(fpath); err != nil {
			p.Logger.Printf("stat(%q) failed\n%s", current_path, err)
			continue
		}

		if fi.IsDir() {
			p.Logger.Printf("Skipping directory: %q", current_path)
			continue
		}

		readEventChan := make(chan fsnotify.Event, 10)
		go p.fileReadLoop(readEventChan, fpath)
		go p.fileWatchLoop(readEventChan, fpath, fsnotify.Create|fsnotify.Write)
	}

	return
}

func (p *processor) fileReadLoop(
	readEventChan chan fsnotify.Event,
	fpath string,
) (err error) {
	var (
		since     *SinceDBInfo
		fp        *os.File
		truncated bool
		ok        bool
		whence    int
		reader    *bufio.Reader
		line      string
		size      int

		buffer = &bytes.Buffer{}
	)

	if fpath, err = filepath.EvalSymlinks(fpath); err != nil {
		p.Logger.Println("Get symlinks failed: %q\n%v", fpath, err)
		return
	}

	if since, ok = p.SinceDBInfos[fpath]; !ok {
		p.SinceDBInfos[fpath] = &SinceDBInfo{}
		since = p.SinceDBInfos[fpath]
	}

	if since.Offset == 0 {
		if p.opt.Start_position == "end" {
			whence = os.SEEK_END
		} else {
			whence = os.SEEK_SET
		}
	} else {
		whence = os.SEEK_SET
	}

	if fp, reader, err = openfile(fpath, since.Offset, whence); err != nil {
		return
	}
	defer fp.Close()

	if truncated, err = isFileTruncated(fp, since); err != nil {
		return
	}
	if truncated {
		p.Logger.Printf("File truncated, seeking to beginning: %q", fpath)
		since.Offset = 0
		if _, err = fp.Seek(since.Offset, os.SEEK_SET); err != nil {
			p.Logger.Println("seek file failed: %q", fpath)
			return
		}
	}

	for {
		if line, size, err = readline(reader, buffer); err != nil {
			if err == io.EOF {
				watchev := <-readEventChan
				p.Logger.Println("fileReadLoop recv:", watchev)
				if watchev.Op&fsnotify.Create == fsnotify.Create {
					p.Logger.Printf("File recreated, seeking to beginning: %q", fpath)
					fp.Close()
					since.Offset = 0
					if fp, reader, err = openfile(fpath, since.Offset, os.SEEK_SET); err != nil {
						return
					}
				}
				if truncated, err = isFileTruncated(fp, since); err != nil {
					return
				}
				if truncated {
					p.Logger.Printf("File truncated, seeking to beginning: %q", fpath)
					since.Offset = 0
					if _, err = fp.Seek(since.Offset, os.SEEK_SET); err != nil {
						p.Logger.Println("seek file failed: %q", fpath)
						return
					}
					continue
				}
				p.Logger.Println("watch %q %q %v", watchev.Name, fpath, watchev)
				continue
			} else {
				return
			}
		}

		host, err := os.Hostname()
		if err != nil {
			p.Logger.Printf("can not get hostname : %s", err.Error())
		}

		ne := p.NewPacket(line, map[string]interface{}{
			"host":   host,
			"path":   fpath,
			"offset": since.Offset,
		})

		since.Offset += int64(size)

		p.Logger.Printf("%q %v", ne.Message(), ne)
		// field.ProcessCommonFields(ne.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
		p.Send(ne)

		//self.SaveSinceDBInfos()
		p.CheckSaveSinceDBInfos()
	}

	return
}

func (self *processor) fileWatchLoop(readEventChan chan fsnotify.Event, fpath string, op fsnotify.Op) (err error) {
	var (
		event fsnotify.Event
	)
	for {
		if event, err = waitWatchEvent(fpath, op); err != nil {
			return
		}
		readEventChan <- event
	}
	return
}

func isFileTruncated(fp *os.File, since *SinceDBInfo) (truncated bool, err error) {
	var (
		fi os.FileInfo
	)
	if fi, err = fp.Stat(); err != nil {
		err = fmt.Errorf("stat file failed: "+fp.Name(), err)
		return
	}
	if fi.Size() < since.Offset {
		truncated = true
	} else {
		truncated = false
	}
	return
}

func openfile(fpath string, offset int64, whence int) (fp *os.File, reader *bufio.Reader, err error) {
	if fp, err = os.Open(fpath); err != nil {
		err = fmt.Errorf("open file failed: "+fpath, err)
		return
	}

	if _, err = fp.Seek(offset, whence); err != nil {
		err = fmt.Errorf("seek file failed: " + fpath)
		return
	}

	reader = bufio.NewReaderSize(fp, 16*1024)
	return
}

func readline(reader *bufio.Reader, buffer *bytes.Buffer) (line string, size int, err error) {
	var (
		segment []byte
	)

	for {
		if segment, err = reader.ReadBytes('\n'); err != nil {
			if err != io.EOF {
				err = fmt.Errorf("read line failed", err)
			}
			return
		}

		if _, err = buffer.Write(segment); err != nil {
			err = fmt.Errorf("write buffer failed", err)
			return
		}

		if isPartialLine(segment) {
			time.Sleep(1 * time.Second)
		} else {
			size = buffer.Len()
			line = buffer.String()
			buffer.Reset()
			line = strings.TrimRight(line, "\r\n")
			return
		}
	}

	return
}

func isPartialLine(segment []byte) bool {
	if len(segment) < 1 {
		return true
	}
	if segment[len(segment)-1] != '\n' {
		return true
	}
	return false
}

var (
	mapWatcher = map[string]*fsnotify.Watcher{}
)

func waitWatchEvent(fpath string, op fsnotify.Op) (event fsnotify.Event, err error) {
	var (
		fdir    string
		watcher *fsnotify.Watcher
		ok      bool
	)

	if fpath, err = filepath.EvalSymlinks(fpath); err != nil {
		err = fmt.Errorf("Get symlinks failed: "+fpath, err)
		return
	}

	fdir = filepath.Dir(fpath)

	if watcher, ok = mapWatcher[fdir]; !ok {
		//		p.Logger.Debugf("create new watcher for %q", fdir)
		if watcher, err = fsnotify.NewWatcher(); err != nil {
			err = fmt.Errorf("create new watcher failed: "+fdir, err)
			return
		}
		mapWatcher[fdir] = watcher
		if err = watcher.Add(fdir); err != nil {
			err = fmt.Errorf("add new watch path failed: "+fdir, err)
			return
		}
	}

	for {
		select {
		case event = <-watcher.Events:
			if event.Name == fpath {
				if op > 0 {
					if event.Op&op > 0 {
						return
					}
				} else {
					return
				}
			}
		case err = <-watcher.Errors:
			err = fmt.Errorf("watcher error", err)
			return
		}
	}

	return
}
