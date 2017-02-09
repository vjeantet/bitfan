package version

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

const (
	VERSION = "0.3.0"
)

var app, version, versionShort string

func App() string {
	return app
}

func Version() string {
	return version
}

func VersionShort() string {
	return versionShort
}

func init() {
	app = path.Base(os.Args[0])
	versionShort = fmt.Sprintf("%s/%s", app, VERSION)
	version = fmt.Sprintf("%s/%s %s/%s %s", app, VERSION, runtime.GOOS, runtime.GOARCH, runtime.Version())
}
