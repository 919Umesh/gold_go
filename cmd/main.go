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

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.InitConfig()

	db := config.ConnectDatabase(cfg)

	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	goldService := gold.NewService(db, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go goldService.StartPriceUpdater(ctx)

	router := api.NewRouter(db, cfg)

	serverAddr := ":" + cfg.ServerPort
	go func() {
		log.Printf("Server starting on %s", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Server exited")
}
