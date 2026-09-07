package config

import (
	"os"
	"time"
)

type Config struct {
	ServerAddress  string
	ServerPort     string
	DatabaseURL    string
	LogLevel       string
	RequestTimeout time.Duration
	MigrationsPath string
}

func Load() *Config {
	return &Config{
		ServerAddress:  getEnv("SERVER_ADDRESS", "0.0.0.0"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/subscriptions?sslmode=disable"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		RequestTimeout: getDurationEnv("REQUEST_TIMEOUT", 30*time.Second),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "./migrations"),
	}
}

func (c *Config) Address() string {
	return c.ServerAddress + ":" + c.ServerPort
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
