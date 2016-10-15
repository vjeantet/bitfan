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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/veino/veino/config"
	"github.com/veino/veino/runtime"
)

var flagTestConfigPath, flagTestConfigContent string

func testConfig(flagConfigPath string, flagConfigContent string, args []string) {
	// runtime.Start()
	logError := log.New(os.Stderr, "ERROR: ", 0)
	logInfo := log.New(os.Stdout, "", 0)
	logWarning := log.New(os.Stdout, "WARNING: ", 0)

	countConfig := 0

	// Load agents from flagConfigContent string
	if flagConfigContent != "" {
		configAgents, err := parseConfig("inline", []byte(flagConfigContent))
		if err != nil {
			logError.Fatalf("%s", err.Error())
		}

		configAgentsOrdered := config.Sort(configAgents, config.SortOutputsFirst)
		for _, configAgent := range configAgentsOrdered {
			_, err := runtime.NewAgent(configAgent, 0)
			if err != nil {
				logError.Fatalf("plugin '%s' may not start : %s", configAgent.Type, err.Error())
			}
		}
		countConfig++
	}

	// Load all agents configuration from conf files
	if flagConfigPath != "" {

		if fi, err := os.Stat(flagConfigPath); err == nil {
			if fi.IsDir() {
				flagConfigPath = flagConfigPath + "/*.*"
			}
		} else {
			logError.Fatalf("error %s", err.Error())
		}

		//List all conf files if flagConfigPath folder
		files, err := filepath.Glob(flagConfigPath)
		if err != nil {
			logError.Fatalf("error %s", err.Error())
		}

		//use each file
		for _, file := range files {

			content, err := ioutil.ReadFile(file)
			if err != nil {
				logWarning.Printf(`Error while reading "%s" [%s]`, file, err)
				continue
			}

			// instance all AgenConfiguration structs from file content
			switch strings.ToLower(filepath.Ext(file)) {
			case ".conf":
				var filename = filepath.Base(file)
				var extension = filepath.Ext(filename)
				var pipelineName = filename[0 : len(filename)-len(extension)]

				configAgents, err := parseConfig(pipelineName, content)
				if err != nil {
					logError.Fatalf("%s", err.Error())
				}
				logInfo.Printf("checking %s", file)
				configAgentsOrdered := config.Sort(configAgents, config.SortOutputsFirst)
				for _, configAgent := range configAgentsOrdered {
					_, err := runtime.NewAgent(configAgent, 0)
					if err != nil {
						logError.Fatalf("plugin '%s' may not start : %s", configAgent.Type, err.Error())
					}
					countConfig++
				}

			default:
				logInfo.Printf("ignored file %s", file)
			}

			if err != nil {
				logError.Fatalf("%s", err.Error())
			}

		}
	}
	if countConfig > 0 {
		logInfo.Printf("Everything is ok\n")
	} else {
		logError.Fatalf("No configuration found")
	}
}
