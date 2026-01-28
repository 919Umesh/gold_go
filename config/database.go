package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func ConnectDatabase(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		// We shouldn't exit here strictly speaking, but for simplicity of migration we panic or let the caller handle nil.
		// Ideally we return (db, error). But since signature is fixed to *gorm.DB, we log Fatal-like behaviour or return nil?
		// The original code Fatalf'ed. Let's keep strict behavior but use idiomatic log if possible,
		// but since we return *gorm.DB, returning nil might crash the app later.
		// Let's stick to panic/Fatal if we can't change signature, OR change signature.
		// Changing signature breaks main.go.
		// Let's use os.Exit(1) to mimic log.Fatalf behavior cleanly.
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get database instance", "error", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("Database connected successfully")

	return db
}
