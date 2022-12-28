package logger

import "go.uber.org/zap"

type Log struct {
	*zap.Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
