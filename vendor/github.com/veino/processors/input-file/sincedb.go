package fileinput

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/tsaikd/KDGoLib/futil"
)

type SinceDBInfo struct {
	Offset int64 `json:"offset,omitempty"`
}

func (p *processor) LoadSinceDBInfos() (err error) {
	var (
		raw []byte
	)
	p.Logger.Println("LoadSinceDBInfos")
	p.SinceDBInfos = map[string]*SinceDBInfo{}

	if p.opt.Sincedb_path == "" || p.opt.Sincedb_path == "/dev/null" {
		p.Logger.Println("No valid sincedb path")
		return
	}

	if !futil.IsExist(p.opt.Sincedb_path) {
		p.Logger.Printf("sincedb not found: %q", p.opt.Sincedb_path)
		return
	}

	if raw, err = ioutil.ReadFile(p.opt.Sincedb_path); err != nil {
		p.Logger.Printf("Read sincedb failed: %q\n%s", p.opt.Sincedb_path, err)
		return
	}

	if err = json.Unmarshal(raw, &p.SinceDBInfos); err != nil {
		p.Logger.Printf("Unmarshal sincedb failed: %q\n%s", p.opt.Sincedb_path, err)
		return
	}

	return
}

func (p *processor) SaveSinceDBInfos() (err error) {
	var (
		raw []byte
	)
	p.Logger.Println("SaveSinceDBInfos")
	p.SinceDBLastSaveTime = time.Now()

	if p.opt.Sincedb_path == "" || p.opt.Sincedb_path == "/dev/null" {
		p.Logger.Println("No valid sincedb path")
		return
	}

	if raw, err = json.Marshal(p.SinceDBInfos); err != nil {
		p.Logger.Printf("Marshal sincedb failed: %s", err)
		return
	}
	p.sinceDBLastInfosRaw = raw

	if err = ioutil.WriteFile(p.opt.Sincedb_path, raw, 0664); err != nil {
		p.Logger.Printf("Write sincedb failed: %q\n%s", p.opt.Sincedb_path, err)
		return
	}

	return
}

func (p *processor) CheckSaveSinceDBInfos() (err error) {
	var (
		raw []byte
	)
	if time.Since(p.SinceDBLastSaveTime) > time.Duration(p.opt.Sincedb_write_interval)*time.Second {
		if raw, err = json.Marshal(p.SinceDBInfos); err != nil {
			p.Logger.Printf("Marshal sincedb failed: %s", err)
			return
		}
		if bytes.Compare(raw, p.sinceDBLastInfosRaw) != 0 {
			err = p.SaveSinceDBInfos()
		}
	}
	return
}

func (p *processor) CheckSaveSinceDBInfosLoop() (err error) {
	for {
		time.Sleep(time.Duration(p.opt.Sincedb_write_interval) * time.Second)
		if err = p.CheckSaveSinceDBInfos(); err != nil {
			return
		}
	}
	return
}
