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
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "logfan",
	Short: "logstash like in go",
	Long:  `Process Any Data, From Any Source`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initSettings(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var configtest = cmd.Flags().Lookup("configtest")
		var eval = cmd.Flags().Lookup("eval")
		var configLocation = cmd.Flags().Lookup("config")
		var version = cmd.Flags().Lookup("version")

		if version.Changed {
			versionCmd.Run(cmd, args)
		} else if configtest.Changed {
			var targs []string
			if configLocation.Changed {
				targs = []string{configLocation.Value.String()}
			} else if eval.Changed {
				targs = []string{eval.Value.String()}
			}
			testCmd.Run(cmd, targs)
		} else if eval.Changed {
			targs := []string{eval.Value.String()}
			runCmd.Run(cmd, targs)
		} else if configLocation.Changed {
			targs := []string{configLocation.Value.String()}
			runCmd.Run(cmd, targs)
		} else {
			cmd.Help()
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

func initSettings(cmd *cobra.Command) {

	viper.SetConfigName("logfan") // name of config file (without extension)

	if cmd.Flags().Changed("settings") {
		settings, _ := cmd.Flags().GetString("settings")
		if _, err := os.Stat(settings); err != nil {
			fmt.Printf("settings: %s, error:%s\n", settings, err)
			os.Exit(2)
		}
		viper.AddConfigPath(settings) // optionally look for config in the working directory
		err := viper.ReadInConfig()   // Find and read the config file
		if err != nil {
			fmt.Printf("settings: can not find logfan.(json|toml|yml) in %s\nerror: %s\n", settings, err)
			os.Exit(2)
		}
	} else {
		viper.AddConfigPath("/etc/logfan/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.logfan") // call multiple times to add many search paths
		viper.AddConfigPath(".")             // optionally look for config in the working directory
		viper.ReadInConfig()                 // Find and read the config file
	}
}

func init() {

	// simulate Logstash flags
	RootCmd.Flags().StringP("config", "f", "", "Load the Logstash config from a file a directory or a url")
	RootCmd.Flags().BoolP("configtest", "t", false, "Test config file or directory")
	RootCmd.Flags().StringP("eval", "e", "", "Use the given string as the configuration data.")
	RootCmd.Flags().BoolP("version", "V", false, "Display version info.")
	// RootCmd.Flags().MarkHidden("configtest")
	// RootCmd.Flags().MarkHidden("eval")
	RootCmd.PersistentFlags().String("settings", "current dir, then ~/.logfan/ then /etc/logfan/", "Set the directory containing the logfan.toml settings")

	RootCmd.PersistentFlags().IntP("filterworkers", "w", runtime.NumCPU(), "number of workers")
	viper.BindPFlag("workers", RootCmd.PersistentFlags().Lookup("filterworkers"))

	RootCmd.PersistentFlags().StringP("log", "l", "", "Log to a given path. Default is to log to stdout.")
	viper.BindPFlag("log", RootCmd.PersistentFlags().Lookup("log"))

	RootCmd.PersistentFlags().Bool("verbose", false, "Increase verbosity to the first level (info), less verbose.")
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))

	RootCmd.PersistentFlags().Bool("debug", false, "Increase verbosity to the last level (trace), more verbose.")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" { // enable ability to specify config file via flag
// 		viper.SetConfigFile(cfgFile)
// 	}

// 	viper.SetConfigName(".logfan") // name of config file (without extension)
// 	viper.AddConfigPath("$HOME")     // adding home directory as first search path
// 	viper.AutomaticEnv()             // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
