package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/clbanning/mxj"
)

var l *log.Logger

func main() {
	l = log.New(os.Stderr, "", 0)
	maxWorker := flag.Int("concurrency", 1, "a number")

	flag.Parse()

	l.Println("START")

	wg := new(sync.WaitGroup)

	done := processStdin(myEventHandler, *maxWorker, wg)
	myProducer(done)
	wg.Wait()
}

var i int

func myProducer(done chan bool) {
	for {
		i = i + 1
		fields := mxj.Map{}
		fields["count"] = i
		fields["generate"] = true
		time.Sleep(time.Second * 1)
		jsonString, _ := fields.Json()
		os.Stdout.Write(jsonString)

		select {
		case <-done:
			return
		default:
		}
	}

}

func myEventHandler(data interface{}) {
	i = i + 1

	fields := mxj.Map(data.(map[string]interface{}))

	fields["count"] = i
	fields["tick"] = "yes"
	jsonString, err := fields.Json()
	if err != nil {
		l.Printf("Error : %s", err.Error())
		return
	}
	os.Stdout.Write(jsonString)
}

func processStdin(f func(interface{}), maxConcurent int, wg *sync.WaitGroup) chan bool {
	ch := make(chan interface{}, maxConcurent)

	if maxConcurent > 0 {
		for i := 0; i < maxConcurent; i++ {
			wg.Add(1)
			go func(ch chan interface{}) {
				for data := range ch {
					f(data)
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
					f(data)
					wg.Done()
				}(data)
			}
			wg.Done()
		}(ch)
	}

	done := make(chan bool, 1)
	go func(ch chan interface{}) {
		dec := json.NewDecoder(os.Stdin)
		for {
			var record interface{}
			if err := dec.Decode(&record); err != nil {
				if err == io.EOF {

				} else {
					l.Printf("codec error : %s", err.Error())
				}
				close(ch)
				close(done)
				return
			} else {
				ch <- record
			}
		}
	}(ch)
	return done
}
