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
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc [plugin]",
	Short: "Display documentation about plugins",
	Long: `Display list of available outputs, filters and outputs
	doc

Display documentation about the "date" plugin
	doc date

Display only a configuration blueprint for the "date" plugin
	doc date -t

Display list of only available filters
	doc --type=filter

Display documentation about the "file" plugin (the output one)
	doc file --type=output
	`,
	Run: func(cmd *cobra.Command, args []string) {
		kind, _ := cmd.Flags().GetString("type")
		if len(args) == 0 {
			switch kind {
			case "input":
				fallthrough
			case "output":
				fallthrough
			case "filter":
				listPlugins(kind)
			default:
				listAllPlugins()
			}
			return
		}

		name := args[0]
		tplOnly, _ := cmd.Flags().GetBool("template")
		if kind == "" {
			if _, ok := plugins["input"][name]; ok {
				kind = "input"
			} else if _, ok := plugins["filter"][name]; ok {
				kind = "filter"
			} else if _, ok := plugins["output"][name]; ok {
				kind = "output"
			}
		}

		// Affiche la doc d'un plugin
		err := displaydoc(kind, name, tplOnly)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("\n\n")
	},
}

func listAllPlugins() {
	listPlugins("input")
	fmt.Print("\n\n")
	listPlugins("filter")
	fmt.Print("\n\n")
	listPlugins("output")
	fmt.Print("\n\n")
}
func listPlugins(kind string) {
	fmt.Printf("# %s\n\n", strings.ToUpper(kind))
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Plugin", "Description"})
	for name, proc := range plugins[kind] {
		if name == "when" {
			continue
		}

		table.Append([]string{
			name, proc().Doc().DocShort,
		})
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.Render()
}

func displaydoc(kind string, name string, tplOnly bool) error {
	if _, ok := plugins[kind][name]; !ok {
		return fmt.Errorf("Unknow plugin %s in %s \n", name, kind)
	}

	p := plugins[kind][name]().Doc()

	if p.Name == "" {
		return fmt.Errorf("no doc available for %s %s\n go to github and open an issue :-(\n", name, kind)
	}

	if tplOnly {
		fmt.Print(string(p.GenExample("logstash")))
		return nil
	}

	w := p.GenMarkdown("logstash")
	fmt.Print(string(w))
	return nil
}

func init() {
	RootCmd.AddCommand(docCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	docCmd.Flags().BoolP("template", "t", false, "show only a template")
	docCmd.Flags().String("type", "", "input ? output ? filter ? (plugin may have the same name in multiple sections)")

}
