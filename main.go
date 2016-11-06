//go:generate go generate github.com/veino/veino/processors/...
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

package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/viper"
	"github.com/veino/logfan/cmd"
)

var version = "dev"
var buildstamp = ""

func main() {
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.AddConfigPath("/etc/logfan/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.logfan") // call multiple times to add many search paths
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	viper.ReadInConfig()                 // Find and read the config file
	// if err != nil {                      // Handle errors reading the config file
	// 	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	// }

	// Service
	if !service.Interactive() {
		s := cmd.GetService()

		slogger, err := s.Logger(nil)
		if err != nil {
			log.Fatal(err)
		}
		err = s.Run()
		if err != nil {
			slogger.Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	//interactive
	cmd.Version = version
	cmd.Buildstamp = buildstamp
	cmd.Execute()
}
