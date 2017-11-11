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

package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bitfan",
	Short: "Produce, transform and consume any data",
	Long:  `Bitfan is an open source data processing pipeline`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initSettings(cmd)
		viper.BindPFlag("workers", cmd.Flags().Lookup("workers"))
		viper.BindPFlag("log", cmd.Flags().Lookup("log"))
		viper.BindPFlag("verbose", cmd.Flags().Lookup("verbose"))
		viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
		viper.BindPFlag("data", cmd.Flags().Lookup("data"))
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
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
		fmt.Println("cmd execute error: ", err)
		os.Exit(-1)
	}
}

func initSettings(cmd *cobra.Command) {

	viper.SetConfigName("bitfan") // name of config file (without extension)

	if cmd.Flags().Changed("settings") {
		settings, _ := cmd.Flags().GetString("settings")
		if _, err := os.Stat(settings); err != nil {
			fmt.Printf("settings: %s, error:%v\n", settings, err)
			os.Exit(2)
		}
		viper.AddConfigPath(settings) // optionally look for config in the working directory
		err := viper.ReadInConfig()   // Find and read the config file
		if err != nil {
			fmt.Printf("settings: can not find bitfan.(json|toml|yml) in %s\nerror: %v\n", settings, err)
			os.Exit(2)
		}
	} else {
		viper.AddConfigPath("/etc/bitfan/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.bitfan") // call multiple times to add many search paths
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
	RootCmd.Flags().IntP("workers", "w", runtime.NumCPU(), "number of workers")
	cwd, _ := os.Getwd()
	RootCmd.Flags().String("data", filepath.Join(cwd, ".bitfan"), "Path to data dir")
	RootCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")

	RootCmd.Flags().MarkDeprecated("config", "use the run command")
	RootCmd.Flags().MarkDeprecated("configtest", "use the test command")
	RootCmd.Flags().MarkDeprecated("eval", "use the run or test command")
	RootCmd.Flags().MarkDeprecated("version", "use the version command")
	RootCmd.Flags().MarkHidden("workers")
	RootCmd.Flags().MarkHidden("data")
	RootCmd.Flags().MarkHidden("host")

	RootCmd.PersistentFlags().String("settings", "current dir, then ~/.bitfan/ then /etc/bitfan/", "Set the directory containing the bitfan.toml settings")
	RootCmd.PersistentFlags().StringP("log", "l", "", "Log to a given path. Default is to log to stdout.")
	RootCmd.PersistentFlags().Bool("verbose", false, "Increase verbosity of logs")
	RootCmd.PersistentFlags().Bool("debug", false, "Increase verbosity to the last level (trace)")
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" { // enable ability to specify config file via flag
// 		viper.SetConfigFile(cfgFile)
// 	}

// 	viper.SetConfigName(".bitfan") // name of config file (without extension)
// 	viper.AddConfigPath("$HOME")     // adding home directory as first search path
// 	viper.AutomaticEnv()             // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
