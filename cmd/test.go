package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/veino/logfan/lib"
	config "github.com/veino/veino/config"
	runtime "github.com/veino/veino/runtime"
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
			locations.Add(v, cwd)
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

	configAgentsOrdered := config.Sort(configAgents, config.SortOutputsFirst)
	for _, configAgent := range configAgentsOrdered {
		_, err := runtime.NewAgent(configAgent, 0)
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
	}

	return nil
}
