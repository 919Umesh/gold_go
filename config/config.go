// package config

// import (
// 	"os"
// 	"strconv"
// 	"sync"
// )

// type Config struct {
// 	DBHost        string
// 	DBUser        string
// 	DBPassword    string
// 	DBName        string
// 	DBPort        string
// 	ServerPort    string
// 	JWTSecret     string
// 	GoldProvider  string
// 	WorkerCount   int
// 	QueueSize     int
// 	RedisAddress  string
// 	RedisPassword string
// 	RedisDB       int
// }

// var (
// 	configInstance *Config
// 	configOnce     sync.Once
// )

// func InitConfig() *Config {
// 	configOnce.Do(func() {
// 		configInstance = &Config{
// 			DBHost:        getEnv("DB_HOST", "localhost"),
// 			DBUser:        getEnv("DB_USER", "postgres"),
// 			DBPassword:    getEnv("DB_PASSWORD", "postgres"),
// 			DBName:        getEnv("DB_NAME", "gold_invest"),
// 			DBPort:        getEnv("DB_PORT", "5432"),
// 			ServerPort:    getEnv("PORT", "8080"),
// 			JWTSecret:     getEnv("JWT_SECRET", "supersecretjwt"),
// 			GoldProvider:  getEnv("GOLD_PROVIDER_URL", "http://localhost:9000"),
// 			WorkerCount:   getEnvAsInt("WORKER_COUNT", 5),
// 			QueueSize:     getEnvAsInt("QUEUE_SIZE", 100),
// 			RedisAddress:  getEnv("REDIS_ADDRESS", "localhost:6379"),
// 			RedisPassword: getEnv("REDIS_PASSWORD", ""),
// 			RedisDB:       getEnvAsInt("REDIS_DB", 0),
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
	"net/url"
	"os"
	"strconv"
	"strings"
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
		dbHost, dbUser, dbPass, dbName, dbPort := parseDatabaseURL(os.Getenv("DATABASE_URL"))
		redisAddr, redisPass := parseRedisURL(os.Getenv("REDIS_URL"))

		configInstance = &Config{
			DBHost:        dbHost,
			DBUser:        dbUser,
			DBPassword:    dbPass,
			DBName:        dbName,
			DBPort:        dbPort,
			ServerPort:    getEnv("PORT", "8080"),
			JWTSecret:     getEnv("JWT_SECRET", "supersecretjwt"),
			GoldProvider:  getEnv("GOLD_PROVIDER_URL", "http://localhost:9000"),
			WorkerCount:   getEnvAsInt("WORKER_COUNT", 5),
			QueueSize:     getEnvAsInt("QUEUE_SIZE", 100),
			RedisAddress:  redisAddr,
			RedisPassword: redisPass,
			RedisDB:       getEnvAsInt("REDIS_DB", 0),
		}
	})
	return configInstance
}

// --- helper functions ---

func parseDatabaseURL(dbURL string) (host, user, password, dbname, port string) {
	if dbURL == "" {
		return "localhost", "postgres", "postgres", "gold_invest", "5432"
	}

	u, err := url.Parse(dbURL)
	if err != nil {
		return
	}

	userInfo := u.User
	user = userInfo.Username()
	password, _ = userInfo.Password()

	hostPort := strings.Split(u.Host, ":")
	host = hostPort[0]
	if len(hostPort) > 1 {
		port = hostPort[1]
	} else {
		port = "5432"
	}

	dbname = strings.TrimPrefix(u.Path, "/")
	return
}

func parseRedisURL(redisURL string) (addr, password string) {
	if redisURL == "" {
		return "localhost:6379", ""
	}
	u, err := url.Parse(redisURL)
	if err != nil {
		return "localhost:6379", ""
	}

	host := u.Host
	pass, _ := u.User.Password()
	return host, pass
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
