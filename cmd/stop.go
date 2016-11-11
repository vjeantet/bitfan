package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veino/logfan/lib"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [pipelineID]",
	Short: "Stop a running pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		s := lib.ApiClient(viper.GetString("host"))
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
	stopCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
	viper.BindPFlag("host", stopCmd.Flags().Lookup("host"))
}
