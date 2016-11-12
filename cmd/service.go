package cmd

import (
	"log"
	"os"

	"github.com/k0kubun/pp"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	runtime "github.com/veino/veino/runtime"
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
	Short:   "Install and manage logfan service",
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
		Name:        "logfan",
		DisplayName: "logfan",
		Description: "logfan",
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
		pp.Println("svcConfig-->", svcConfig)
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
	slogger.Infof("Logfan service started")
	return nil
}

func (p *sprogram) Stop(s service.Service) error {
	runtime.Stop()
	slogger.Info("Logfan Stopped")
	return nil
}
