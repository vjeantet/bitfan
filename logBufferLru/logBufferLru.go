package logBufferLru

import (
	"fmt"

	"github.com/armon/circbuf"
	"github.com/sirupsen/logrus"
)

// HookConfig stores configuration needed to setup the hook
type HookConfig struct {
	Size int64
}

// RedisHook to sends logs to buffer server
type BufferLruHook struct {
	buf *circbuf.Buffer
	c   chan []byte
}

// NewHook creates a hook to be added to an instance of logger
func NewHook(config HookConfig) (*BufferLruHook, error) {
	if config.Size == 0 {
		config.Size = 100 * 50
	}
	cbuffer, err := circbuf.NewBuffer(config.Size)
	return &BufferLruHook{
		buf: cbuffer,
		c:   make(chan []byte),
	}, err

}

func (b *BufferLruHook) AddChan(c chan []byte) {
	b.c = c
}

func (b *BufferLruHook) String() string {
	return b.buf.String()
}

// Fire is called when a log event is fired.
func (hook *BufferLruHook) Fire(entry *logrus.Entry) error {
	msg := fmt.Sprintf("%s %s\n", entry.Time.Format("02-01-06 15:04:05 - "), entry.Message)
	go func(m []byte) { hook.c <- m }([]byte(msg))
	hook.buf.Write([]byte(msg))
	return nil
}

// Levels returns the available logging levels.
func (hook *BufferLruHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
