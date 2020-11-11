package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(writer io.Writer) *Logger {
	baseLogger := logrus.New()
	lgr := &Logger{baseLogger}
	Formatter := new(logrus.JSONFormatter)
	Formatter.TimestampFormat = "2006-01-02T15:04:05"
	lgr.SetFormatter(Formatter)
	lgr.SetOutput(writer)
	return lgr
}

func (l *Logger) LogInfo(pkg string, method string, msg string) {
	l.WithFields(logrus.Fields{
		"package":  pkg,
		"function": method,
	}).Info(msg)
}

func (l *Logger) LogWarning(pkg string, method string, msg string) {
	l.WithFields(logrus.Fields{
		"package":  pkg,
		"function": method,
	}).Warn(msg)
}

func (l *Logger) LogError(pkg string, method string, err error) {
	l.WithFields(logrus.Fields{
		"package":  pkg,
		"function": method,
	}).Error(err)
}

func (l *Logger) LogAccess(r *http.Request, worktime time.Duration) {
	l.WithFields(logrus.Fields{
		"method":      r.Method,
		"remote_addr": r.RemoteAddr,
		"work_time":   worktime,
	}).Info(r.URL.Path)
}
