package core

type Options struct {
	AutoStart    bool
	Host         string
	HttpHandlers []fnMux
	Debug        bool
	VerboseLog   bool
	LogFile      string
	DataLocation string
}
