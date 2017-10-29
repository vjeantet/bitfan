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

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vjeantet/bitfan/api"
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
		cli := api.NewRestClient(viper.GetString("host"))
		pipelines, err := cli.ListPipelines()
		if err != nil {
			fmt.Printf("list error: %v\n", err.Error())
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{
				"UUID",
				"name",
				"configuration",
			})

			for _, pipeline := range pipelines {
				host := ""
				if pipeline.ConfigHostLocation != fqdn.Get() {
					host = pipeline.ConfigHostLocation + "@"
				}

				table.Append([]string{
					pipeline.Uuid,
					pipeline.Label,
					fmt.Sprintf("%s%s",
						host,
						pipeline.ConfigLocation),
				})
			}
			//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("+")
			table.Render()

		}
	},
}
