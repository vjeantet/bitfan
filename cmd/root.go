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
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var configPath, configString, logPath string
var filterworkers int
var verbose, debug, version, configtest bool

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "logstack",
	Short: "a logstash fork in go",
	Long: `Logstack is a logstash fork.

Process Any Data, From Any Source
Centralize data processing of all types
Normalize varying schema and formats
Quickly extend to custom log formats`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			printVersion()
		} else if configtest {
			testConfig(configPath, configString, args)
		} else if configString != "" {
			startLogstack("", configString, args)
		} else {
			startLogstack(configPath, "", args)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.Flags().StringVarP(&configPath, "config", "f", "", "Load the Logstash config from a file or directory")
	RootCmd.Flags().StringVarP(&configString, "eval", "e", "", "Use the given string as the configuration data.")
	RootCmd.Flags().IntVarP(&filterworkers, "filterworkers", "w", runtime.NumCPU(), "number of workers")

	RootCmd.Flags().StringVarP(&logPath, "log", "l", "", "Log to a given path. Default is to log to stdout.")
	RootCmd.Flags().BoolVarP(&verbose, "verbose", "", false, "Increase verbosity to the first level (info), less verbose.")
	RootCmd.Flags().BoolVarP(&debug, "debug", "", false, "Increase verbosity to the last level (trace), more verbose.")

	RootCmd.Flags().BoolVarP(&version, "version", "V", false, "Display the version of Logstash.")
	RootCmd.Flags().BoolVarP(&configtest, "configtest", "t", false, "Test config file or directory")

}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" { // enable ability to specify config file via flag
// 		viper.SetConfigFile(cfgFile)
// 	}

// 	viper.SetConfigName(".logstack") // name of config file (without extension)
// 	viper.AddConfigPath("$HOME")     // adding home directory as first search path
// 	viper.AutomaticEnv()             // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
