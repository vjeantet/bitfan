package tail

// // Code source from github.com/tsaikd/gogstash
// // @see https://github.com/tsaikd/gogstash/tree/master/input/file

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io/ioutil"
// 	"os"
// 	"sync"
// 	"time"
// )

// type sinceDBInfo struct {
// 	Offset int64 `json:"offset,omitempty"`
// }

// func (p *processor) loadSinceDBInfos() (err error) {
// 	var (
// 		raw []byte
// 	)
// 	p.sinceDBInfosMutex = &sync.Mutex{}

// 	p.sinceDBInfos = map[string]*sinceDBInfo{}

// 	if p.opt.SincedbPath == "" || p.opt.SincedbPath == "/dev/null" {
// 		p.Logger.Debugf("No valid sincedb path : %s", p.opt.SincedbPath)
// 		return
// 	}

// 	if _, err := os.Stat(p.opt.SincedbPath); os.IsNotExist(err) {
// 		p.Logger.Warnf("sincedb not found: %q", p.opt.SincedbPath)
// 		return err
// 	}

// 	if raw, err = ioutil.ReadFile(p.opt.SincedbPath); err != nil {
// 		p.Logger.Warnf("Read sincedb failed: %q\n%v", p.opt.SincedbPath, err)
// 		return
// 	}

// 	if err = json.Unmarshal(raw, &p.sinceDBInfos); err != nil {
// 		p.Logger.Warnf("Unmarshal sincedb failed: %q\n%v", p.opt.SincedbPath, err)
// 		return
// 	}

// 	return
// }

// func (p *processor) saveSinceDBInfos() (err error) {
// 	var (
// 		raw []byte
// 	)

// 	p.sinceDBLastSaveTime = time.Now()

// 	if p.opt.SincedbPath == "" || p.opt.SincedbPath == "/dev/null" {
// 		p.Logger.Debugf("No valid sincedb path : %s", p.opt.SincedbPath)
// 		return
// 	}

// 	p.sinceDBInfosMutex.Lock()
// 	if raw, err = json.Marshal(p.sinceDBInfos); err != nil {
// 		p.sinceDBInfosMutex.Unlock()
// 		p.Logger.Warnf("Marshal sincedb failed: %v", err)
// 		return
// 	}
// 	p.sinceDBInfosMutex.Unlock()

// 	p.sinceDBLastInfosRaw = raw

// 	if err = ioutil.WriteFile(p.opt.SincedbPath, raw, 0664); err != nil {
// 		p.Logger.Warnf("Write sincedb failed: %q\n%v", p.opt.SincedbPath, err)
// 		return
// 	}

// 	return
// }

// func (p *processor) checkSaveSinceDBInfos() (err error) {
// 	var (
// 		raw []byte
// 	)
// 	if time.Since(p.sinceDBLastSaveTime) > time.Duration(p.opt.SincedbWriteInterval)*time.Second {
// 		if raw, err = json.Marshal(p.sinceDBInfos); err != nil {
// 			p.Logger.Warnf("Marshal sincedb failed: %v", err)
// 			return
// 		}
// 		if !bytes.Equal(raw, p.sinceDBLastInfosRaw) {
// 			err = p.saveSinceDBInfos()
// 		}
// 	}
// 	return
// }

// func (p *processor) checkSaveSinceDBInfosLoop() (err error) {
// 	for {
// 		time.Sleep(time.Duration(p.opt.SincedbWriteInterval) * time.Second)
// 		if err = p.checkSaveSinceDBInfos(); err != nil {
// 			return
// 		}
// 	}
// }
