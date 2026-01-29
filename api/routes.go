package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/919Umesh/gold_go/config"
	"github.com/919Umesh/gold_go/internal/auth"
	"github.com/919Umesh/gold_go/internal/gold"
	"github.com/919Umesh/gold_go/internal/stock"
	"github.com/919Umesh/gold_go/internal/trading"
	"github.com/919Umesh/gold_go/internal/wallet"
	"github.com/919Umesh/gold_go/pkg/middleware"
	"github.com/919Umesh/gold_go/pkg/queue"
	"github.com/919Umesh/gold_go/pkg/redis"
	"github.com/gin-contrib/cors"
)

type Router struct {
	db          *gorm.DB
	cfg         *config.Config
	engine      *gin.Engine
	redisClient *redis.Client
	goldService *gold.Service
	workerPool  *queue.WorkerPool
}

func NewRouter(db *gorm.DB, cfg *config.Config, goldService *gold.Service, wp *queue.WorkerPool) *Router {
	router := &Router{
		db:          db,
		cfg:         cfg,
		engine:      gin.Default(),
		goldService: goldService,
		workerPool:  wp,
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
		panic("Failed to connect to redis: " + err.Error())
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
			authRepo := auth.NewRepository(r.db)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			public.POST("/auth/register", rateLimiter.RateLimit(), authHandler.Register)
			public.POST("/auth/login", rateLimiter.RateLimit(), authHandler.Login)

			goldHandler := gold.NewHandler(r.goldService)

			public.GET("/gold/price", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), goldHandler.GetCurrentPrice)
			public.GET("/gold/history", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), goldHandler.GetPriceHistory)

			// Stock Market Public Routes
			stockRepo := stock.NewRepository(r.db)
			stockService := stock.NewService(stockRepo)
			stockHandler := NewStockHandler(stockService)

			public.GET("/stocks", rateLimiter.RateLimit(), cacheMiddleware.Cache(2*time.Minute), stockHandler.ListCompanies)
			public.GET("/stocks/search", rateLimiter.RateLimit(), stockHandler.SearchCompanies)
			public.GET("/stocks/sector/:sector", rateLimiter.RateLimit(), cacheMiddleware.Cache(2*time.Minute), stockHandler.GetCompaniesBySector)
			public.GET("/stocks/market-overview", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), stockHandler.GetMarketOverview)
			public.GET("/stocks/top-gainers", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), stockHandler.GetTopGainers)
			public.GET("/stocks/top-losers", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), stockHandler.GetTopLosers)
			public.GET("/stocks/most-active", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), stockHandler.GetMostActive)
			public.GET("/stocks/:symbol", rateLimiter.RateLimit(), cacheMiddleware.Cache(2*time.Minute), stockHandler.GetCompany)
			public.GET("/stocks/:symbol/price", rateLimiter.RateLimit(), cacheMiddleware.Cache(30*time.Second), stockHandler.GetCurrentPrice)
			public.GET("/stocks/:symbol/history", rateLimiter.RateLimit(), cacheMiddleware.Cache(5*time.Minute), stockHandler.GetPriceHistory)
			public.GET("/stocks/:symbol/events", rateLimiter.RateLimit(), cacheMiddleware.Cache(5*time.Minute), stockHandler.GetUpcomingEvents)

		}

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(r.cfg))
		{
			authRepo := auth.NewRepository(r.db)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			protected.GET("/auth/profile", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), authHandler.GetProfile)
			protected.PUT("/auth/profile/update", rateLimiter.RateLimit(), authHandler.UpdateProfile)

			walletRepo := wallet.NewRepository(r.db)
			walletService := wallet.NewService(walletRepo, r.workerPool)
			walletHandler := wallet.NewHandler(walletService)

			protected.GET("/wallet", rateLimiter.RateLimit(), walletHandler.GetWallet)
			protected.GET("/transaction", rateLimiter.RateLimit(), walletHandler.GetUserTransaction)
			protected.POST("/wallet/topup", rateLimiter.RateLimit(), walletHandler.TopUp)
			protected.POST("/wallet/buy", rateLimiter.RateLimit(), walletHandler.BuyGold)
			protected.POST("/wallet/sell", rateLimiter.RateLimit(), walletHandler.SellGold)

			// Stock Trading Protected Routes
			stockRepo := stock.NewRepository(r.db)
			tradingRepo := trading.NewRepository(r.db)
			tradingService := trading.NewService(tradingRepo, stockRepo)
			tradingHandler := NewTradingHandler(tradingService)

			protected.GET("/trading/wallet", rateLimiter.RateLimit(), tradingHandler.GetWallet)
			protected.GET("/trading/portfolio", rateLimiter.RateLimit(), tradingHandler.GetPortfolio)
			protected.POST("/trading/buy", rateLimiter.RateLimit(), tradingHandler.BuyStock)
			protected.POST("/trading/sell", rateLimiter.RateLimit(), tradingHandler.SellStock)
			protected.GET("/trading/transactions", rateLimiter.RateLimit(), tradingHandler.GetTransactionHistory)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(r.cfg))
		admin.Use(middleware.AdminAuth(r.db))
		{
			authRepo := auth.NewRepository(r.db)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			admin.PUT("/users/:user_id/kyc", rateLimiter.RateLimit(), authHandler.UpdateKYC)

			// Admin seed endpoint (use once to populate stock data)
			adminHandler := NewAdminHandler(r.db)
			admin.POST("/seed-stocks", adminHandler.SeedStockData)

		}
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
