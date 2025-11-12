package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
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
			DBHost:        getEnv("DB_HOST", "localhost"),
			DBUser:        getEnv("DB_USER", "postgres"),
			DBPassword:    getEnv("DB_PASSWORD", "postgres"),
			DBName:        getEnv("DB_NAME", "gold_invest"),
			DBPort:        getEnv("DB_PORT", "5432"),
			ServerPort:    getEnv("PORT", "8080"),
			JWTSecret:     getEnv("JWT_SECRET", "supersecretjwt"),
			GoldProvider:  getEnv("GOLD_PROVIDER_URL", "http://localhost:9000"),
			WorkerCount:   getEnvAsInt("WORKER_COUNT", 5),
			QueueSize:     getEnvAsInt("QUEUE_SIZE", 100),
			RedisAddress:  getEnv("REDIS_ADDRESS", "localhost:6379"),
			RedisPassword: getEnv("REDIS_PASSWORD", ""),
			RedisDB:       getEnvAsInt("REDIS_DB", 0),
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
