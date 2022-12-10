package config

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		c := Config{
			Logger: ConfigLogger{
				Level:            "debug",
				OutputPaths:      []string{"stdout", "/tmp/logs"},
				ErrorOutputPaths: []string{"stderr"},
			},
			Server: ConfigServer{
				Host: "localhost",
				Port: "3000",
			},
			Storage: ConfigStorage{Type: "SQL", SQLDriver: "pgx", DSN: "Connection string"},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		config, err := New(tmpFilename)
		require.NoError(t, err)
		require.Equal(t, c, config)
	})

	t.Run("invalid storage type - used default", func(t *testing.T) {
		c := Config{
			Logger: ConfigLogger{
				Level:            "debug",
				OutputPaths:      []string{"stdout", "/tmp/logs"},
				ErrorOutputPaths: []string{"stderr"},
			},
			Server: ConfigServer{
				Host: "localhost",
				Port: "3000",
			},
			Storage: ConfigStorage{Type: "WRONG_TYPE", DSN: "Connection string"},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		config, err := New(tmpFilename)
		require.NoError(t, err)
		require.Equal(t, config.Logger, c.Logger)
		require.Equal(t, config.Server, c.Server)
		require.Equal(t, config.Storage.Type, StorageInMemory)
	})

	t.Run("can not open conig file", func(t *testing.T) {
		_, err := New("UNKNOWN_PATH")
		require.Error(t, err)
		require.ErrorContains(t, err, ErrFailedOpenConfigFile)
	})

	t.Run("can not read config (not json)", func(t *testing.T) {
		tmpFilename := createTmpFile(t, "")
		defer os.Remove(tmpFilename)
		_, err := New(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrFailedReadConfig)
	})

	t.Run("missing properties in json", func(t *testing.T) {
		c := "\"{\"Logger\":{\"Level\":\"debug\"},\"Server\":{},\"Storage\":{}}\""
		tmpFilename := createTmpFile(t, c)
		defer os.Remove(tmpFilename)
		_, err := New(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrFailedReadConfig)
	})

	t.Run("missing DSN for SQL storage", func(t *testing.T) {
		c := Config{
			Logger: ConfigLogger{
				Level:            "debug",
				OutputPaths:      []string{"stdout", "/tmp/logs"},
				ErrorOutputPaths: []string{"stderr"},
			},
			Server: ConfigServer{
				Host: "localhost",
				Port: "3000",
			},
			Storage: ConfigStorage{Type: StorageSql, DSN: ""},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		_, err = New(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrMissingDSN)
	})

	t.Run("missing output paths for logger", func(t *testing.T) {
		c := Config{
			Logger: ConfigLogger{
				Level:            "debug",
				OutputPaths:      []string{},
				ErrorOutputPaths: []string{"stderr"},
			},
			Server: ConfigServer{
				Host: "localhost",
				Port: "3000",
			},
			Storage: ConfigStorage{Type: StorageInMemory, DSN: ""},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		_, err = New(tmpFilename)
		require.Error(t, err)
		require.ErrorContains(t, err, ErrMissingOutputPaths)
	})

	t.Run("missing output paths for logger", func(t *testing.T) {
		c := Config{
			Logger: ConfigLogger{
				Level:            "debug",
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{},
			},
			Server: ConfigServer{
				Host: "localhost",
				Port: "3000",
			},
			Storage: ConfigStorage{Type: StorageInMemory, DSN: ""},
		}
		b, err := json.Marshal(c)
		require.NoError(t, err)
		tmpFilename := createTmpFile(t, string(b))
		defer os.Remove(tmpFilename)
		_, err = New(tmpFilename)
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
	defer f.Close()
	if _, err = f.WriteString(content); err != nil {
		log.Fatal()
	}
	return f.Name()
}
