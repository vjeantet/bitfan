package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(serviceRestartCmd)
}

// serviceRestartCmd represents the serviceRestart command
var serviceRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart a bitfan service",

	Run: func(cmd *cobra.Command, args []string) {

		s := getService(nil)

		log.Println("stopping bitfan service...")
		err := s.Stop()
		if err != nil {
			log.Printf("stop service error : %s", err)
		} else {
			// log.Println("service bitfan stopped")
		}

		log.Println("starting bitfan service...")
		err = s.Start()
		if err != nil {
			log.Printf("start service error : %s", err)
		} else {
			// log.Println("service bitfan started")
		}

	},
}
