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
	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/pkg/logger"
	"github.com/919Umesh/stock_market_sim/pkg/queue"
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

	workerPool := queue.NewWorkerPool(cfg.WorkerCount, cfg.QueueSize)
	workerPool.Start()
	slog.Info("Worker pool started", "workers", cfg.WorkerCount)

	stockRepo := stock.NewRepository(supabaseClient)
	marketSim := market.NewSimulator(stockRepo)
	marketSim.Start(6 * time.Hour)
	slog.Info("Market simulator started", "interval", "6h")

	router := api.NewRouter(supabaseClient, cfg, workerPool)

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

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelShutdown()

	done := make(chan struct{})
	go func() {
		workerPool.Stop()
		marketSim.Stop()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("Worker pool stopped")
	case <-ctxShutdown.Done():
		slog.Warn("Shutdown timed out")
	}

	slog.Info("Server exited")
}
