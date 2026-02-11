package api

import (
	"github.com/gin-gonic/gin"

	"github.com/919Umesh/stock_market_sim/config"
	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/919Umesh/stock_market_sim/internal/auth"
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

	//To allow the all origin for development purpose
	router.engine.Use(cors.Default())

	router.setupRoutes()
	return router
}

func (r *Router) setupRoutes() {
	// Health check endpoint for Render deployment (at root level)
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

			public.POST("/auth/register", authHandler.Register)
			public.POST("/auth/login", authHandler.Login)

			// Stock Market
			stockRepo := stock.NewRepository(r.client)
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
		}

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(r.cfg))
		{
			authRepo := auth.NewRepository(r.client)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.PUT("/auth/profile/update", authHandler.UpdateProfile)

			walletRepo := wallet.NewRepository(r.client)
			walletService := wallet.NewService(walletRepo, r.workerPool)
			walletHandler := wallet.NewHandler(walletService)

			protected.GET("/wallet", walletHandler.GetWallet)
			protected.GET("/transaction", walletHandler.GetUserTransaction)
			protected.POST("/wallet/topup", walletHandler.TopUp)

			// Stock Trading Protected Routes
			stockRepoProtected := stock.NewRepository(r.client)
			tradingRepo := trading.NewRepository(r.client)
			tradingService := trading.NewService(tradingRepo, stockRepoProtected)
			tradingHandler := NewTradingHandler(tradingService)

			protected.GET("/trading/wallet", tradingHandler.GetWallet)
			protected.GET("/trading/portfolio", tradingHandler.GetPortfolio)
			protected.POST("/trading/buy", tradingHandler.BuyStock)
			protected.POST("/trading/sell", tradingHandler.SellStock)
			protected.GET("/trading/transactions", tradingHandler.GetTransactionHistory)
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(r.cfg))

		authRepoAdmin := auth.NewRepository(r.client)
		admin.Use(middleware.AdminAuth(authRepoAdmin))
		{
			authService := auth.NewService(authRepoAdmin, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			admin.PUT("/users/:user_id/kyc", authHandler.UpdateKYC)

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
