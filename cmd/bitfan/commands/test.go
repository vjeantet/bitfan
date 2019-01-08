package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/awillis/bitfan/entrypoint"
)

func init() {
	RootCmd.AddCommand(testCmd)
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test configurations (files, url, directories)",
	Run: func(cmd *cobra.Command, args []string) {

		var locations entrypoint.EntrypointList
		cwd, _ := os.Getwd()
		for _, v := range args {
			var loc *entrypoint.Entrypoint
			var err error
			loc, err = entrypoint.New(v, cwd, entrypoint.CONTENT_REF)
			if err != nil {
				loc, err = entrypoint.New(v, cwd, entrypoint.CONTENT_INLINE)
				if err != nil {
					return
				}
			}

			locations.AddEntrypoint(loc)
		}

		var cko int
		var ctot int
		for _, loc := range locations.Items {
			ctot++
			err := testConfigContent(loc)
			if err != nil {
				fmt.Printf("%s\n -> %v\n\n", loc.Path, err)
				cko++
			}
		}

		if ctot == 0 {
			fmt.Println("No configuration available to test")
		} else if cko == 0 {
			fmt.Printf("Everything is ok, %d configurations checked\n", ctot)
		}

	},
}

func testConfigContent(loc *entrypoint.Entrypoint) error {
	// TODO : refactor with pipeline
	// configAgents, err := loc.ConfigAgents()
	// if err != nil {
	// 	return err
	// }

	// configAgentsOrdered := config.Sort(configAgents, config.SortInputsFirst)
	// for _, configAgent := range configAgentsOrdered {
	// 	if _, err := core.NewAgent(configAgent); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
