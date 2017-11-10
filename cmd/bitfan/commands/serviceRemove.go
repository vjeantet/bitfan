package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(serviceRemoveCmd)
}

// serviceRemoveCmd represents the serviceRemove command
var serviceRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove existing service",
	Run: func(cmd *cobra.Command, args []string) {
		s := getService(nil)
		s.Stop()
		if err := s.Uninstall(); err != nil {
			log.Fatal(err)
		}
		log.Println("service bitfan successfully removed")
		os.Exit(0)
	},
}
