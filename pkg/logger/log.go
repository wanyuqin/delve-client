package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var GlobalLogger = logrus.New()

type Fields = map[string]interface{}

func init() {
	GlobalLogger.SetOutput(os.Stdout)
	GlobalLogger.SetLevel(logrus.DebugLevel)
}

func ErrorFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Error(message)
}

func SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "info":
		GlobalLogger.SetLevel(logrus.InfoLevel)
	case "debug":
		GlobalLogger.SetLevel(logrus.DebugLevel)
	case "warn":
		GlobalLogger.SetLevel(logrus.WarnLevel)
	case "error":
		GlobalLogger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		GlobalLogger.SetLevel(logrus.FatalLevel)
	case "panic":
		GlobalLogger.SetLevel(logrus.PanicLevel)
	default:
		GlobalLogger.SetLevel(logrus.DebugLevel)
		GlobalLogger.Warnf("Log level '%s' not recognised. Setting to Debug.", level)
	}
}

func Info(message string) {
	GlobalLogger.Info(message)
}

func Infof(message string, args ...interface{}) {
	GlobalLogger.Infof(message, args...)
}

func InfoFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Info(message)
}

func Debug(message string) {
	GlobalLogger.Debug(message)
}

func Debugf(message string, args ...interface{}) {
	GlobalLogger.Debugf(message, args...)
}

func DebugFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Debug(message)
}

func Warn(message string) {
	GlobalLogger.Warn(message)
}

func Warnf(message string, args ...interface{}) {
	GlobalLogger.Warnf(message, args...)
}

func WarnFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Warn(message)
}

func Error(message string) {
	GlobalLogger.Error(message)
}

func Errorf(message string, args ...interface{}) {
	GlobalLogger.Errorf(message, args...)
}

func Fatal(message string) {
	GlobalLogger.Fatal(message)
}

func Fatalf(message string, args ...interface{}) {
	GlobalLogger.Fatalf(message, args...)
}

func FatalFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Fatal(message)
}

func Panic(message string) {
	GlobalLogger.Panic(message)
}

func Panicf(message string, args ...interface{}) {
	GlobalLogger.Panicf(message, args...)
}

func PanicFields(message string, fields Fields) {
	GlobalLogger.WithFields(map[string]interface{}(fields)).Panic(message)
}
