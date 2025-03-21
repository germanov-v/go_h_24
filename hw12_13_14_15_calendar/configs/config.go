package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerConfig `json:"Server"`
	DbBaseConfig `json:"Database"`
	LogConfig    `json:"Logs"`
}
type ServerConfig struct {
	Host string `json:"Host"`
	Port int    `json:"Port"`
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

type LogConfig struct {
	Level     LogLevel    `json:"Level"`
	Providers LogProvider `json:"Providers"`
}

type LogProvider struct {
	File string `json:"File"`
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

	switch config.LogConfig.Level {
	case Info, Debug, Error, Warn, Trace:
	default:
		return nil, fmt.Errorf("unknown config level: %s", config.LogConfig.Level)
	}

	// из стека в кучу
	return &config, nil
}
