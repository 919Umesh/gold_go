package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	ServerPort    string
	JWTSecret     string
	GoldProvider  string
	WorkerCount   int
	QueueSize     int
	RedisAddress  string
	RedisPassword string
	RedisDB       int
}

var (
	configInstance *Config
	configOnce     sync.Once
)

func InitConfig() *Config {
	configOnce.Do(func() {
		configInstance = &Config{
			ServerPort:    getEnv("PORT", "8080"),
			JWTSecret:     getEnv("JWT_SECRET", "supersecretjwt"),
			GoldProvider:  getEnv("GOLD_PROVIDER_URL", "http://localhost:9000"),
			WorkerCount:   getEnvAsInt("WORKER_COUNT", 5),
			QueueSize:     getEnvAsInt("QUEUE_SIZE", 100),
			RedisAddress:  getRedisAddress(),
			RedisPassword: getEnv("REDIS_PASSWORD", ""),
			RedisDB:       getEnvAsInt("REDIS_DB", 0),
		}
	})
	return configInstance
}

func getRedisAddress() string {
	if url := os.Getenv("REDIS_URL"); url != "" {
		return url
	}
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	if host != "" && port != "" {
		return host + ":" + port
	}
	return getEnv("REDIS_ADDRESS", "localhost:6379")
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
