package api

import (
	"github.com/gin-gonic/gin"

	"github.com/919Umesh/stock_market_sim/config"
	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/919Umesh/stock_market_sim/internal/auth"
	"github.com/919Umesh/stock_market_sim/internal/ml"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/internal/trading"
	"github.com/919Umesh/stock_market_sim/internal/wallet"
	"github.com/919Umesh/stock_market_sim/pkg/middleware"
	"github.com/919Umesh/stock_market_sim/pkg/queue"
	"github.com/gin-contrib/cors"
)

type Router struct {
	client     *appwrite.Client
	cfg        *config.Config
	engine     *gin.Engine
	workerPool *queue.WorkerPool
}

func NewRouter(client *appwrite.Client, cfg *config.Config, wp *queue.WorkerPool) *Router {
	router := &Router{
		client:     client,
		cfg:        cfg,
		engine:     gin.Default(),
		workerPool: wp,
	}

	router.engine.Use(cors.Default())

	router.setupRoutes()
	return router
}

func (r *Router) setupRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize repositories once for reuse
	authRepo := auth.NewRepository(r.client)
	stockRepo := stock.NewRepository(r.client)

	v1 := r.engine.Group("/api/v1")
	{

		public := v1.Group("")
		{
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			public.POST("/auth/register", authHandler.Register)
			public.POST("/auth/login", authHandler.Login)

			// Stock Market
			stockService := stock.NewService(stockRepo)
			stockHandler := NewStockHandler(stockService)

			public.GET("/stocks", stockHandler.ListCompanies)
			public.GET("/stocks/search", stockHandler.SearchCompanies)
			public.GET("/stocks/:symbol", stockHandler.GetCompany)
			public.GET("/stocks/:symbol/price", stockHandler.GetCurrentPrice)
			public.GET("/stocks/:symbol/history", stockHandler.GetPriceHistory)
			public.GET("/stocks/market-overview", stockHandler.GetMarketOverview)
			public.GET("/stocks/top-gainers", stockHandler.GetTopGainers)
			public.GET("/stocks/top-losers", stockHandler.GetTopLosers)
			public.GET("/stocks/most-active", stockHandler.GetMostActive)
			public.GET("/stocks/:symbol/events", stockHandler.GetUpcomingEvents)

			// Sector APIs
			public.GET("/sectors", stockHandler.GetAllSectors)
			public.GET("/sectors/:sector/companies", stockHandler.GetCompaniesBySector)
			public.GET("/sectors/:sector/stats", stockHandler.GetSectorStats)

			// Prediction API
			mlService := ml.NewService(stockRepo)
			predictionHandler := NewPredictionHandler(mlService)
			public.GET("/prediction/:symbol", predictionHandler.GetPrediction)
		}

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(r.cfg))
		{
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.PUT("/auth/profile/update", authHandler.UpdateProfile)
			protected.POST("/auth/profile/image", authHandler.UploadProfileImage)

			walletRepo := wallet.NewRepository(r.client)
			walletService := wallet.NewService(walletRepo, r.workerPool)
			walletHandler := wallet.NewHandler(walletService)

			protected.GET("/wallet", walletHandler.GetWallet)
			protected.GET("/transaction", walletHandler.GetUserTransaction)
			protected.POST("/wallet/topup", walletHandler.TopUp)

			tradingRepo := trading.NewRepository(r.client)
			tradingService := trading.NewService(tradingRepo, stockRepo)
			tradingHandler := NewTradingHandler(tradingService)

			protected.GET("/trading/wallet", tradingHandler.GetWallet)
			protected.GET("/trading/portfolio", tradingHandler.GetPortfolio)
			protected.POST("/trading/buy", tradingHandler.BuyStock)
			protected.POST("/trading/sell", tradingHandler.SellStock)
			protected.GET("/trading/transactions", tradingHandler.GetTransactionHistory)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(r.cfg))
		admin.Use(middleware.AdminAuth(authRepo))
		{
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			admin.PUT("/users/:user_id/kyc", authHandler.UpdateKYC)

			adminHandler := NewAdminHandler(stockRepo)
			admin.POST("/seed-stocks", adminHandler.SeedStockData)

		}
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
