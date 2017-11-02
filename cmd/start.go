package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vjeantet/bitfan/api/client"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start [pipelineUUID]",
	Aliases: []string{"restart"},
	Short:   "Start a pipeline to the running bitfan",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	},
	Run: func(cmd *cobra.Command, args []string) {

		cli := client.New(viper.GetString("host"))

		for _, uuid := range args {
			// Send a request & read result
			pipeline, err := cli.StartPipeline(uuid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error : %s\n", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("Started (UUID:%s) - %s\n", pipeline.Uuid, pipeline.Label)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
}
