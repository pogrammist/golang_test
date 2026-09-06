package main

import (
	"github.com/pogrammist/golang_test/internal/config"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger := setupLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("starting subscription service", zap.String("address", cfg.Address()))
}

func setupLogger(level string) *zap.Logger {
	switch level {
	case "debug":
		logger, _ := zap.NewDevelopment()
		return logger
	default:
		logger, _ := zap.NewProduction()
		return logger
	}
}
