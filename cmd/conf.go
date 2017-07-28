package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vjeantet/bitfan/api"
)

// stopCmd represents the stop command
var confCmd = &cobra.Command{
	Use:   "conf [pipelineID]",
	Short: "Retreive configuration file and its related files of a running pipeline",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli := api.NewRestClient(viper.GetString("host"))

		for _, ID := range args {
			// Send a request & read result
			pipeline, err := cli.Pipeline(ID, true)
			if err != nil {
				fmt.Printf("error : %s\n", err.Error())
				os.Exit(1)
			}
			uid := filepath.Base(filepath.Dir(pipeline.ConfigLocation))

			cwd, _ := os.Getwd()
			cwd = filepath.Join(cwd, uid)
			os.MkdirAll(cwd, os.ModePerm)

			//Save assets to cwd + pipeline uuid
			for _, asset := range pipeline.Assets {
				dest := filepath.Join(cwd, asset.Path)
				dir := filepath.Dir(dest)
				os.MkdirAll(dir, os.ModePerm)
				b64Decode(asset.Content, dest)
				// fmt.Printf("%s\n", filepath.Join(uid, asset.Path))
				fmt.Printf("%s\n", dest)
			}
		}
	},
}

func b64Decode(code string, dest string) error {
	buff, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(dest, buff, 07440); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCmd.AddCommand(confCmd)
	confCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
}
