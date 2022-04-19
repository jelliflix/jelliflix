package logger

import (
	"github.com/jelliflix/jelliflix/infrastructure/environment"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var Log = newLogger()

func newLogger() (l *logrus.Logger) {
	l = logrus.New()

	if environment.IsDEV() {
		l.SetLevel(logrus.DebugLevel)
		l.SetFormatter(&nested.Formatter{HideKeys: true})
	} else {
		l.SetLevel(logrus.ErrorLevel)
		l.SetFormatter(&logrus.JSONFormatter{})
	}

	return
}
