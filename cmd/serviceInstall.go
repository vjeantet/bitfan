package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	serviceCmd.AddCommand(serviceInstallCmd)
	initRunFlags(serviceInstallCmd)
}

// serviceInstallCmd represents the serviceInstall command
var serviceInstallCmd = &cobra.Command{
	Use:   "install [config1] [config2]",
	Short: "install logfan as a service",
	PreRun: func(cmd *cobra.Command, args []string) {
		initRunConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		svcConfig := getServiceConfig()
		svcConfig.WorkingDirectory = cwd
		svcConfig.Arguments = append([]string{"run"}, args...)
		svcConfig.Option = service.KeyValue{
			"RunAtLoad": true,
			"KeepAlive": false,
		}

		cmd.Flags().Visit(func(f *pflag.Flag) {
			svcConfig.Arguments = append(svcConfig.Arguments, fmt.Sprintf("--%s=%s", f.Name, f.Value))
		})

		s := getService(svcConfig)

		if err := s.Install(); err != nil {
			log.Fatal(err)
		}

		log.Printf("service '%s' successfully installed", svcConfig.Name)
	},
}
