package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vjeantet/bitfan/api/client"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:     "stop [pipelineUUID]",
	Short:   "Stop a running pipeline",
	Aliases: []string{"remove", "rm", "delete"},
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli := client.New(viper.GetString("host"))

		for _, uuid := range args {
			// Send a request & read result
			_, err := cli.StopPipeline(uuid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error : %s\n", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("pipeline %s stopped\n", uuid)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
}
