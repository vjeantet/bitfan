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

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [config1] [config2] [config...]",
	Short: "Run logfan",
	Long: `Load and run pipelines configured in configuration files (logstash format)
you can set multiples files, urls, diretories, or a configuration content as a string (mimic the logstash -e flag)

When no configuration is passed to the command, logfan use the config set in global settings file logfan.(toml|yml|json)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var locations lib.Locations
		cwd, _ := os.Getwd()
		for _, v := range args {
			locations.Add(v, cwd)
		}

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

		for _, loc := range locations.Items {
			agt, err := loc.ConfigAgents()

			if err != nil {
				fmt.Printf("Error : %s %s", loc.Path, err)
				os.Exit(2)
			}
			ppl := loc.ConfigPipeline()

			_, err = runtime.StartPipeline(&ppl, agt)
			if err != nil {
				fmt.Printf("error : %s\n", err.Error())
				os.Exit(1)
			}
		}

		lib.ApiServe(viper.GetString("host"))
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

	runCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
	viper.BindPFlag("host", runCmd.Flags().Lookup("host"))

	RootCmd.AddCommand(runCmd)
}
