package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func NewSchedulerConfig(configFilePath string) (SchedulerConfig, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return SchedulerConfig{}, fmt.Errorf("%v: %w", ErrFailedOpenConfigFile, err)
	}
	defer file.Close()

	c, err := parseSchedulerConfigFromFile(file)
	if err != nil {
		return SchedulerConfig{}, err
	}

	return c, nil
}

func parseSchedulerConfigFromFile(f *os.File) (SchedulerConfig, error) {
	var config SchedulerConfig
	err := json.NewDecoder(f).Decode(&config)
	if err != nil {
		return SchedulerConfig{}, fmt.Errorf("%v: %w", ErrFailedReadConfig, err)
	}
	return config, nil
}
