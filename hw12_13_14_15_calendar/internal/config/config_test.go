package config

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCalendarConfig(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		c := CalendarConfig{
			Logger: Logger{
				Level:            "debug",
				OutputPaths:      []string{"stdout", "/tmp/logs"},
				ErrorOutputPaths: []string{"stderr"},
			},
			Server:  Server{Host: "localhost", HTTPPort: "3000", GrpcPort: "3001"},
			Storage: Storage{Type: "sql", Driver: "postgres", DSN: "Connection string"},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		config, err := NewCalendarConfig(tmpFilename)
		require.NoError(t, err)
		require.Equal(t, c, config)
	})

	t.Run("invalid storage type - used default", func(t *testing.T) {
		c := CalendarConfig{
			Logger: Logger{
				Level:            "debug",
				OutputPaths:      []string{"stdout", "/tmp/logs"},
				ErrorOutputPaths: []string{"stderr"},
			},
			Server:  Server{Host: "localhost", HTTPPort: "3000", GrpcPort: "3001"},
			Storage: Storage{Type: "WRONG_TYPE", DSN: "Connection string"},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		config, err := NewCalendarConfig(tmpFilename)
		require.NoError(t, err)
		require.Equal(t, config.Logger, c.Logger)
		require.Equal(t, config.Server, c.Server)
		require.Equal(t, config.Storage.Type, StorageInMemory)
	})

	t.Run("can not open conig file", func(t *testing.T) {
		_, err := NewCalendarConfig("UNKNOWN_PATH")
		require.Error(t, err)
		require.ErrorContains(t, err, ErrFailedOpenConfigFile)
	})

	t.Run("can not read config (not json)", func(t *testing.T) {
		tmpFilename := createTmpFile(t, "")
		defer os.Remove(tmpFilename)
		_, err := NewCalendarConfig(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrFailedReadConfig)
	})

	t.Run("missing properties in json", func(t *testing.T) {
		c := "\"{\"Logger\":{\"Level\":\"debug\"},\"Server\":{},\"Storage\":{}}\""
		tmpFilename := createTmpFile(t, c)
		defer os.Remove(tmpFilename)
		_, err := NewCalendarConfig(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrFailedReadConfig)
	})

	t.Run("missing DSN for SQL storage", func(t *testing.T) {
		c := CalendarConfig{
			Logger: Logger{
				Level:            "debug",
				OutputPaths:      []string{"stdout", "/tmp/logs"},
				ErrorOutputPaths: []string{"stderr"},
			},
			Server:  Server{Host: "localhost", HTTPPort: "3000", GrpcPort: "3001"},
			Storage: Storage{Type: StorageSQL, DSN: ""},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		_, err = NewCalendarConfig(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrMissingDSN)
	})

	t.Run("missing output paths for logger", func(t *testing.T) {
		c := CalendarConfig{
			Logger: Logger{
				Level:            "debug",
				OutputPaths:      []string{},
				ErrorOutputPaths: []string{"stderr"},
			},
			Server:  Server{Host: "localhost", HTTPPort: "3000", GrpcPort: "3001"},
			Storage: Storage{Type: StorageInMemory, DSN: ""},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		_, err = NewCalendarConfig(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrMissingOutputPaths)
	})

	t.Run("missing output paths for logger", func(t *testing.T) {
		c := CalendarConfig{
			Logger: Logger{
				Level:            "debug",
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{},
			},
			Server:  Server{Host: "localhost", HTTPPort: "3000", GrpcPort: "3001"},
			Storage: Storage{Type: StorageInMemory, DSN: ""},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		_, err = NewCalendarConfig(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrMissingErrOutputPaths)
	})
}

func createTmpFile(t *testing.T, content string) (filepath string) {
	t.Helper()
	f, err := os.CreateTemp("", "test-config")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = f.WriteString(content); err != nil {
		log.Fatal()
	}
	return f.Name()
}
