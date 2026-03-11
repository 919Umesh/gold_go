package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/919Umesh/stock_market_sim/api"
	"github.com/919Umesh/stock_market_sim/config"
	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	cfg := config.InitConfig()

	supabaseClient, err := supabase.NewClient()
	if err != nil {
		slog.Error("Failed to initialize Supabase client", "error", err)
		os.Exit(1)
	}

	router := api.NewRouter(supabaseClient, cfg)

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

	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("Server exited")
}
