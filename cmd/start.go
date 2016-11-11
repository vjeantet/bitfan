package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veino/logfan/lib"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new pipeline in a running logfan",
	Run: func(cmd *cobra.Command, args []string) {

		var locations lib.Locations
		cwd, _ := os.Getwd()
		for _, v := range args {
			locations.Add(v, cwd)
		}

		for _, loc := range locations.Items {
			agt, err := loc.ConfigAgents()

			if err != nil {
				fmt.Printf("Error : %s %s", loc.Path, err)
				os.Exit(2)
			}
			ppl := loc.ConfigPipeline()
			if cmd.Flags().Changed("name") {
				ppl.Name, _ = cmd.Flags().GetString("name")
			}

			if cmd.Flags().Changed("id") {
				ppl.ID, _ = cmd.Flags().GetString("id")
			}

			starter := &lib.ApiStarter{
				Pipeline: ppl,
				Agents:   agt,
			}

			s := lib.ApiClient(viper.GetString("host"))
			ID := ""
			if err := s.Request("startPipeline", starter, &ID); err != nil {
				fmt.Printf("error : %s\n", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("Started (PID:%s) - %s\n", ID, loc.Path)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().String("name", "", "set pipeline's name")
	startCmd.Flags().String("id", "", "set pipeline's id")
	startCmd.Flags().String("force", "", "force start even if duplicate")
	startCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
	viper.BindPFlag("host", startCmd.Flags().Lookup("host"))
}
