package commands

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
	Short: "restart a bitfan ui service",

	Run: func(cmd *cobra.Command, args []string) {

		s := getService(nil)

		log.Println("stopping bitfan ui service...")
		err := s.Stop()
		if err != nil {
			log.Printf("stop service error : %v", err)
		} else {
			// log.Println("service bitfan stopped")
		}

		log.Println("starting bitfan ui service...")
		err = s.Start()
		if err != nil {
			log.Printf("start service error : %v", err)
		} else {
			// log.Println("service bitfan started")
		}

	},
}
