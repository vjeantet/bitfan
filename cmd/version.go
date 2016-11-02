package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "No Version Provided"
var Buildstamp = ""

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
	fmt.Println("LogFan Version : " + Version + "\nUTC Build Time : " + Buildstamp)
}
