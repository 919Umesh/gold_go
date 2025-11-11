// package config

// import (
// 	"os"
// 	"strconv"
// 	"sync"
// )

// type Config struct {
// 	DBHost       string
// 	DBUser       string
// 	DBPassword   string
// 	DBName       string
// 	DBPort       string
// 	ServerPort   string
// 	JWTSecret    string
// 	GoldProvider string
// 	WorkerCount  int
// 	QueueSize    int
// }

// var (
// 	configInstance *Config
// 	configOnce     sync.Once
// )

// func InitConfig() *Config {
// 	configOnce.Do(func() {
// 		configInstance = &Config{
// 			DBHost:       getEnv("DB_HOST", "localhost"),
// 			DBUser:       getEnv("DB_USER", "postgres"),
// 			DBPassword:   getEnv("DB_PASSWORD", "postgres"),
// 			DBName:       getEnv("DB_NAME", "gold_invest"),
// 			DBPort:       getEnv("DB_PORT", "5432"),
// 			ServerPort:   getEnv("PORT", "8080"),
// 			JWTSecret:    getEnv("JWT_SECRET", "supersecretjwt"),
// 			GoldProvider: getEnv("GOLD_PROVIDER_URL", "http://localhost:9000"),
// 			WorkerCount:  getEnvAsInt("WORKER_COUNT", 5),
// 			QueueSize:    getEnvAsInt("QUEUE_SIZE", 100),
// 		}
// 	})
// 	return configInstance
// }

// func getEnv(key, defaultValue string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}
// 	return defaultValue
// }

// func getEnvAsInt(key string, defaultValue int) int {
// 	if value := os.Getenv(key); value != "" {
// 		if intValue, err := strconv.Atoi(value); err == nil {
// 			return intValue
// 		}
// 	}
// 	return defaultValue
// }

package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	DBHost       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBPort       string
	ServerPort   string
	JWTSecret    string
	GoldProvider string
	WorkerCount  int
	QueueSize    int
	DatabaseURL  string // Add this for Railway
}

var (
	configInstance *Config
	configOnce     sync.Once
)

func InitConfig() *Config {
	configOnce.Do(func() {
		configInstance = &Config{
			DBHost:       getEnv("DB_HOST", "localhost"),
			DBUser:       getEnv("DB_USER", "postgres"),
			DBPassword:   getEnv("DB_PASSWORD", "password"), // Changed from "postgres"
			DBName:       getEnv("DB_NAME", "gold_investment_db"),
			DBPort:       getEnv("DB_PORT", "5432"),
			ServerPort:   getEnv("PORT", "8080"),
			JWTSecret:    getEnv("JWT_SECRET", "dev_jwt_secret_min_32_chars_long_here"), // Better default
			GoldProvider: getEnv("GOLD_PROVIDER_URL", ""),                               // Empty default
			WorkerCount:  getEnvAsInt("WORKER_COUNT", 5),
			QueueSize:    getEnvAsInt("QUEUE_SIZE", 100),
			DatabaseURL:  getEnv("DATABASE_URL", ""),
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
