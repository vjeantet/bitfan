package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vjeantet/bitfan/api"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:     "stop [pipelineID]",
	Short:   "Stop a running pipeline",
	Aliases: []string{"remove", "rm", "delete"},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli := api.NewRestClient(viper.GetString("host"))

		for _, ID := range args {
			// Send a request & read result
			err := cli.StopPipeline(ID)
			if err != nil {
				fmt.Printf("error : %v\n", err)
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
}
