package logger

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string, outputPaths []string, errOutputPath []string) (*Log, error) {
	parsedLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, errors.New(ErrWrongLevel)
	}

	cfg := zap.Config{
		Encoding:         "json",
		Level:            parsedLevel,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errOutputPath,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("%v: %w", ErrLoggerBuild, err)
	}
	return &Log{logger}, nil
}

func (l Log) Debug(args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf("%v", args))
}

func (l Log) Info(args ...interface{}) {
	l.Logger.Info(fmt.Sprintf("%v", args))
}

func (l Log) Warn(args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf("%v", args))
}

func (l Log) Error(args ...interface{}) {
	l.Logger.Error(fmt.Sprintf("%v", args))
}

func (l Log) Fatal(args ...interface{}) {
	l.Logger.Fatal(fmt.Sprintf("%v", args))
}
