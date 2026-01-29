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
	"github.com/919Umesh/gold_go/pkg/queue"
	"github.com/joho/godotenv"
)

func main() {

	logger.InitLogger()

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	cfg := config.InitConfig()

	db := config.ConnectDatabase(cfg)

	// Initialize Worker Pool
	workerPool := queue.NewWorkerPool(cfg.WorkerCount, cfg.QueueSize)
	workerPool.Start()
	slog.Info("Worker pool started", "workers", cfg.WorkerCount)

	// Initialize Gold Service with Worker Pool
	goldService := gold.NewService(db, cfg, workerPool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Price Updater (Runs in its own goroutine)
	go goldService.StartPriceUpdater(ctx)

	// Initialize Router with injected dependencies
	router := api.NewRouter(db, cfg, goldService, workerPool)

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

	// Graceful shutdown
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelShutdown()

	// Stop cancellation context for services
	cancel()

	done := make(chan struct{})
	go func() {
		// Stop Worker Pool and wait for jobs to finish
		workerPool.Stop()
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
