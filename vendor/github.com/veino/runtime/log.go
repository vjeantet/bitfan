package runtime

import (
	"github.com/Sirupsen/logrus"
	"github.com/veino/veino"
)

var logger = &vienoLogger{logrus.New()}

func init() {
	logger.Level = logrus.WarnLevel
}

func Logger() *vienoLogger {
	return logger
}

type vienoLogger struct {
	*logrus.Logger
}

func (v *vienoLogger) WithEvent(e veino.IPacket) *logrus.Entry {
	return v.WithFields(logrus.Fields{
		"kind":          e.Kind(),
		"field.message": e.Fields().ValueOrEmptyForPathString("message"),
	})
}

func (v *vienoLogger) SetDebugMode(debug bool) {
	if debug == true {
		v.Level = logrus.DebugLevel
	}
}
