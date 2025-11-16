package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if config.IntervalHours == 0 {
		config.IntervalHours = 6
	}

	if config.Language == "" {
		config.Language = "ja"
	}

	return &config, nil
}

