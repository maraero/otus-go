package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func NewConfig(configFilePath string) (Config, error) {
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

// getStorageType returns valid storage type. It is "in memory" by default
func getStorageType(t string) string {
	if t == StorageInMemory || t == StorageSql {
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
	if c.Storage.Type == StorageSql && c.Storage.DSN == "" {
		return errors.New(ErrMissingDSN)
	}
	return nil
}
