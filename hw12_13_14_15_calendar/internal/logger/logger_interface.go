package logger

import "go.uber.org/zap"

type Log struct {
	*zap.Logger
}

type ConfigLogger struct {
	Level            string   `json:"level"`
	OutputPaths      []string `json:"outputPaths"`
	ErrorOutputPaths []string `json:"errorOutputPaths"`
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
