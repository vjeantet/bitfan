package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type sinceDBInfo struct {
	Files []string `json:"files,omitempty"`
}

func (s *sinceDBInfo) has(str string) bool {
	for _, v := range s.Files {
		if v == str {
			return true
		}
	}
	return false
}

func (p *processor) markFileReaded(filepath string) {
	if p.opt.SincedbPath == "" || p.opt.SincedbPath == "/dev/null" {
		return
	}
	p.sinceDBInfos.Files = append(p.sinceDBInfos.Files, filepath)
}

func (p *processor) loadSinceDBInfos() (err error) {
	var (
		raw []byte
	)

	p.sinceDBInfos = &sinceDBInfo{}
	p.sinceDBInfosMutex = &sync.Mutex{}

	if p.opt.SincedbPath == "" || p.opt.SincedbPath == "/dev/null" {
		p.Logger.Warn("null sincedb path : since infos will not be persisted")
		return
	}

	if _, err := os.Stat(p.opt.SincedbPath); os.IsNotExist(err) {
		p.Logger.Warnf("sincedb not found: %q", p.opt.SincedbPath)
		return err
	}

	if raw, err = ioutil.ReadFile(p.opt.SincedbPath); err != nil {
		p.Logger.Warnf("Read sincedb failed: %q\n%v", p.opt.SincedbPath, err)
		return
	}

	if err = json.Unmarshal(raw, &p.sinceDBInfos); err != nil {
		p.Logger.Warnf("Unmarshal sincedb failed: %q\n%v", p.opt.SincedbPath, err)
		return
	}

	return
}

func (p *processor) saveSinceDBInfos() (err error) {
	var (
		raw []byte
	)

	p.sinceDBLastSaveTime = time.Now()

	if p.opt.SincedbPath == "" || p.opt.SincedbPath == "/dev/null" {
		return
	}

	p.sinceDBInfosMutex.Lock()
	if raw, err = json.Marshal(p.sinceDBInfos); err != nil {
		p.sinceDBInfosMutex.Unlock()
		p.Logger.Warnf("Marshal sincedb failed: %v", err)
		return
	}
	p.sinceDBInfosMutex.Unlock()

	if err = ioutil.WriteFile(p.opt.SincedbPath, raw, 0664); err != nil {
		p.Logger.Warnf("Write sincedb failed: %q\n%v", p.opt.SincedbPath, err)
		return
	}

	return
}
