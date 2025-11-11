package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/umesh/gold_investment/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

type DB struct {
	*gorm.DB
}

func ConnectDatabase(cfg *Config) *gorm.DB {
	dbOnce.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
		)

		var err error
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: true,
		})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		sqlDB, err := dbInstance.DB()
		if err != nil {
			log.Fatalf("Failed to get database instance: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		log.Println("Database connected successfully")
	})
	return dbInstance
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
		&models.GoldPrice{},
	)
}
