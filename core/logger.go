package core

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/vjeantet/bitfan/processors"
)

var logger = logrus.New()
var loggerProcessor = logrus.New()

func init() {
	logger.Level = logrus.WarnLevel
	loggerProcessor.Level = logrus.WarnLevel
	loggerProcessor.Formatter = &ProcessorFormatter{formatter: &logrus.TextFormatter{}}
}

type ProcessorFormatter struct {
	formatter  logrus.Formatter
	Xpipeline  string
	Xagent     string
	Xprocessor string
}

func (p *ProcessorFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if v, ok := entry.Data["Xprocessor"]; ok {
		entry.Message = fmt.Sprintf(`[%s] %s`,
			v,
			entry.Message,
		)
	}
	return p.formatter.Format(entry)
}

func Log() *logrus.Logger {
	return logger
}

func LogWithEvent(e processors.IPacket) *logrus.Entry {
	return Log().WithFields(logrus.Fields{
		"field.message": e.Fields().ValueOrEmptyForPathString("message"),
	})
}

func SetLogDebugMode(components []string) {

	for _, c := range components {
		switch c {
		case "core":
			Log().Level = logrus.DebugLevel
		case "processors":
			loggerProcessor.Level = logrus.DebugLevel
		}
	}
}

func SetLogVerboseMode(components []string) {
	for _, c := range components {
		switch c {
		case "core":
			Log().Level = logrus.InfoLevel
		case "processors":
			loggerProcessor.Level = logrus.InfoLevel
		}
	}
}

func SetLogOutputFile(fileLocation string) {
	Log().Out = ioutil.Discard
	Log().Formatter = &logrus.TextFormatter{DisableColors: true}
	f, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

	}
	Log().Out = f
}

func SetProcessorLogOutputFile(fileLocation string) {
	loggerProcessor.Formatter = &ProcessorFormatter{formatter: &logrus.TextFormatter{DisableColors: true}}
	f, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

	}
	loggerProcessor.Out = f
}

type processorLogger struct {
	e *logrus.Entry
}

func newProcessorLogger(procType, agentID, pipelineID string) *processorLogger {
	return &processorLogger{
		e: loggerProcessor.WithFields(map[string]interface{}{
			"Xprocessor": procType,
			"pipeline":   pipelineID,
			"agent":      agentID,
		}),
	}
}

func (p *processorLogger) Debug(args ...interface{}) {
	p.e.Debug(args...)
}
func (p *processorLogger) Debugf(format string, args ...interface{}) {
	p.e.Debugf(format, args...)
}
func (p *processorLogger) Debugln(args ...interface{}) {
	p.e.Debugln(args...)
}

func (p *processorLogger) Error(args ...interface{}) {
	p.e.Error(args...)
}
func (p *processorLogger) Errorf(format string, args ...interface{}) {
	p.e.Errorf(format, args...)
}
func (p *processorLogger) Errorln(args ...interface{}) {
	p.e.Errorln(args...)
}

func (p *processorLogger) Fatal(args ...interface{}) {
	p.e.Fatal(args...)
}
func (p *processorLogger) Fatalf(format string, args ...interface{}) {
	p.e.Fatalf(format, args...)
}
func (p *processorLogger) Fatalln(args ...interface{}) {
	p.e.Fatalln(args...)
}

func (p *processorLogger) Info(args ...interface{}) {
	p.e.Info(args...)
}
func (p *processorLogger) Infof(format string, args ...interface{}) {
	p.e.Infof(format, args...)
}
func (p *processorLogger) Infoln(args ...interface{}) {
	p.e.Infoln(args...)
}

func (p *processorLogger) Panic(args ...interface{}) {
	p.e.Panic(args...)
}
func (p *processorLogger) Panicf(format string, args ...interface{}) {
	p.e.Panicf(format, args...)
}
func (p *processorLogger) Panicln(args ...interface{}) {
	p.e.Panicln(args...)
}

func (p *processorLogger) Print(args ...interface{}) {
	p.e.Print(args...)
}
func (p *processorLogger) Printf(format string, args ...interface{}) {
	p.e.Printf(format, args...)
}
func (p *processorLogger) Println(args ...interface{}) {
	p.e.Println(args...)
}

func (p *processorLogger) Warn(args ...interface{}) {
	p.e.Warn(args...)
}

func (p *processorLogger) Warnf(format string, args ...interface{}) {
	p.e.Warnf(format, args...)
}
func (p *processorLogger) Warning(args ...interface{}) {
	p.e.Warning(args...)
}
func (p *processorLogger) Warningf(format string, args ...interface{}) {
	p.e.Warningf(format, args...)
}
func (p *processorLogger) Warningln(args ...interface{}) {
	p.e.Warningln(args...)
}
func (p *processorLogger) Warnln(args ...interface{}) {
	p.e.Warnln(args...)
}
