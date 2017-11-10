package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vjeantet/bitfan/api/client"
)

// stopCmd represents the stop command
var confCmd = &cobra.Command{
	Use:   "conf [pipelineUUID]",
	Short: "Retrieve configuration file and its related files of a running pipeline",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli := client.New(viper.GetString("host"))

		for _, ID := range args {
			// Send a request & read result
			pipeline, err := cli.Pipeline(ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error : %v\n", err)
				os.Exit(1)
			}
			uid := pipeline.Uuid

			cwd, _ := os.Getwd()
			cwd = filepath.Join(cwd, uid)
			os.MkdirAll(cwd, os.ModePerm)

			//Save assets to cwd + pipeline uuid
			for _, asset := range pipeline.Assets {
				dest := filepath.Join(cwd, asset.Name)
				dir := filepath.Dir(dest)
				os.MkdirAll(dir, os.ModePerm)

				if err := ioutil.WriteFile(dest, asset.Value, 07440); err != nil {
					fmt.Fprintf(os.Stderr, "error : '%s' - %s \n", err.Error(), dest)
					os.Exit(1)
				}
				fmt.Printf("%s\n", dest)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(confCmd)
	confCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
}
