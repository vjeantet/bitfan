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
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veino/logfan/lib"
	config "github.com/veino/veino/config"
)

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List running pipelines",
	Aliases: []string{"ls"},
	Long:    ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		s := lib.ApiClient(viper.GetString("host"))

		// Send a request & read result
		var pipelines = config.PipelineList{}
		if err := s.Request("findPipelines", "", &pipelines); err != nil {
			fmt.Printf("list error: %v\n", err.Error())
		} else {

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{
				"ID",
				"name",
				"Started at",
				"processors",
				"description",
				"configuration",
			})
			for _, pipeline := range pipelines {
				table.Append([]string{
					pipeline.ID,
					pipeline.Name,
					pipeline.StartedAt.Format("2006-01-02 15:04:05"),
					strconv.Itoa(len(pipeline.AgentsID)),
					pipeline.Description,
					fmt.Sprintf("%s@%s",
						pipeline.ConfigHostLocation,
						pipeline.ConfigLocation),
				})
			}
			// table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("+")
			table.Render()

		}
	},
}
