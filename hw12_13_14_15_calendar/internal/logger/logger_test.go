package logger

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("builds logger", func(t *testing.T) {
		logger, err := New("debug", []string{"stdout"}, []string{"stderr"})
		require.NoError(t, err)
		require.NotNil(t, logger)
	})

	t.Run("wrong level error", func(t *testing.T) {
		_, err := New("WRONG_LEVEL", []string{"stdout"}, []string{"stderr"})
		require.Error(t, err)
		require.ErrorContains(t, err, ErrWrongLevel)
	})

	t.Run("logs info with info level", func(t *testing.T) {
		tmpFilename := createTmpFile(t)
		defer os.Remove(tmpFilename)

		logger, err := New("info", []string{tmpFilename}, []string{"stderr"})
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

		logger, err := New("info", []string{tmpFilename}, []string{"stderr"})
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

		logger, err := New("warn", []string{tmpFilename}, []string{"stderr"})
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
