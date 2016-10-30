package cmd

import (
	"log"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
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
	serviceCmd.PersistentFlags().StringP("name", "n", "com.github.veino.logfan", "Logfan service's name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getServiceConfig() *service.Config {
	return &service.Config{
		Name:        "com.github.veino.logfan",
		DisplayName: "Logfan",
		Description: "Logfan is Logstash implementation on Golang",
	}
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

	var err error

	if configPath != "" {
		err = startLogfan(configPath, "", stats, []string{})
	} else {
		log.Fatalln("missing configuration location")
	}

	if err == nil {
		slogger.Info("logfan service started")
	}

	return err
}

func (p *sprogram) Stop(s service.Service) error {
	runtime.Stop()
	slogger.Info("Logfan Stopped")
	return nil
}
