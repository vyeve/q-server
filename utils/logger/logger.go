package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	// EnvLogLevel is environment to change log level
	EnvLogLevel = "APP_LOG_LEVEL"
	// logTimeFormat represents time format in log messages
	logTimeFormat = "2006-01-02 15:04:05.99"
	// defaultLogLevel set as INFO
	defaultLogLevel = logrus.InfoLevel
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type loggerImpl struct {
	*logrus.Logger
}

func getLogLevel() logrus.Level {
	lvl, err := logrus.ParseLevel(os.Getenv(EnvLogLevel))
	if err != nil {
		return defaultLogLevel
	}
	return lvl
}

func New() Logger {
	lvl := getLogLevel()
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = logTimeFormat
	customFormatter.FullTimestamp = true
	log := logrus.New()
	log.SetFormatter(customFormatter)
	log.SetLevel(lvl)
	log.SetNoLock()
	return &loggerImpl{
		Logger: log,
	}
}
