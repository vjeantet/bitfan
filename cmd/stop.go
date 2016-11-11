package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/veino/logfan/lib"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [pipelineID]",
	Short: "stop a running pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		s := lib.ApiClient("127.0.0.1:1234")
		for _, ID := range args {
			// Send a request & read result
			retour := false
			if err := s.Request("stopPipeline", ID, &retour); err != nil {
				fmt.Printf("error : %s\n", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("pipeline %s stopped\n", ID)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
