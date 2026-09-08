package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pogrammist/golang_test/internal/config"
	"github.com/pogrammist/golang_test/internal/handler"
	"github.com/pogrammist/golang_test/internal/model"
	"github.com/pogrammist/golang_test/internal/repository"
	"github.com/pogrammist/golang_test/internal/service"
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

	repo := repository.NewSubscriptionRepository(db)
	if err := repo.RunMigrations(cfg.MigrationsPath, cfg.DatabaseURL); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}
	logger.Info("database migrations completed")

	subscriptionService := service.NewSubscriptionService(repo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, logger)

	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", subscriptionHandler.HealthCheck)

	v1 := router.Group("/")
	{
		v1.POST("/subscriptions", subscriptionHandler.CreateSubscription)
		v1.GET("/subscriptions", subscriptionHandler.ListSubscriptions)
		v1.GET("/subscriptions/:id", subscriptionHandler.GetSubscription)
		v1.PUT("/subscriptions/:id", subscriptionHandler.UpdateSubscription)
		v1.DELETE("/subscriptions/:id", subscriptionHandler.DeleteSubscription)
		v1.GET("/subscriptions/total", subscriptionHandler.CalculateTotalCost)
	}

	logger.Info("server started", zap.String("address", cfg.Address()))
	if err := router.Run(cfg.Address()); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
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
