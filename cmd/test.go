package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vjeantet/bitfan/core"
	config "github.com/vjeantet/bitfan/core/config"
	"github.com/vjeantet/bitfan/lib"
)

func init() {
	RootCmd.AddCommand(testCmd)
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test configurations (files, url, directories)",
	Run: func(cmd *cobra.Command, args []string) {

		var locations lib.Locations
		cwd, _ := os.Getwd()
		for _, v := range args {
			var loc *lib.Location
			var err error
			loc, err = lib.NewLocation(v, cwd)
			if err != nil {
				loc, err = lib.NewLocationContent(v, cwd)
				if err != nil {
					return
				}
			}

			locations.AddLocation(loc)
		}

		var cko int
		var ctot int
		for _, loc := range locations.Items {
			ctot++
			err := testConfigContent(loc)
			if err != nil {
				fmt.Printf("%s\n -> %s\n\n", loc.Path, err)
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

func testConfigContent(loc *lib.Location) error {
	configAgents, err := loc.ConfigAgents()
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	configAgentsOrdered := config.Sort(configAgents, config.SortInputsFirst)
	for _, configAgent := range configAgentsOrdered {
		_, err := core.NewAgent(configAgent)
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
	}

	return nil
}
