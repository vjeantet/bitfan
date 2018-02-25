// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element.

// used as processor
// # Inputs :
// - html content (field name)
// - ids of elements to image
// - StripCSSANDJSLink

// # output :
// same event with html content replaced with html with embeded images or local images

// SingleHTMLStatic{
//   fieldName => "content"
//   FixIds =>
//   EmbedasB64 =>
//   StripLink => css, js
// }

package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/clbanning/mxj"
)

var l *log.Logger
var fieldName string
var embedImage bool
var fixIds []string

// {"array":[1,2,3],"boolean":true,"null":null,"number":123,"object":{"a":"b","c":"d","e":"f"},"output":"<div id=\"oo\" style=\"color:red\">test</div><div id=\"ooo\" style=\"color:red\">sdjflkqs jdfmlkqsj dflkmqj sdflmqksjdf</div>"}
func getEventHandler(cs *ChromeShot) func(interface{}) {
	return func(data interface{}) {
		fields := mxj.Map(data.(map[string]interface{}))

		htmlContent := fields.ValueOrEmptyForPathString(fieldName)
		// l.Println("htmlContent-->", htmlContent)

		var err error
		htmlContent, err = cs.EmbedImageForDomElements(htmlContent, fixIds)
		if err != nil {
			l.Printf("Error : %s", err.Error())
			return
		}

		// Return resulting HTML
		fields.SetValueForPath(htmlContent, fieldName)
		jsonBytes, err := fields.Json()
		if err != nil {
			l.Printf("Error : %s", err.Error())
			return
		}
		os.Stdout.Write(jsonBytes)
	}
}

func main() {
	l = log.New(os.Stderr, "", 0)

	flag.StringVar(&fieldName, "fieldName", "output", "field name with html")
	flag.BoolVar(&embedImage, "embedImage", true, "embed image to html")
	fixIdsTmp := flag.String("ids", "", "id of element to fix")
	flag.Parse()

	l.Println("START")

	fixIds = strings.Split(*fixIdsTmp, ",")

	l.Println("fieldName-->", fieldName)
	l.Println("embedImage-->", embedImage)

	wg := new(sync.WaitGroup)

	cs, err := NewChromeShot()
	if err != nil {
		l.Println("error ", err.Error())
		os.Exit(2)
	}

	processStdin(getEventHandler(cs), 1, wg)
	wg.Wait()

	cs.Stop()
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
					close(done)
				} else {
					l.Printf("codec error : %s", err.Error())
				}
				close(ch)
				return
			} else {
				ch <- record
			}
		}
	}(ch)
	return done
}
