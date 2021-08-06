package log

import (
	fmtlog "baas/common/log/fmt"
)

// Logger 日志
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

var std Logger

func init() {
	std = &fmtlog.PrintLog{}
}

// ResetLogger 重设默认 logger
func ResetLogger(l Logger) {
	std = l
}

// GetLogger 返回默认 logger
func GetLogger() Logger {
	return std
}

// Debug ...
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Info ...
func Info(args ...interface{}) {
	std.Info(args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warn ...
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Warnf ...
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Error ...
func Error(args ...interface{}) {
	std.Error(args...)
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}
