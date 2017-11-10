package commands

import (
	"log"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

type sprogram struct{}

func init() {
	RootCmd.AddCommand(serviceCmd)
}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"s", "svc"},
	Short:   "Install and manage bitfan service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func getServiceConfig() *service.Config {
	return &service.Config{
		Name:        "bitfan-ui",
		DisplayName: "bitfan-ui",
		Description: "bitfan-ui",
	}
}
func GetService() service.Service {
	return getService(nil)
}
func getService(svcConfig *service.Config) service.Service {
	if svcConfig == nil {
		svcConfig = getServiceConfig()
	}

	cwd, _ := os.Getwd()
	svcConfig.WorkingDirectory = cwd
	prg := &sprogram{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal("getService New err : ", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return s

}

func (p *sprogram) Start(s service.Service) error {
	go Execute()
	return nil
}

func (p *sprogram) Stop(s service.Service) error {
	os.Exit(0)
	return nil
}
