// Package logger provides structured logging using zap.
package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

func NewLogger(level string, env string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	if env == "dev" {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.Level = zap.NewAtomicLevelAt(parseLogLevel(level))

	return cfg.Build()
}

func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}