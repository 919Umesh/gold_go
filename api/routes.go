package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"log/slog"

	"github.com/919Umesh/stock_market_sim/config"
	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/919Umesh/stock_market_sim/internal/auth"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/internal/trading"
	"github.com/919Umesh/stock_market_sim/internal/wallet"
	"github.com/919Umesh/stock_market_sim/pkg/middleware"
	"github.com/919Umesh/stock_market_sim/pkg/queue"
	"github.com/919Umesh/stock_market_sim/pkg/redis"
	"github.com/gin-contrib/cors"
)

type Router struct {
	client      *appwrite.Client
	cfg         *config.Config
	engine      *gin.Engine
	redisClient *redis.Client
	workerPool  *queue.WorkerPool
}

func NewRouter(client *appwrite.Client, cfg *config.Config, wp *queue.WorkerPool) *Router {
	router := &Router{
		client:     client,
		cfg:        cfg,
		engine:     gin.Default(),
		workerPool: wp,
	}

	//To allow the all origin for development purpose
	router.engine.Use(cors.Default())

	router.redisClient = redis.NewRedisClient(
		cfg.RedisAddress,
		cfg.RedisPassword,
		cfg.RedisDB,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := router.redisClient.Ping(ctx); err != nil {
		slog.Warn("Failed to connect to redis. Caching and Rate-limiting will be disabled or fail.", "error", err)
		slog.Info("To start redis locally, run: docker-compose up -d redis")
	}

	router.setupRoutes()
	return router
}

func (r *Router) setupRoutes() {
	rateLimiter := middleware.NewRateLimiter(r.redisClient)
	cacheMiddleware := middleware.NewCacheMiddleware(r.redisClient)

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	v1 := r.engine.Group("/api/v1")
	{

		public := v1.Group("")
		{
			authRepo := auth.NewRepository(r.client)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			public.POST("/auth/register", rateLimiter.RateLimit(), authHandler.Register)
			public.POST("/auth/login", rateLimiter.RateLimit(), authHandler.Login)

			// Stock Market Public Routes - Commented out until Stock Service is refactored
			// Stock Market Public Routes
			stockRepo := stock.NewRepository(r.client)
			stockService := stock.NewService(stockRepo)
			stockHandler := NewStockHandler(stockService)

			public.GET("/stocks", rateLimiter.RateLimit(), cacheMiddleware.Cache(2*time.Minute), stockHandler.ListCompanies)
			public.GET("/stocks/search", rateLimiter.RateLimit(), stockHandler.SearchCompanies)
			public.GET("/stocks/:symbol", rateLimiter.RateLimit(), stockHandler.GetCompany)
			public.GET("/stocks/:symbol/price", rateLimiter.RateLimit(), stockHandler.GetCurrentPrice)
			public.GET("/stocks/:symbol/history", rateLimiter.RateLimit(), stockHandler.GetPriceHistory)
			public.GET("/stocks/market-overview", rateLimiter.RateLimit(), stockHandler.GetMarketOverview)
			public.GET("/stocks/top-gainers", rateLimiter.RateLimit(), stockHandler.GetTopGainers)
			public.GET("/stocks/top-losers", rateLimiter.RateLimit(), stockHandler.GetTopLosers)
			public.GET("/stocks/most-active", rateLimiter.RateLimit(), stockHandler.GetMostActive)
			public.GET("/stocks/:symbol/events", rateLimiter.RateLimit(), stockHandler.GetUpcomingEvents)
		}

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(r.cfg))
		{
			authRepo := auth.NewRepository(r.client)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			protected.GET("/auth/profile", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), authHandler.GetProfile)
			protected.PUT("/auth/profile/update", rateLimiter.RateLimit(), authHandler.UpdateProfile)

			walletRepo := wallet.NewRepository(r.client)
			walletService := wallet.NewService(walletRepo, r.workerPool)
			walletHandler := wallet.NewHandler(walletService)

			protected.GET("/wallet", rateLimiter.RateLimit(), walletHandler.GetWallet)
			protected.GET("/transaction", rateLimiter.RateLimit(), walletHandler.GetUserTransaction)
			protected.POST("/wallet/topup", rateLimiter.RateLimit(), walletHandler.TopUp)

			// Stock Trading Protected Routes
			// Re-initialize stock components for protected group if needed, or reuse variable if scope allows.
			// Scoping in Go blocks: 'stockRepo' above is in 'public' block. Here we are in 'protected'.
			// We should probably lift repo initialization to 'v1' scope or re-initialize.
			// Since r.client is stateless/shared, re-init is cheap.
			stockRepoProtected := stock.NewRepository(r.client)
			tradingRepo := trading.NewRepository(r.client)
			tradingService := trading.NewService(tradingRepo, stockRepoProtected)
			tradingHandler := NewTradingHandler(tradingService)

			protected.GET("/trading/wallet", rateLimiter.RateLimit(), tradingHandler.GetWallet)
			protected.GET("/trading/portfolio", rateLimiter.RateLimit(), tradingHandler.GetPortfolio)
			protected.POST("/trading/buy", rateLimiter.RateLimit(), tradingHandler.BuyStock)
			protected.POST("/trading/sell", rateLimiter.RateLimit(), tradingHandler.SellStock)
			protected.GET("/trading/transactions", rateLimiter.RateLimit(), tradingHandler.GetTransactionHistory)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(r.cfg))

		authRepoAdmin := auth.NewRepository(r.client)
		admin.Use(middleware.AdminAuth(authRepoAdmin))
		{
			authService := auth.NewService(authRepoAdmin, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			admin.PUT("/users/:user_id/kyc", rateLimiter.RateLimit(), authHandler.UpdateKYC)

			// Admin seed endpoint
			// Admin seed endpoint
			stockRepoAdmin := stock.NewRepository(r.client)
			adminHandler := NewAdminHandler(stockRepoAdmin)
			admin.POST("/seed-stocks", adminHandler.SeedStockData)

		}
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
