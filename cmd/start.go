// Copyright © 2016 Valere JEANTET <valere.jeantet@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/kardianos/service"
	"github.com/veino/veino/config"
	"github.com/veino/veino/runtime/metrics"

	"github.com/veino/veino/runtime"
)

var flagConfigPath string

func startLogfan(flagConfigPath string, flagConfigContent string, stats metrics.IStats, args []string) error {

	runtime.SetIStat(stats)
	runtime.Start(webhookListen)

	runtime.Logger().SetVerboseMode(verbose)
	runtime.Logger().SetDebugMode(debug)

	var configAgents = []config.Agent{}

	// Load agents from flagConfigContent string
	if flagConfigContent != "" {
		pwd, _ := os.Getwd()
		fileConfigAgents, err := parseConfig("inline", []byte(flagConfigContent), pwd)
		if err != nil {
			log.Fatalln("ERROR while using config ", err.Error())
		}
		configAgents = append(configAgents, fileConfigAgents...)
	}

	// Load all agents configuration from conf files
	if flagConfigPath != "" {

		if v, _ := url.Parse(flagConfigPath); v.Scheme == "http" || v.Scheme == "https" {
			uriSegments := strings.Split(flagConfigPath, "/")
			var baseUrl = strings.Join(uriSegments[:len(uriSegments)-1], "/") + "/"
			var pipelineName = strings.Join(uriSegments[2:], ".")
			fileLocation := map[string]interface{}{"url": flagConfigPath}
			var err error
			configAgents, err = parseConfigLocation(pipelineName, fileLocation, baseUrl)
			if err != nil {
				log.Fatalf("error %s", err.Error())
			}
		} else {

			if fi, err := os.Stat(flagConfigPath); err == nil {
				if fi.IsDir() {
					flagConfigPath = flagConfigPath + "/*.*"
				}
			} else {
				log.Fatalf("error %s", err.Error())
			}

			//List all conf files if flagConfigPath folder
			files, err := filepath.Glob(flagConfigPath)
			if err != nil {
				log.Fatalf("error %s", err.Error())
			}

			//use each file
			for _, file := range files {
				var fileConfigAgents []config.Agent

				// instance all AgenConfiguration structs from file content
				switch strings.ToLower(filepath.Ext(file)) {
				case ".conf":
					var filename = filepath.Base(file)
					var extension = filepath.Ext(filename)
					var pipelineName = filename[0 : len(filename)-len(extension)]

					fileLocation := map[string]interface{}{"path": file}
					fileConfigAgents, err = parseConfigLocation(pipelineName, fileLocation, filepath.Dir(file))
					if err != nil {
						break
					}
					log.Printf("using config file : %s\n", file)

				default:
					log.Printf("ignored file %s", file)
				}

				if err != nil {
					log.Fatalf("error %s", err.Error())
				}

				configAgents = append(configAgents, fileConfigAgents...)
			}
		}
	}

	err := runtime.StartAgents(configAgents)

	runtime.Logger().Infoln("logfan ready !")
	if service.Interactive() {
		// Wait for signal CTRL+C for send a stop event to all AgentProcessor
		// When CTRL+C, SIGINT and SIGTERM signal occurs
		// Then stop server gracefully
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		close(ch)

		fmt.Println("")
		log.Printf("stopping...")
		runtime.Stop()
		log.Printf("Everything stopped gracefully. Goodbye!\n")

	}
	return err
}
