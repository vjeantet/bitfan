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
	"os"
	"os/signal"
	"syscall"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/veino/logfan/lib"
	"github.com/veino/veino/runtime"
	"github.com/veino/veino/runtime/metrics"
)

func init() {
	RootCmd.AddCommand(runCmd)
	initRunFlags(runCmd)
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [config1] [config2] [config...]",
	Short: "Run logfan",
	Long: `Load and run pipelines configured in configuration files (logstash format)
you can set multiples files, urls, diretories, or a configuration content as a string (mimic the logstash -e flag)

When no configuration is passed to the command, logfan use the config set in global settings file logfan.(toml|yml|json)
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initRunConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runtime.Logger().SetVerboseMode(viper.GetBool("verbose"))
		runtime.Logger().SetDebugMode(viper.GetBool("debug"))
		if viper.IsSet("log") {
			runtime.Logger().SetOutputFile(viper.GetString("log"))
		}
		log := runtime.Logger()

		var locations lib.Locations
		cwd, _ := os.Getwd()

		if len(args) == 0 {
			for _, v := range viper.GetStringSlice("config") {
				locations.Add(v, cwd)
			}
		} else {
			for _, v := range args {
				locations.Add(v, cwd)
			}
		}

		var stats metrics.IStats
		if true == viper.IsSet("prometheus") {
			stats = metrics.NewPrometheus(viper.GetString("prometheus.listen"), viper.GetString("prometheus.path"))
		} else {
			stats = &metrics.StatsVoid{}
		}
		runtime.SetIStat(stats)

		if !viper.GetBool("no-network") {
			runtime.Start("")
		} else {
			runtime.Start(viper.GetString("webhook.listen"))
		}

		for _, loc := range locations.Items {
			agt, err := loc.ConfigAgents()

			if err != nil {
				log.Printf("Error : %s %s", loc.Path, err)
				os.Exit(2)
			}
			ppl := loc.ConfigPipeline()

			// Allow pipeline customisation only when only one location was provided by user
			if len(locations.Items) == 1 {
				if cmd.Flags().Changed("name") {
					ppl.Name, _ = cmd.Flags().GetString("name")
				}
				if cmd.Flags().Changed("id") {
					ppl.ID, _ = cmd.Flags().GetString("id")
				}
			}

			_, err = runtime.StartPipeline(&ppl, agt)
			if err != nil {
				log.Printf("error : %s\n", err.Error())
				os.Exit(1)
			}
		}

		if viper.GetBool("no-network") {
			log.Println("Veino API disabled")
		} else {
			lib.ApiServe(viper.GetString("host"))
			log.Println("Veino API listening on", viper.GetString("host"))
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

			log.Println("")
			log.Printf("LogFan is stopping...")
			runtime.Stop()
			log.Printf("Everything stopped gracefully. Goodbye!\n")
		}

	},
}

func initRunConfig(cmd *cobra.Command) {
	viper.BindPFlag("prometheus", cmd.Flags().Lookup("prometheus"))
	viper.BindPFlag("prometheus.listen", cmd.Flags().Lookup("prometheus.listen"))
	viper.BindPFlag("prometheus.path", cmd.Flags().Lookup("prometheus.path"))
	viper.BindPFlag("webhook.listen", cmd.Flags().Lookup("webhook.listen"))
	viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	viper.BindPFlag("no-network", cmd.Flags().Lookup("no-network"))
}

func initRunFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("prometheus", false, "Export stats using prometheus output")
	cmd.Flags().String("prometheus.listen", "0.0.0.0:24232", "Address and port to bind Prometheus metrics")
	cmd.Flags().String("prometheus.path", "/metrics", "Expose Prometheus metrics at specified path.")
	cmd.Flags().String("webhook.listen", "127.0.0.1:19090", "Address and port to bind webhooks")
	cmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
	cmd.Flags().Bool("no-network", false, "Disable network (api and webhook)")
	cmd.Flags().String("name", "", "set pipeline's name")
	cmd.Flags().String("id", "", "set pipeline's id")
}
