package cmd

import (
	"log"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

// serviceStopCmd represents the serviceStop command
var serviceStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop a running logfan service",

	Run: func(cmd *cobra.Command, args []string) {
		servicename, _ := cmd.Flags().GetString("name")
		svcConfig := getServiceConfig()
		svcConfig.Name = servicename
		svcConfig.DisplayName = servicename
		s := getService(svcConfig)

		if service.Interactive() {
			s.Stop()
			log.Println("service logfan started")
		}
	},
}

func init() {
	serviceCmd.AddCommand(serviceStopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceStopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceStopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
