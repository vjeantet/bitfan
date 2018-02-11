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

package commands

import (
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vjeantet/bitfan/api"
	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/entrypoint"
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

		opt := core.Options{
			VerboseLog:   viper.GetBool("verbose"),
			Debug:        viper.GetBool("debug"),
			LogFile:      viper.GetString("log"),
			DataLocation: viper.GetString("data"),
			Host:         viper.GetString("host"),
		}

		if !viper.GetBool("no-network") {
			opt.HttpHandlers = append(opt.HttpHandlers, core.HTTPHandler("/api/v2/", api.Handler("api/v2")))
			opt.HttpHandlers = append(opt.HttpHandlers, core.HTTPHandler("/public/",
				http.StripPrefix("/public/", http.FileServer(http.Dir(viper.GetString("public")))),
			))

			if viper.IsSet("prometheus") {
				opt.Prometheus = viper.GetString("prometheus.path")
			}
		}

		core.Start(opt)
		core.Log().Infoln("bitfan ready")

		// Start Pipelines

		// Prepare entrypoints
		var entrypoints entrypoint.EntrypointList
		//	From Storage when len == 0
		if len(args) == 0 {
			pipelinesToStart := core.Storage().FindPipelinesWithAutoStart(true)
			for _, p := range pipelinesToStart {
				entryPointPath, err := core.Storage().PreparePipelineExecutionStage(&p)
				if err != nil {
					core.Log().Fatalln(err)
				}

				var loc *entrypoint.Entrypoint
				loc, err = entrypoint.New(entryPointPath, "", entrypoint.CONTENT_REF)
				loc.PipelineName = p.Label
				loc.PipelineUuid = p.Uuid
				if err != nil {
					core.Log().Fatalln(err)
				}
				entrypoints.AddEntrypoint(loc)
			}
		}

		//	From config when config == 0
		cwd, _ := os.Getwd()
		if len(args) == 0 {
			for _, v := range viper.GetStringSlice("config") {
				loc, _ := entrypoint.New(v, cwd, entrypoint.CONTENT_REF)
				entrypoints.AddEntrypoint(loc)
			}
		}

		//	From args when config > 0
		if len(args) > 0 {
			for _, v := range args {
				var loc *entrypoint.Entrypoint
				var err error
				loc, err = entrypoint.New(v, cwd, entrypoint.CONTENT_REF)
				if err != nil {
					// is a content ?
					loc, err = entrypoint.New(v, cwd, entrypoint.CONTENT_INLINE)
					if err != nil {
						core.Log().Fatalln(err)
					}
				}
				entrypoints.AddEntrypoint(loc)
			}
		}

		for _, ep := range entrypoints.Items {
			ppl, err := ep.Pipeline()
			if err != nil {
				core.Log().Fatalln(err)
			}

			nUUID, err := ppl.Start()
			if err != nil {
				core.Log().Errorf("error: %v", err)
				os.Exit(1)
			}
			core.Log().Infof("Pipeline started %s (%s)(%s)", ppl.Label, ppl.Uuid, nUUID)

			// agt, err := ep.ConfigAgents()

			// if err != nil {
			// 	core.Log().Errorf("Error : %s %v", ep.Path, err)
			// 	os.Exit(2)
			// }
			// ppl := ep.ConfigPipeline()
			// _, err = core.StartPipeline(&ppl, agt)
			// if err != nil {
			// 	core.Log().Errorf("error: %v", err)
			// 	os.Exit(1)
			// }
			// core.Log().Infof("Pipeline started %s (%s)", ppl.Name, ppl.Uuid)
		}

		if service.Interactive() {
			// Wait for signal CTRL+C for send a stop event to all AgentProcessor
			// When CTRL+C, SIGINT and SIGTERM signal occurs
			// Then stop server gracefully
			ch := make(chan os.Signal)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			<-ch
			close(ch)

			core.Log().Println("")
			core.Log().Printf("BitFan is stopping...")
			core.Stop()
			core.Log().Printf("Everything stopped gracefully. Goodbye!")
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
	viper.BindPFlag("public", cmd.Flags().Lookup("public"))
}

func initRunFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")

	cmd.Flags().Bool("no-network", false, "Disable network (api and webhook)")
	cwd, _ := os.Getwd()
	cmd.Flags().String("data", filepath.Join(cwd, ".bitfan"), "Path to data dir")
	cmd.Flags().String("public", filepath.Join(cwd, "public"), "Path to public dir with served as /public/")
	cmd.Flags().Bool("api", true, "Expose REST Api")
	cmd.Flags().Bool("prometheus", false, "Export stats using prometheus output")
	cmd.Flags().String("prometheus.path", "/metrics", "Expose Prometheus metrics at specified path.")
}
