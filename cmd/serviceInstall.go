package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// serviceInstallCmd represents the serviceInstall command
var serviceInstallCmd = &cobra.Command{
	Use:   "install [conf location]",
	Short: "install logfan as a service",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		cwd, _ := os.Getwd()
		configLocation := ""

		if len(args) == 1 {
			configLocation = args[0]
		} else {
			configLocation, _ = os.Getwd()
		}

		log.Println("Configuration location : ", configLocation)

		servicename, _ := cmd.Flags().GetString("name")
		svcConfig := getServiceConfig()
		svcConfig.Name = servicename
		svcConfig.DisplayName = servicename
		svcConfig.Arguments = []string{"-f", configLocation}
		svcConfig.WorkingDirectory = cwd

		s := getService(svcConfig)

		// if _, err := os.Stat(configPath); err != nil {
		// 	log.Fatalf("ERROR file or directory does not exist [%s]", absConfigPath)
		// }

		if err := s.Install(); err != nil {
			log.Fatal(err)
		}
		log.Println("service logfan successfully installed")
		os.Exit(0)

	},
}

func init() {
	serviceCmd.AddCommand(serviceInstallCmd)
	// serviceInstallCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
