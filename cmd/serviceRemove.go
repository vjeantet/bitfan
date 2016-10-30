package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// serviceRemoveCmd represents the serviceRemove command
var serviceRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove existing service",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		servicename, _ := cmd.Flags().GetString("name")
		svcConfig := getServiceConfig()
		svcConfig.Name = servicename
		svcConfig.DisplayName = servicename
		s := getService(svcConfig)

		s.Stop()
		if err := s.Uninstall(); err != nil {
			log.Fatal(err)
		}
		log.Println("service logfan successfully removed")
		os.Exit(0)
	},
}

func init() {
	serviceCmd.AddCommand(serviceRemoveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceRemoveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceRemoveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
