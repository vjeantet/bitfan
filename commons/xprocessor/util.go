package xprocessor

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

func processStdinDataWith(recv func(interface{}) error, maxConcurrent int) *sync.WaitGroup {
	// func processStdinDataWith(p *processor) (*sync.WaitGroup, chan bool) {
	wg := new(sync.WaitGroup)
	ch := make(chan interface{}, maxConcurrent)

	if maxConcurrent > 0 {
		for i := 0; i < maxConcurrent; i++ {
			wg.Add(1)
			go func(ch chan interface{}) {
				for data := range ch {
					if err := recv(data); err != nil {
						Logger.Printf(err.Error())
					}
				}
				wg.Done()
			}(ch)
		}
	} else { // unlimited go routines
		wg.Add(1)
		go func(ch chan interface{}) {
			for data := range ch {
				wg.Add(1)
				go func(data interface{}) {
					if err := recv(data); err != nil {
						Logger.Printf(err.Error())
					}
					wg.Done()
				}(data)
			}
			wg.Done()
		}(ch)
	}

	// done := make(chan bool, 1)
	go func(ch chan interface{}) {
		dec := json.NewDecoder(os.Stdin)
		for {
			var record interface{}
			if err := dec.Decode(&record); err != nil {
				if err == io.EOF {
					// close(done)
				} else {
					Logger.Printf("codec error : %s", err.Error())
				}
				close(ch)
				return
			} else {
				ch <- record
			}
		}
	}(ch)
	return wg //, done
}
