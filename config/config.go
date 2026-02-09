package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	ServerPort   string
	JWTSecret    string
	GoldProvider string
	WorkerCount  int
	QueueSize    int
}

var (
	configInstance *Config
	configOnce     sync.Once
)

func InitConfig() *Config {
	configOnce.Do(func() {
		configInstance = &Config{
			ServerPort:   getEnv("PORT", "8080"),
			JWTSecret:    getEnv("JWT_SECRET", "supersecretjwt"),
			GoldProvider: getEnv("GOLD_PROVIDER_URL", "http://localhost:9000"),
			WorkerCount:  getEnvAsInt("WORKER_COUNT", 5),
			QueueSize:    getEnvAsInt("QUEUE_SIZE", 100),
		}
	})
	return configInstance
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
