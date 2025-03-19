package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbBaseConfig `json:"Database"`
	LogLevel     `json:"LogLevel"`
}

type DbBaseConfig struct {
	ConnectionString string `json:"ConnectionString"`
}

type LogLevel string

const (
	Info  LogLevel = "INFO"
	Debug LogLevel = "DEBUG"
	Error LogLevel = "ERROR"
	Warn  LogLevel = "WARN"
	Trace LogLevel = "TRACE"
)

type LogLevelsConfig struct {
	level LogLevel
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// из стека в кучу
	return &config, nil
}
