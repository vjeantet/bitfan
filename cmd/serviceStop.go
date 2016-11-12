package cmd

import (
	"log"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

func init() {
	serviceCmd.AddCommand(serviceStopCmd)
}

// serviceStopCmd represents the serviceStop command
var serviceStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running logfan service",

	Run: func(cmd *cobra.Command, args []string) {
		s := getService(nil)
		if service.Interactive() {
			s.Stop()
			log.Println("stop signal sent to logfan service")
		}
	},
}
