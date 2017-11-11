package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "No Version Provided"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Display version informations",

	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printVersion() {
	fmt.Println("BitFan UI Version : " + Version + "\n")
}
