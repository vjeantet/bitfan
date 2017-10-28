package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vjeantet/bitfan/api"
	"github.com/vjeantet/bitfan/lib"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start [config]",
	Aliases: []string{"add", "create"},
	Short:   "Start a new pipeline to the running bitfan",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli := api.NewRestClient(viper.GetString("host"))

		var locations lib.Locations
		cwd, _ := os.Getwd()
		for _, v := range args {
			var loc *lib.Location
			var err error
			loc, err = lib.NewLocation(v, cwd)
			if err != nil {
				// is a content ?
				loc, err = lib.NewLocationContent(v, cwd)
				if err != nil {
					return
				}
			}
			locations.AddLocation(loc)
		}

		for _, loc := range locations.Items {

			nPipeline := &api.Pipeline{}
			// Allow pipeline customisation only when only one location was provided by user
			if len(locations.Items) == 1 {
				if cmd.Flags().Changed("name") {
					nPipeline.Label, _ = cmd.Flags().GetString("name")
				}
				if cmd.Flags().Changed("id") {
					nPipeline.ID, _ = cmd.Flags().GetInt("id")
				}
			}

			if loc.Kind == lib.CONTENT_INLINE {
				nPipeline.Content = loc.Content
			} else if loc.Kind == lib.CONTENT_FS {
				nPipeline.ConfigLocation = filepath.Base(loc.Path)
				for path, b64Content := range loc.AssetsContent() {
					nPipeline.Assets = append(nPipeline.Assets, api.Asset{Path: path, Content: b64Content})
				}
			} else {
				nPipeline.ConfigLocation = loc.Path
			}

			pipeline, err := cli.AddPipeline(nPipeline)

			if err != nil {
				fmt.Printf("error : %v\n", err)
				os.Exit(1)
			} else {
				fmt.Printf("Started (ID:%d) - %s\n", pipeline.ID, pipeline.Label)
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
}
