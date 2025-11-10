// package main

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"github.com/umesh/gold_investment/api"
// 	"github.com/umesh/gold_investment/config"
// 	"github.com/umesh/gold_investment/internal/gold"
// )

// func main() {

// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using environment variables")
// 	}

// 	cfg := config.InitConfig()

// 	db := config.ConnectDatabase(cfg)

// 	if err := config.AutoMigrate(db); err != nil {
// 		log.Fatalf("Failed to run migrations: %v", err)
// 	}

// 	goldService := gold.NewService(db, cfg)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	go goldService.StartPriceUpdater(ctx)

// 	router := api.NewRouter(db, cfg)

// 	serverAddr := ":" + cfg.ServerPort
// 	go func() {
// 		log.Printf("Server starting on %s", serverAddr)
// 		if err := router.Run(serverAddr); err != nil {
// 			log.Fatalf("Failed to start server: %v", err)
// 		}
// 	}()

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	log.Println("Shutting down server...")

// 	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	log.Println("Server exited")
// }

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	log.Println("✅ Database connection established")

	// Run database migrations
	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("❌ Database migrations failed: %v", err)
	}

	// Initialize gold price service
	goldService := gold.NewService(db, cfg)

	// Start background services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start gold price updater
	go goldService.StartPriceUpdater(ctx)
	log.Println("✅ Gold price updater started")

	// Initialize and start HTTP server
	router := api.NewRouter(db, cfg)

	// Graceful shutdown
	serverAddr := ":" + cfg.ServerPort
	if serverAddr == ":" {
		serverAddr = ":8080"
	}

	go func() {
		log.Printf("🌐 Server starting on %s", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("✅ Server exited gracefully")
}
