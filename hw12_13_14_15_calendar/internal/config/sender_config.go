package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func NewSenderConfig(configFilePath string) (SenderConfig, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return SenderConfig{}, fmt.Errorf("%v: %w", ErrFailedOpenConfigFile, err)
	}
	defer file.Close()

	c, err := parseSenderConfigFromFile(file)
	if err != nil {
		return SenderConfig{}, err
	}

	return c, nil
}

func parseSenderConfigFromFile(f *os.File) (SenderConfig, error) {
	var config SenderConfig
	err := json.NewDecoder(f).Decode(&config)
	if err != nil {
		return SenderConfig{}, fmt.Errorf("%v: %w", ErrFailedReadConfig, err)
	}
	return config, nil
}
