package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// serviceStartCmd represents the serviceStart command
var serviceStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start an installed logfan service",
	Run: func(cmd *cobra.Command, args []string) {
		servicename, _ := cmd.Flags().GetString("name")
		svcConfig := getServiceConfig()
		svcConfig.Name = servicename
		svcConfig.DisplayName = servicename
		s := getService(svcConfig)
		err := s.Start()
		if err != nil {
			log.Printf("start service error : %s", err)
		} else {
			log.Println("service logfan started")
		}

	},
}

func init() {
	serviceCmd.AddCommand(serviceStartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceStartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceStartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
