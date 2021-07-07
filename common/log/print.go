package log

import "fmt"

// PrintLog .
type PrintLog struct {
}

// Debug ...
func (pl *PrintLog) Debug(args ...interface{}) {
	print(args...)
}

// Debugf ...
func (pl *PrintLog) Debugf(format string, args ...interface{}) {
	printf(format, args...)
}

// Info ...
func (pl *PrintLog) Info(args ...interface{}) {
	print(args...)
}

// Infof ...
func (pl *PrintLog) Infof(format string, args ...interface{}) {
	printf(format, args...)
}

// Warn ...
func (pl *PrintLog) Warn(args ...interface{}) {
	print(args...)
}

// Warnf ...
func (pl *PrintLog) Warnf(format string, args ...interface{}) {
	printf(format, args...)
}

// Error ...
func (pl *PrintLog) Error(args ...interface{}) {
	print(args...)
}

// Errorf ...
func (pl *PrintLog) Errorf(format string, args ...interface{}) {
	printf(format, args...)
}

func print(args ...interface{}) {
	fmt.Println(args...)
}

func printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
