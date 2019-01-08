//go:generate go-bindata -o server/assets.go -pkg server -ignore ".DS_Store" -ignore ".scss" assets/...
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/awillis/bitfan/cmd/bitfanUI/commands"
	"github.com/kardianos/service"
)

var version = "master"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Service
	if !service.Interactive() {

		// PASS Service
		s := commands.GetService()

		slogger, err := s.Logger(nil)
		if err != nil {
			log.Fatal("EOOR", err)
		}
		err = s.Run()
		if err != nil {
			slogger.Error("service startup error : ", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	//interactive
	commands.Version = version
	commands.Execute()
}
