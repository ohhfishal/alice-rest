package config

import (
	"log/slog"
)

type LogLevel string

func (logLevel LogLevel) Level() slog.Level {
	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}

}

type Config struct {
	Host     string
	Port     string
	LogLevel slog.Level
}

func NewConfig(_ []string, getenv func(string) string) *Config {
	logLevel := (LogLevel)(getenv("LOG_LEVEL"))
	config := &Config{
		Host:     getenv("HOST"),
		Port:     getenv("PORT"),
		LogLevel: logLevel.Level(),
	}
	config.UseDefaults()
	return config
}

func (c *Config) UseDefaults() {
	if c.Host == "" {
		c.Host = ""
	}
	if c.Port == "" {
		c.Port = "8000"
	}
}
