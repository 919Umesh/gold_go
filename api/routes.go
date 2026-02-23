package api

import (
	"github.com/gin-gonic/gin"

	"github.com/919Umesh/stock_market_sim/config"
	"github.com/919Umesh/stock_market_sim/internal/auth"
	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/internal/ml"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/internal/trading"
	"github.com/919Umesh/stock_market_sim/internal/wallet"
	"github.com/919Umesh/stock_market_sim/pkg/middleware"
	"github.com/919Umesh/stock_market_sim/pkg/queue"
	"github.com/gin-contrib/cors"
)

type Router struct {
	client      *supabase.Client
	cfg         *config.Config
	engine      *gin.Engine
	workerPool  *queue.WorkerPool
	eventHub    *market.EventHub
	priceEngine *market.PriceEngine
}

func NewRouter(client *supabase.Client, cfg *config.Config, wp *queue.WorkerPool) *Router {
	stockRepo := stock.NewRepository(client)
	eventHub := market.NewEventHub()
	priceEngine := market.NewPriceEngine(stockRepo, eventHub)

	router := &Router{
		client:      client,
		cfg:         cfg,
		engine:      gin.Default(),
		workerPool:  wp,
		eventHub:    eventHub,
		priceEngine: priceEngine,
	}

	// global CORS settings – allow the frontend on localhost:3000 and any other origins you need
	router.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// Allow all subdomains or additional origins by adding them to the slice above (including http://localhost:8080 if you run client there)
	}))

	router.setupRoutes()
	return router
}

func (r *Router) setupRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authRepo := auth.NewRepository(r.client)
	stockRepo := stock.NewRepository(r.client)
	walletRepo := wallet.NewRepository(r.client)
	tradingRepo := trading.NewRepository(r.client)

	authService := auth.NewService(authRepo, r.cfg.JWTSecret)
	stockService := stock.NewService(stockRepo)
	mlService := ml.NewService(stockRepo)
	walletService := wallet.NewService(walletRepo, r.workerPool)
	tradingService := trading.NewService(tradingRepo, stockRepo, r.priceEngine)

	authHandler := auth.NewHandler(authService)
	stockHandler := NewStockHandler(stockService, r.priceEngine, r.eventHub)
	predictionHandler := NewPredictionHandler(mlService)
	walletHandler := wallet.NewHandler(walletService)
	tradingHandler := NewTradingHandler(tradingService)
	adminHandler := NewAdminHandler(stockRepo)

	v1 := r.engine.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		stockRoutes := v1.Group("/stocks")
		{
			stockRoutes.GET("", stockHandler.ListCompanies)
			stockRoutes.GET("/search", stockHandler.SearchCompanies)
			stockRoutes.GET("/market-overview", stockHandler.GetMarketOverview)
			stockRoutes.GET("/top-gainers", stockHandler.GetTopGainers)
			stockRoutes.GET("/top-losers", stockHandler.GetTopLosers)
			stockRoutes.GET("/most-active", stockHandler.GetMostActive)

			// New: Live trading data (like NEPSE live board)
			stockRoutes.GET("/live-trading", stockHandler.GetLiveTradingData)

			// New: Market index (overall market value)
			stockRoutes.GET("/market-index", stockHandler.GetMarketIndex)

			// New: Comprehensive market summary
			stockRoutes.GET("/market-summary", stockHandler.GetMarketSummary)

			// New: SSE stream for real-time price updates
			stockRoutes.GET("/stream", stockHandler.StreamPrices)

			stockRoutes.GET("/:symbol", stockHandler.GetCompany)
			stockRoutes.GET("/:symbol/price", stockHandler.GetCurrentPrice)
			stockRoutes.GET("/:symbol/history", stockHandler.GetPriceHistory)
			stockRoutes.GET("/:symbol/events", stockHandler.GetUpcomingEvents)

			// New: Candlestick OHLCV data for charting
			stockRoutes.GET("/:symbol/candles", stockHandler.GetCandlestickData)
		}

		sectorRoutes := v1.Group("/sectors")
		{
			sectorRoutes.GET("", stockHandler.GetAllSectors)
			sectorRoutes.GET("/:sector/companies", stockHandler.GetCompaniesBySector)
			sectorRoutes.GET("/:sector/stats", stockHandler.GetSectorStats)
		}

		predictionRoutes := v1.Group("/prediction")
		{
			predictionRoutes.GET("/algorithms", predictionHandler.ListAlgorithms)
			predictionRoutes.GET("/:symbol", predictionHandler.GetPrediction)
			predictionRoutes.GET("/:symbol/compare", predictionHandler.CompareAlgorithms)
		}

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(r.cfg))
		{
			profileRoutes := protected.Group("/auth")
			{
				profileRoutes.GET("/profile", authHandler.GetProfile)
				profileRoutes.PUT("/profile/update", authHandler.UpdateProfile)
				profileRoutes.POST("/profile/image", authHandler.UploadProfileImage)
			}

			walletRoutes := protected.Group("/wallet")
			{
				walletRoutes.GET("", walletHandler.GetWallet)
				walletRoutes.POST("/topup", walletHandler.TopUp)
			}
			protected.GET("/transaction", walletHandler.GetUserTransaction)

			tradingRoutes := protected.Group("/trading")
			{
				tradingRoutes.GET("/wallet", tradingHandler.GetWallet)
				tradingRoutes.GET("/portfolio", tradingHandler.GetPortfolio)
				tradingRoutes.POST("/buy", tradingHandler.BuyStock)
				tradingRoutes.POST("/sell", tradingHandler.SellStock)
				tradingRoutes.GET("/transactions", tradingHandler.GetTransactionHistory)
			}
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(r.cfg))
		admin.Use(middleware.AdminAuth(authRepo))
		{
			admin.PUT("/users/:user_id/kyc", authHandler.UpdateKYC)
			admin.POST("/seed-stocks", adminHandler.SeedStockData)
		}
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
