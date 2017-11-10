package api

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// hookConfig stores configuration needed to setup the hook
type hookConfig struct {
	Size int
}

// logHook to sends logs to buffer server
type logHook struct {
	config *hookConfig
	slice  [][]byte
	c      chan []byte
}

// newHook creates a hook to be added to an instance of logger
func newHook(config hookConfig) (*logHook, error) {
	if config.Size == 0 {
		config.Size = 50
	}

	return &logHook{
		c:      make(chan []byte),
		config: &config,
	}, nil
}

func (l *logHook) AddChan(c chan []byte) {
	l.c = c
}

func (l *logHook) String() [][]byte {
	return l.slice
}

// Fire is called when a log event is fired.
func (l *logHook) Fire(entry *logrus.Entry) error {
	serialized, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	go func(serialized []byte) {
		l.c <- serialized
	}(serialized)
	l.slice = append(l.slice, serialized)
	if len(l.slice) > l.config.Size {
		l.slice = l.slice[1:]
	}

	return nil
}

// Levels returns the available logging levels.
func (l *logHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
