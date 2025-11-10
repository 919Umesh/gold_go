package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/umesh/gold_investment/api"
	"github.com/umesh/gold_investment/config"
	"github.com/umesh/gold_investment/internal/gold"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize configuration
	cfg := config.InitConfig()

	// Connect to database
	db := config.ConnectDatabase(cfg)

	// Run database migrations
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize gold price service
	goldService := gold.NewService(db, cfg)

	// Start background services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start gold price updater
	go goldService.StartPriceUpdater(ctx)

	// Initialize and start HTTP server
	router := api.NewRouter(db, cfg)

	// Graceful shutdown
	serverAddr := ":" + cfg.ServerPort
	go func() {
		log.Printf("Server starting on %s", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Server exited")
}
