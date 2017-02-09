package handler

import "log"

type LoggerHandler struct {
	logger *log.Logger
}

func (hnd *LoggerHandler) Deliver(message string) error {
	hnd.logger.Printf("Message:\n%q", message)

	return nil
}

func (hnd *LoggerHandler) Describe() string {
	return "Logger Handler"
}

func NewLoggerHandler(out *log.Logger) *LoggerHandler {
	return &LoggerHandler{logger: out}
}
