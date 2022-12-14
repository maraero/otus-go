package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func New(configFilePath string) (Config, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("%v: %w", ErrFailedOpenConfigFile, err)
	}
	defer file.Close()

	c, err := parseConfigFromFile(file)
	if err != nil {
		return Config{}, err
	}

	err = validateConfig(c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

// getStorageType returns valid storage type. It is "in memory" by default.
func getStorageType(t string) string {
	if t == StorageInMemory || t == StorageSQL {
		return t
	}
	return StorageInMemory
}

func parseConfigFromFile(f *os.File) (Config, error) {
	var config Config
	err := json.NewDecoder(f).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("%v: %w", ErrFailedReadConfig, err)
	}
	config.Storage.Type = getStorageType(config.Storage.Type)
	return config, nil
}

func validateConfig(c Config) error {
	err := validateConfigLogger(c.Logger)
	if err != nil {
		return err
	}

	err = validateConfigStorage(c.Storage)
	if err != nil {
		return err
	}

	return nil
}

func validateConfigLogger(c Logger) error {
	if len(c.OutputPaths) == 0 {
		return errors.New(ErrMissingOutputPaths)
	}

	if len(c.ErrorOutputPaths) == 0 {
		return errors.New(ErrMissingErrOutputPaths)
	}

	return nil
}

func validateConfigStorage(c Storage) error {
	if c.Type == StorageSQL {
		return validateSQLConfig(c.DSN, c.Driver)
	}
	return nil
}

func validateSQLConfig(dsn string, database string) error {
	if dsn == "" {
		return errors.New(ErrMissingDSN)
	}

	for _, d := range AllowedDatabases {
		if d == database {
			return nil
		}
	}

	return errors.New(ErrWrongSQLDriver)
}
