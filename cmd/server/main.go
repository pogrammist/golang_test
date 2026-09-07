package main

import (
	"github.com/pogrammist/golang_test/internal/config"
	"github.com/pogrammist/golang_test/internal/model"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger := setupLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("starting subscription service", zap.String("address", cfg.Address()))

	db, err := model.NewDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()
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
