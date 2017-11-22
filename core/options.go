package core

type Options struct {
	Host         string
	HttpHandlers []fnMux
	Debug        bool
	VerboseLog   bool
	LogFile      string
	DataLocation string
}
