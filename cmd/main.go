package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/919Umesh/gold_go/api"
	"github.com/919Umesh/gold_go/config"
	"github.com/919Umesh/gold_go/internal/gold"
	"github.com/919Umesh/gold_go/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {

	logger.InitLogger()

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	cfg := config.InitConfig()

	db := config.ConnectDatabase(cfg)

	goldService := gold.NewService(db, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go goldService.StartPriceUpdater(ctx)

	router := api.NewRouter(db, cfg)

	serverAddr := ":" + cfg.ServerPort
	go func() {
		slog.Info("Server starting", "address", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("Server exited")
}
