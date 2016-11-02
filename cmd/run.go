// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/veino/logfan/lib"
	"github.com/veino/veino/config"
	"github.com/veino/veino/runtime"
	"github.com/veino/veino/runtime/metrics"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [config1] [config2]",
	Short: "Run logfan",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		var cwd string

		cwd, _ = os.Getwd()

		var locs locations

		// Si il n'y a pas argument alors utilise viper.config
		if len(args) == 0 {
			cwd = filepath.Dir(viper.ConfigFileUsed())
			for _, conf := range viper.GetStringSlice("config") {
				if filepath.IsAbs(conf) {
					locs.add(conf, cwd)
				} else {
					locs.add(filepath.Join(cwd, conf), cwd)
				}
			}
		} else {
			for _, v := range args {
				locs.add(v, cwd)
			}
		}

		// Si location est un dossier
		//   calcul les autres locations
		var configAgents = []config.Agent{}
		var locConfigAgents = []config.Agent{}

		for _, loc := range locs.items {
			if loc.kind == "url" {
				content, ncwl, err := lib.GetContentFromLocation(loc.path, loc.workingpath)
				if err != nil {
					fmt.Printf("error %s\n", err)
				}
				uriSegments := strings.Split(loc.path, "/")
				pipelineName := strings.Join(uriSegments[2:], ".")
				locConfigAgents, err = lib.ParseConfig(pipelineName, content, ncwl)
				if err != nil {
					log.Fatalf("error %s", err.Error())
				}
				log.Printf("using config url : %s", loc.path)
				configAgents = append(configAgents, locConfigAgents...)
			}

			if loc.kind == "file" {
				locsexpanded, err := loc.expand()
				if err != nil {
					log.Fatalln(err)
				}

				for _, file := range locsexpanded {
					content, ncwl, err := lib.GetContentFromLocation(file, loc.workingpath)
					if err != nil {
						fmt.Printf("error %s\n", err)
					}
					filename := filepath.Base(file)
					extension := filepath.Ext(filename)
					pipelineName := filename[0 : len(filename)-len(extension)]

					locConfigAgents, err = lib.ParseConfig(pipelineName, content, ncwl)
					if err != nil {
						break
					}
					log.Printf("using config file : %s", file)
					configAgents = append(configAgents, locConfigAgents...)
				}
			}
		}

		// pp.Println("configAgents-->", configAgents)
		var stats metrics.IStats
		if true == viper.IsSet("prometheus") {
			stats = metrics.NewPrometheus(viper.GetString("prometheus.listen"), viper.GetString("prometheus.path"))
		} else {
			stats = &metrics.StatsVoid{}
		}
		runtime.SetIStat(stats)
		runtime.Logger().SetVerboseMode(viper.GetBool("verbose"))
		runtime.Logger().SetDebugMode(viper.GetBool("debug"))
		if viper.IsSet("log") {
			runtime.Logger().SetOutputFile(viper.GetString("log"))
		}
		log := runtime.Logger()

		runtime.Start(viper.GetString("webhook.listen"))
		err := runtime.StartAgents(configAgents)

		if err != nil {
			log.Fatalln(err)
		}

		log.Println("logfan ready")

		if service.Interactive() {
			// Wait for signal CTRL+C for send a stop event to all AgentProcessor
			// When CTRL+C, SIGINT and SIGTERM signal occurs
			// Then stop server gracefully
			ch := make(chan os.Signal)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			<-ch
			close(ch)

			fmt.Println("")
			log.Printf("LogFan is stopping...")
			runtime.Stop()
			log.Printf("Everything stopped gracefully. Goodbye!\n")
		}

	},
}

func init() {
	runCmd.Flags().Bool("prometheus", false, "Export stats using prometheus output")
	viper.BindPFlag("prometheus", runCmd.Flags().Lookup("prometheus"))

	runCmd.Flags().String("prometheus.listen", "0.0.0.0:24232", "Address and port to bind Prometheus metrics")
	viper.BindPFlag("prometheus.listen", runCmd.Flags().Lookup("prometheus.listen"))

	runCmd.Flags().String("prometheus.path", "/metrics", "Expose Prometheus metrics at specified path.")
	viper.BindPFlag("prometheus.path", runCmd.Flags().Lookup("prometheus.path"))

	runCmd.Flags().String("webhook.listen", "127.0.0.1:19090", "Address and port to bind webhooks")
	viper.BindPFlag("webhook.listen", runCmd.Flags().Lookup("webhook.listen"))

	RootCmd.AddCommand(runCmd)
}

func startLogfan_(
	flagConfigPath string,
	flagConfigContent string,
	stats metrics.IStats,
	webhookListen string,
	verbose bool,
	debug bool,
	logPath string,
	args []string) error {

	runtime.SetIStat(stats)
	runtime.Start(webhookListen)

	if logPath != "" {
		runtime.Logger().SetOutputFile(logPath)
	}

	runtime.Logger().SetVerboseMode(verbose)
	runtime.Logger().SetDebugMode(debug)
	log := runtime.Logger()

	var configAgents = []config.Agent{}

	// Load agents from flagConfigContent string
	if flagConfigContent != "" {
		pwd, _ := os.Getwd()
		fileConfigAgents, err := lib.ParseConfig("inline", []byte(flagConfigContent), pwd)
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
			configAgents, err = lib.ParseConfigLocation(pipelineName, fileLocation, baseUrl)
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
					fileConfigAgents, err = lib.ParseConfigLocation(pipelineName, fileLocation, filepath.Dir(file))
					if err != nil {
						break
					}
					log.Printf("using config file : %s", file)

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

	log.Println("logfan ready")
	if service.Interactive() {
		// Wait for signal CTRL+C for send a stop event to all AgentProcessor
		// When CTRL+C, SIGINT and SIGTERM signal occurs
		// Then stop server gracefully
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		close(ch)

		fmt.Println("")
		log.Printf("LogFan is stopping...")
		runtime.Stop()
		log.Printf("Everything stopped gracefully. Goodbye!\n")

	}
	return err
}
