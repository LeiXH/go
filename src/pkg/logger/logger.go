package logger

import (
	"go.uber.org/zap"
)

var _logger *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	_logger = logger.Sugar()
}

func Debug(msg string) {
	_logger.Debug(msg)
}
func Debugf(format string, args ...interface{}) {
	_logger.Debugf(format, args...)
}
func Debugw(msg string, kv ...interface{}) {
	_logger.Debugw(msg, kv...)
}

func Info(msg string) {
	_logger.Info(msg)
}
func Infof(format string, args ...interface{}) {
	_logger.Infof(format, args...)
}
func Infow(msg string, kv ...interface{}) {
	_logger.Infow(msg, kv...)
}

func Warn(msg string) {
	_logger.Warn(msg)
}
func Warnf(format string, args ...interface{}) {
	_logger.Warnf(format, args...)
}
func Warnw(msg string, kv ...interface{}) {
	_logger.Warnw(msg, kv...)
}

func Error(msg string) {
	_logger.Error(msg)
}
func Errorf(format string, args ...interface{}) {
	_logger.Errorf(format, args...)
}
func Errorw(msg string, kv ...interface{}) {
	_logger.Errorw(msg, kv...)
}

func Fatal(msg string) {
	_logger.Fatal(msg)
}
func Fatalf(format string, args ...interface{}) {
	_logger.Fatalf(format, args...)
}
func Fatalw(msg string, kv ...interface{}) {
	_logger.Fatalw(msg, kv...)
}
