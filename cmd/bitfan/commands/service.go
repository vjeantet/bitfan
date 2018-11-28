package commands

import (
	"log"
	"os"

	"bitfan/core"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var slogger service.Logger

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
		initSettings(cmd)
		viper.BindPFlag("workers", cmd.Flags().Lookup("filterworkers"))
		viper.BindPFlag("log", cmd.Flags().Lookup("log"))
		viper.BindPFlag("verbose", cmd.Flags().Lookup("verbose"))
		viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func getServiceConfig() *service.Config {
	return &service.Config{
		Name:        "bitfan",
		DisplayName: "bitfan",
		Description: "bitfan",
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

	slogger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	return s

}

func (p *sprogram) Start(s service.Service) error {
	go Execute()
	slogger.Infof("Bitfan service started")
	return nil
}

func (p *sprogram) Stop(s service.Service) error {
	core.Stop()
	slogger.Info("Bitfan Stopped")
	return nil
}
