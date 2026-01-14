package helper

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

// SetLogLevel sets the logging level
func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

// SetLogFormatter sets the log formatter
func SetLogFormatter(formatter logrus.Formatter) {
	log.SetFormatter(formatter)
}

// LogRequest logs HTTP request/response information
func LogRequest(source, api string, payload, responseBody interface{}, statusCode int) {
	_, fn, line, _ := runtime.Caller(1)
	log.WithFields(logrus.Fields{
		"filename":    fn,
		"line":        line,
		"url":         api,
		"payload":     payload,
		"response":    responseBody,
		"status_code": statusCode,
	}).Info(source)
}

// LogError logs error information
func LogError(source string, err error, api string, payload interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	log.WithFields(logrus.Fields{
		"filename": fn,
		"line":     line,
		"url":      api,
		"payload":  payload,
		"error":    err.Error(),
	}).Error(source)
}

// LogInfo logs general information
func LogInfo(source string, fields map[string]interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	logFields := logrus.Fields{
		"filename": fn,
		"line":     line,
	}
	for k, v := range fields {
		logFields[k] = v
	}
	log.WithFields(logFields).Info(source)
}

// LogDebug logs debug information
func LogDebug(source string, fields map[string]interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	logFields := logrus.Fields{
		"filename": fn,
		"line":     line,
	}
	for k, v := range fields {
		logFields[k] = v
	}
	log.WithFields(logFields).Debug(source)
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	return log
}
