package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	instance, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("error initializing zap logger: %w", err))
	}

	logger = instance.Sugar()
}

func Debug(msg ...interface{}) {
	logger.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(msg ...interface{}) {
	logger.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(msg ...interface{}) {
	logger.Warn(msg...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(msg ...interface{}) {
	logger.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(msg ...interface{}) {
	logger.Fatal(msg...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
