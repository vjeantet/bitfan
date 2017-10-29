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

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart [ID] [config]",
	Short: "Restart a pipeline",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("host", cmd.Flags().Lookup("host"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli := api.NewRestClient(viper.GetString("host"))

		// bitfan restart 23 simple.com
		if len(args) > 2 {
			fmt.Println("this command allow only one configuration file")
			os.Exit(1)
		}
		if len(args) == 0 {
			fmt.Println("give an pipeline uuid identifier")
			os.Exit(1)
		}

		var pipelineRef string
		pipelineRef = args[0]

		var confLocation string
		if len(args) == 2 {
			confLocation = args[1]
		}

		// get Pipeline by pipelineRef
		oldPipeline, err := cli.Pipeline(pipelineRef, false)
		if err != nil {
			fmt.Printf("Pipeline %s not found - %v\n", pipelineRef, err)
			os.Exit(1)
		}

		// stop pipeline

		err = cli.StopPipeline(pipelineRef)
		if err != nil {
			fmt.Printf("error : %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("pipeline %s stopped\n", pipelineRef)
		}

		// start pipeline

		if confLocation != "" { // with given conf
			nPipeline, _ := getPipeline(confLocation)

			pipeline, err := cli.AddPipeline(nPipeline)
			if err != nil {
				fmt.Printf("error : %v\n", err)
				os.Exit(1)
			} else {
				fmt.Printf("Started (UUID:%s) - %s\n", pipeline.Uuid, pipeline.Label)
			}
		} else { // with old conf
			nPipeline := &api.Pipeline{
				Label:              oldPipeline.Label,
				ConfigLocation:     oldPipeline.ConfigLocation,
				ConfigHostLocation: oldPipeline.ConfigHostLocation,
			}

			pipeline, err := cli.AddPipeline(nPipeline)
			if err != nil {
				fmt.Printf("error : %v\n", err)
				os.Exit(1)
			} else {
				fmt.Printf("Started (UUID:%s) - %s\n", pipeline.Uuid, pipeline.Label)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(restartCmd)
	restartCmd.Flags().StringP("host", "H", "127.0.0.1:5123", "Service Host to connect to")
}

func getPipeline(v string) (*api.Pipeline, error) {
	cwd, _ := os.Getwd()
	var locations lib.Locations
	var loc *lib.Location
	var err error

	loc, err = lib.NewLocation(v, cwd)
	if err != nil {
		// is a content ?
		loc, err = lib.NewLocationContent(v, cwd)
		if err != nil {
			return nil, err
		}
	}
	locations.AddLocation(loc)

	if len(locations.Items) > 1 {
		fmt.Println("this command allow only one configuration file")
	}

	loc = locations.Items[0]
	nPipeline := &api.Pipeline{}

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
	return nPipeline, nil
}
