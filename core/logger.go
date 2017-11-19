package core

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *Logger

func init() {
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetFormatter(&bitfanFormatter{formatter: &logrus.TextFormatter{}})
	logger = NewLogger("core", nil)
}

type bitfanFormatter struct {
	formatter logrus.Formatter
}

func (p *bitfanFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if v, ok := entry.Data["component"]; ok {
		entry.Message = fmt.Sprintf(`[%s] %s`,
			v,
			entry.Message,
		)
	} else {
		entry.Message = fmt.Sprintf(`%s`, entry.Message)
	}
	return p.formatter.Format(entry)
}

//Should be unexported
func Log() *Logger {
	return logger
}

func setLogDebugMode() {
	logrus.SetLevel(logrus.DebugLevel)
}

func setLogVerboseMode() {
	logrus.SetLevel(logrus.InfoLevel)
}

func setLogOutputFile(fileLocation string) {
	logrus.SetOutput(ioutil.Discard)
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	f, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Log().Errorf("Error while opening log file %v", err)
	}
	logrus.SetOutput(f)
}

type Logger struct {
	e *logrus.Entry
}

func NewLogger(component string, data map[string]interface{}) *Logger {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["component"] = component
	return &Logger{
		e: logrus.WithFields(data),
	}
}

func (p *Logger) Debug(args ...interface{}) {
	p.e.Debug(args...)
}
func (p *Logger) Debugf(format string, args ...interface{}) {
	p.e.Debugf(format, args...)
}
func (p *Logger) Debugln(args ...interface{}) {
	p.e.Debugln(args...)
}

func (p *Logger) Error(args ...interface{}) {
	p.e.Error(args...)
}
func (p *Logger) Errorf(format string, args ...interface{}) {
	p.e.Errorf(format, args...)
}
func (p *Logger) Errorln(args ...interface{}) {
	p.e.Errorln(args...)
}

func (p *Logger) Fatal(args ...interface{}) {
	p.e.Fatal(args...)
}
func (p *Logger) Fatalf(format string, args ...interface{}) {
	p.e.Fatalf(format, args...)
}
func (p *Logger) Fatalln(args ...interface{}) {
	p.e.Fatalln(args...)
}

func (p *Logger) Info(args ...interface{}) {
	p.e.Info(args...)
}
func (p *Logger) Infof(format string, args ...interface{}) {
	p.e.Infof(format, args...)
}
func (p *Logger) Infoln(args ...interface{}) {
	p.e.Infoln(args...)
}

func (p *Logger) Panic(args ...interface{}) {
	p.e.Panic(args...)
}
func (p *Logger) Panicf(format string, args ...interface{}) {
	p.e.Panicf(format, args...)
}
func (p *Logger) Panicln(args ...interface{}) {
	p.e.Panicln(args...)
}

func (p *Logger) Print(args ...interface{}) {
	p.e.Print(args...)
}
func (p *Logger) Printf(format string, args ...interface{}) {
	p.e.Printf(format, args...)
}
func (p *Logger) Println(args ...interface{}) {
	p.e.Println(args...)
}

func (p *Logger) Warn(args ...interface{}) {
	p.e.Warn(args...)
}

func (p *Logger) Warnf(format string, args ...interface{}) {
	p.e.Warnf(format, args...)
}
func (p *Logger) Warning(args ...interface{}) {
	p.e.Warning(args...)
}
func (p *Logger) Warningf(format string, args ...interface{}) {
	p.e.Warningf(format, args...)
}
func (p *Logger) Warningln(args ...interface{}) {
	p.e.Warningln(args...)
}
func (p *Logger) Warnln(args ...interface{}) {
	p.e.Warnln(args...)
}
