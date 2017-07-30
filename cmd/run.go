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
	"path/filepath"
	"syscall"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vjeantet/bitfan/api"
	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/lib"
)

func init() {
	RootCmd.AddCommand(runCmd)
	initRunFlags(runCmd)
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [config1] [config2] [config...]",
	Short: "Run bitfan",
	Long: `Load and run pipelines configured in configuration files (logstash format)
you can set multiples files, urls, diretories, or a configuration content as a string (mimic the logstash -e flag)

When no configuration is passed to the command, bitfan use the config set in global settings file bitfan.(toml|yml|json)
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initRunConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {

		core.SetLogVerboseMode(viper.GetStringSlice("verbose"))
		core.SetLogDebugMode(viper.GetStringSlice("debug"))

		if viper.IsSet("log") {
			core.SetLogOutputFile(viper.GetString("log"))
			core.SetProcessorLogOutputFile(viper.GetString("log"))
		}
		log := core.Log()

		var locations lib.Locations
		cwd, _ := os.Getwd()

		if len(args) == 0 {
			for _, v := range viper.GetStringSlice("config") {
				loc, _ := lib.NewLocation(v, cwd)
				locations.AddLocation(loc)
			}
		} else {
			for _, v := range args {
				var loc *lib.Location
				var err error
				loc, err = lib.NewLocation(v, cwd)
				if err != nil {
					// is a content ?
					loc, err = lib.NewLocationContent(v, cwd)
					if err != nil {
						return
					}
				}

				locations.AddLocation(loc)
			}
		}

		if true == viper.IsSet("data") {
			core.SetDataLocation(viper.GetString("data"))
		} else {
			core.SetDataLocation(filepath.Join(cwd, "data"))
		}

		if false == viper.GetBool("no-network") {

			handlers := []core.FnMux{}
			handlers = append(handlers, core.WebHookServer())
			handlers = append(handlers, core.ApiServer(api.Handler(plugins)))

			if true == viper.IsSet("prometheus") {
				handlers = append(handlers, core.PrometheusServer(viper.GetString("prometheus.path")))
			}

			core.ListenAndServe(viper.GetString("host"), handlers...)
		}

		core.Log().Debugln("bitfan started")

		for _, loc := range locations.Items {
			agt, err := loc.ConfigAgents()

			if err != nil {
				log.Errorf("Error : %s %s", loc.Path, err)
				os.Exit(2)
			}
			ppl := loc.ConfigPipeline()

			// Allow pipeline customisation only when only one location was provided by user
			if len(locations.Items) == 1 {
				if cmd.Flags().Changed("name") {
					ppl.Name, _ = cmd.Flags().GetString("name")
				}
				if cmd.Flags().Changed("id") {
					ppl.ID, _ = cmd.Flags().GetInt("id")
				}
			}

			_, err = core.StartPipeline(&ppl, agt)
			if err != nil {
				log.Errorf("error : %s\n", err.Error())
				os.Exit(1)
			}
		}

		core.Log().Infoln("bitfan ready")

		if service.Interactive() {
			// Wait for signal CTRL+C for send a stop event to all AgentProcessor
			// When CTRL+C, SIGINT and SIGTERM signal occurs
			// Then stop server gracefully
			ch := make(chan os.Signal)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			<-ch
			close(ch)

			log.Println("")
			log.Printf("BitFan is stopping...")
			core.Stop()
			log.Printf("Everything stopped gracefully. Goodbye!\n")
		}

	},
}

func initRunConfig(cmd *cobra.Command) {
	viper.BindPFlag("api", cmd.Flags().Lookup("api"))
	viper.BindPFlag("prometheus", cmd.Flags().Lookup("prometheus"))
	viper.BindPFlag("prometheus.listen", cmd.Flags().Lookup("prometheus.listen"))
	viper.BindPFlag("prometheus.path", cmd.Flags().Lookup("prometheus.path"))
	viper.BindPFlag("webhook.listen", cmd.Flags().Lookup("webhook.listen"))
	viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	viper.BindPFlag("no-network", cmd.Flags().Lookup("no-network"))
	viper.BindPFlag("data", cmd.Flags().Lookup("data"))
}

func initRunFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")

	cmd.Flags().Bool("no-network", false, "Disable network (api and webhook)")
	cmd.Flags().String("name", "", "set pipeline's name")
	cmd.Flags().String("id", "", "set pipeline's id")
	cmd.Flags().String("data", "", "Path to data dir")

	cmd.Flags().Bool("api", true, "Expose REST Api")
	cmd.Flags().Bool("prometheus", false, "Export stats using prometheus output")
	cmd.Flags().String("prometheus.path", "/metrics", "Expose Prometheus metrics at specified path.")
}
