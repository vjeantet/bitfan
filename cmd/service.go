package cmd

import (
	"log"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	runtime "github.com/veino/veino/runtime"
)

var slogger service.Logger

type sprogram struct{}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"s", "svc"},
	Short:   "Install and manage logfan service",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serviceCmd.PersistentFlags().StringP("name", "n", "logfan", "service name")
	viper.BindPFlag("service.name", serviceCmd.PersistentFlags().Lookup("name"))

	serviceCmd.PersistentFlags().String("description", "Logfan is Logstash implementation on Golang", "service description")
	viper.BindPFlag("service.description", serviceCmd.PersistentFlags().Lookup("name"))
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getServiceConfig() *service.Config {
	return &service.Config{
		Name:        viper.GetString("service.name"),
		DisplayName: viper.GetString("service.name"),
		Description: viper.GetString("service.description"),
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
		log.Fatal(err)
	}

	slogger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	return s

}

func (p *sprogram) Start(s service.Service) error {
	go Execute()
	slogger.Info("logfan service started")
	return nil
}

func (p *sprogram) Stop(s service.Service) error {
	runtime.Stop()
	slogger.Info("Logfan Stopped")
	return nil
}
