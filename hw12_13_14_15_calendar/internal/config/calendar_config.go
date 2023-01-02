package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func NewCalendarConfig(configFilePath string) (CalendarConfig, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return CalendarConfig{}, fmt.Errorf("%v: %w", ErrFailedOpenConfigFile, err)
	}
	defer file.Close()

	c, err := parseCalendarConfigFromFile(file)
	if err != nil {
		return CalendarConfig{}, err
	}

	err = validateCalendarConfig(c)
	if err != nil {
		return CalendarConfig{}, err
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

func parseCalendarConfigFromFile(f *os.File) (CalendarConfig, error) {
	var config CalendarConfig
	err := json.NewDecoder(f).Decode(&config)
	if err != nil {
		return CalendarConfig{}, fmt.Errorf("%v: %w", ErrFailedReadConfig, err)
	}
	config.Storage.Type = getStorageType(config.Storage.Type)
	return config, nil
}

func validateCalendarConfig(c CalendarConfig) error {
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
