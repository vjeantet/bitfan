// +build windows

package cmd

import (
	"log"
	"os"

	"github.com/kardianos/service"
	"github.com/veino/runtime"
	"path/filepath"
)

var logger service.Logger

type program struct{}

func Execute() {
	parseFlags()

	workingDirectory, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    absConfigPath := filepath.Join(workingDirectory, configPath)

	svcConfig := &service.Config{
		Name:             "logfan",
		DisplayName:      "logfan",
		Description:      "Logfan is Logstash implementation on Golang",
		Arguments:        []string{"-f", absConfigPath},

	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if installWindowsServiceFlag {
		if _, err := os.Stat(configPath); err != nil {
			log.Fatalf("ERROR file or directory does not exist [%s]", absConfigPath)
		}
		if err := s.Install(); err != nil {
			log.Fatal(err)
		}
		log.Print("service logfan successfully installed")
		os.Exit(0)
	}
	if removeWindowsServiceFlag {
		s.Stop()
		if err := s.Uninstall(); err != nil {
			log.Fatal(err)
		}
		log.Print("service logfan successfully removed")
		os.Exit(0)
	}

	if !service.Interactive() {
		logger, err = s.Logger(nil)
		if err != nil {
			log.Fatal(err)
		}
		err = s.Run()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	}
	// Started in console
	startLogfanAndWait(configPath, configString, stats)
}

func (p *program) Start(s service.Service) error {
	err := startLogfan(configPath, configString, stats)
	if err != nil {
		logger.Info("Logfan Started")
	}
	return err
}

func (p *program) Stop(s service.Service) error {
	runtime.Stop()
	logger.Info("Logfan Stopped")
	return nil
}
