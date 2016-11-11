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
	"github.com/veino/logfan/lib"
	config "github.com/veino/veino/config"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List pipelines",
	Aliases: []string{"ls"},
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		s := lib.ApiClient("127.0.0.1:1234")

		// Send a request & read result
		var pipelines = config.PipelineList{}
		if err := s.Request("findPipelines", "", &pipelines); err != nil {
			fmt.Printf("list error: %v\n", err.Error())
		} else {

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{
				"ID",
				"name",
				"status",
				"processors",
				"description",
				"configuration",
			})
			for _, pipeline := range pipelines {
				table.Append([]string{
					pipeline.ID,
					pipeline.Name,
					pipeline.Status.String(),
					strconv.Itoa(len(pipeline.AgentsID)),
					pipeline.Description,
					pipeline.ConfigLocation,
				})
			}
			// table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("+")
			table.Render()

		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
