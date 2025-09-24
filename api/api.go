package api

import (
	"fmt"
	"log/slog"
	"os"
)

type CliConfig struct {
	// shared parameters
	LogLevel string `mapstructure:"log-level"`
}

func ConfigureLogging(level string) error {
	var logLevel slog.Level

	switch level {
	case "error":
		logLevel = slog.LevelError
	case "warn":
		logLevel = slog.LevelWarn
	case "info":
		logLevel = slog.LevelInfo
	case "debug":
		logLevel = slog.LevelDebug
	default:
		return fmt.Errorf("unknown log-level: %s", level)
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return nil
}

func (c *CliConfig) GetTimes() error {
	if err := ConfigureLogging(c.LogLevel); err != nil {
		return err
	}

	slog.Info("GetTimes() called")

	return nil
}

func (c *CliConfig) SetTimes() error {
	if err := ConfigureLogging(c.LogLevel); err != nil {
		return err
	}

	slog.Info("SetTimes() called")

	return nil
}
