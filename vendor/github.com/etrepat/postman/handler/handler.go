package handler

import "log"

const (
	POSTBACK_HANDLER = 1 << iota
	LOGGER_HANDLER
)

type MessageHandler interface {
	Deliver(message string) error
	Describe() string
}

func New(t uint, args ...interface{}) (hnd MessageHandler) {
	switch t {
	case POSTBACK_HANDLER:
		hnd = NewPostBackHandler(args[0].(string), args[1].(bool), args[2].(string))

	case LOGGER_HANDLER:
		hnd = NewLoggerHandler(args[0].(*log.Logger))
	}

	return hnd
}
