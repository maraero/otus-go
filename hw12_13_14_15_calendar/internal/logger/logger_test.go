package logger

import (
	"log"
	"os"
	"testing"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("builds logger", func(t *testing.T) {
		logger, err := New(config.Logger{
			Level:            "debug",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		})
		require.NoError(t, err)
		require.NotNil(t, logger)
	})

	t.Run("wrong level error", func(t *testing.T) {
		_, err := New(config.Logger{
			Level:            "WRONG_LEVEL",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		})
		require.Error(t, err)
		require.ErrorContains(t, err, ErrWrongLevel)
	})

	t.Run("logs info with info level", func(t *testing.T) {
		tmpFilename := createTmpFile(t)
		defer os.Remove(tmpFilename)

		logger, err := New(config.Logger{
			Level:            "info",
			OutputPaths:      []string{tmpFilename},
			ErrorOutputPaths: []string{"stderr"},
		})
		require.NoError(t, err)
		logText := "test log"
		logger.Info(logText)
		logContent, err := os.ReadFile(tmpFilename)
		require.NoError(t, err)
		require.Contains(t, string(logContent), logText)
	})

	t.Run("logs error with info level", func(t *testing.T) {
		tmpFilename := createTmpFile(t)
		defer os.Remove(tmpFilename)

		logger, err := New(config.Logger{
			Level:            "info",
			OutputPaths:      []string{tmpFilename},
			ErrorOutputPaths: []string{"stderr"},
		})
		require.NoError(t, err)
		logText := "test err log"
		logger.Error(logText)
		logContent, err := os.ReadFile(tmpFilename)
		require.NoError(t, err)
		require.Contains(t, string(logContent), logText)
	})

	t.Run("does not logs info with warn level", func(t *testing.T) {
		tmpFilename := createTmpFile(t)
		defer os.Remove(tmpFilename)

		logger, err := New(config.Logger{
			Level:            "warn",
			OutputPaths:      []string{tmpFilename},
			ErrorOutputPaths: []string{"stderr"},
		})
		require.NoError(t, err)
		logText := "test err log"
		logger.Info(logText)
		logContent, err := os.ReadFile(tmpFilename)
		require.NoError(t, err)
		require.Contains(t, string(logContent), "")
	})
}

func createTmpFile(t *testing.T) (filepath string) {
	t.Helper()
	f, err := os.CreateTemp("", "test-config")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	return f.Name()
}
