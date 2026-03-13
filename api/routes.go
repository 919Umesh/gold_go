package api

import (
	"github.com/919Umesh/stock_market_sim/config"
	"github.com/919Umesh/stock_market_sim/internal/auth"
	"github.com/919Umesh/stock_market_sim/internal/event"
	"github.com/919Umesh/stock_market_sim/internal/ipo"
	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/internal/orderbook"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/internal/wallet"
	"github.com/919Umesh/stock_market_sim/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func NewRouter(supabaseClient *supabase.Client, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// ─── Repositories ───
	authRepo := auth.NewRepository(supabaseClient)
	stockRepo := stock.NewRepository(supabaseClient)
	walletRepo := wallet.NewRepository(supabaseClient)
	orderRepo := orderbook.NewRepository(supabaseClient)
	ipoRepo := ipo.NewRepository(supabaseClient)
	eventRepo := event.NewRepository(supabaseClient)

	// ─── Services ───
	authService := auth.NewService(authRepo, cfg.JWTSecret)
	stockService := stock.NewService(stockRepo)
	walletService := wallet.NewService(walletRepo)
	eventService := event.NewService(eventRepo)

	// ─── Market Infrastructure ───
	eventHub := market.NewEventHub()

	priceEngine := market.NewPriceEngine(stockRepo, eventHub, orderRepo)
	triggerWorker := market.NewTriggerWorker(supabaseClient)

	// ─── Order Book ───
	orderService := orderbook.NewService(orderRepo, walletService, stockRepo, priceEngine)
	orderEngine := orderbook.NewEngine(orderRepo)

	// Connect trigger worker to price engine
	priceEngine.SetOnPriceUpdate(func(companyID string, newPrice decimal.Decimal) {
		triggerWorker.CheckTriggers(companyID, newPrice)
	})

	_ = orderEngine // engine processes orders internally

	// ─── IPO ───
	ipoService := ipo.NewService(ipoRepo, stockRepo, walletService)

	// ─── Handlers ───
	authHandler := auth.NewHandler(authService)
	walletHandler := wallet.NewHandler(walletService)
	orderHandler := orderbook.NewHandler(orderService)
	ipoHandler := ipo.NewHandler(ipoService)
	eventHandler := event.NewHandler(eventService)
	marketHandler := NewMarketHandler(priceEngine, stockService, eventHub, triggerWorker)

	// ─── Middleware ───
	authMiddleware := middleware.JWTAuth(cfg)
	adminMiddleware := middleware.AdminAuth(authRepo)

	// ─── Routes ───
	// Base level health check for Render
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "stock-market-simulator"})
	})

	v1 := r.Group("/api/v1")

	// Health
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "stock-market-simulator"})
	})

	// Auth
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/profile", authMiddleware, authHandler.GetProfile)
		authGroup.PUT("/profile/update", authMiddleware, authHandler.UpdateProfile)
		authGroup.POST("/profile/image", authMiddleware, authHandler.UploadProfileImage)
	}

	// Admin
	adminGroup := v1.Group("/admin", authMiddleware, adminMiddleware)
	{
		// User Management
		adminGroup.GET("/users", authHandler.GetAllUsers)
		adminGroup.PUT("/users/:user_id/kyc", authHandler.UpdateKYC)

		// IPO/Company Management
		adminGroup.POST("/companies", ipoHandler.CreateCompany)
		adminGroup.POST("/ipos", ipoHandler.LaunchIPO)
		adminGroup.POST("/ipos/:id/allocate", ipoHandler.AllocateIPO)
		adminGroup.GET("/ipos/:id/applications", ipoHandler.GetIPOApplications)
	}

	// Wallet
	walletGroup := v1.Group("/wallet", authMiddleware)
	{
		walletGroup.GET("/", walletHandler.GetWallets)
		walletGroup.GET("/main", walletHandler.GetMainWallet)
		walletGroup.GET("/trading", walletHandler.GetTradingWallet)
		walletGroup.POST("/topup", walletHandler.TopUp)
		walletGroup.POST("/transfer", walletHandler.Transfer)
		walletGroup.GET("/transfers", walletHandler.GetTransferHistory)
	}

	// Orders / Trading
	orderGroup := v1.Group("/orders", authMiddleware)
	{
		orderGroup.POST("/buy", orderHandler.PlaceBuyOrder)
		orderGroup.POST("/sell", orderHandler.PlaceSellOrder)
		orderGroup.PUT("/:id/cancel", orderHandler.CancelOrder)
		orderGroup.GET("/book/:company_id", orderHandler.GetOrderBook)
		orderGroup.GET("/my-orders", orderHandler.GetUserOrders)
		orderGroup.GET("/portfolio", orderHandler.GetPortfolio)
		orderGroup.GET("/trades", orderHandler.GetUserTrades)
	}

	// IPO
	ipoGroup := v1.Group("/ipo")
	{
		ipoGroup.GET("/", ipoHandler.ListIPOs)
		ipoGroup.GET("/:id", ipoHandler.GetIPO)
		ipoGroup.POST("/:id/apply", authMiddleware, ipoHandler.ApplyForIPO)
	}

	// Market Data (public)
	marketGroup := v1.Group("/market")
	{
		marketGroup.GET("/companies", marketHandler.ListCompanies)
		marketGroup.GET("/companies/new", marketHandler.GetNewCompanies)
		marketGroup.GET("/companies/old", marketHandler.GetOldCompanies)
		marketGroup.GET("/companies/:id", marketHandler.GetCompanyDetail)
		marketGroup.GET("/companies/:id/prediction", marketHandler.GetPricePrediction)
		marketGroup.GET("/companies/:id/trades", orderHandler.GetCompanyTrades)

		marketGroup.GET("/live", marketHandler.GetLiveTradingData)
		marketGroup.GET("/index", marketHandler.GetMarketIndex)
		marketGroup.GET("/candlestick", marketHandler.GetCandlestickData)
		marketGroup.GET("/candlestick/1d", marketHandler.GetCandlestickData)
		marketGroup.GET("/chart/1d", marketHandler.GetCandlestickData) // Professional Charting Endpoint

		marketGroup.GET("/top-gainers", marketHandler.GetTopGainers)
		marketGroup.GET("/top-losers", marketHandler.GetTopLosers)
		marketGroup.GET("/most-active", marketHandler.GetMostActive)
		marketGroup.GET("/top-turnover", marketHandler.GetTopTurnover)

		marketGroup.GET("/sectors", marketHandler.GetTopSectors)
		marketGroup.GET("/sectors/:sector/companies", marketHandler.GetCompaniesBySector)

		marketGroup.GET("/stream", marketHandler.StreamPrices)

		// Price triggers (authenticated)
		marketGroup.POST("/triggers", authMiddleware, marketHandler.CreateTrigger)
		marketGroup.GET("/triggers", authMiddleware, marketHandler.GetUserTriggers)
		marketGroup.PUT("/triggers/:id/cancel", authMiddleware, marketHandler.CancelTrigger)
	}

	// Company Events (public)
	eventGroup := v1.Group("/market")
	{
		eventGroup.GET("/events", eventHandler.GetAllEvents)
		eventGroup.GET("/events/company/:company_id", eventHandler.GetCompanyEvents)
		eventGroup.GET("/events/upcoming", eventHandler.GetUpcomingEvents)
		eventGroup.GET("/events/type/:event_type", eventHandler.GetEventsByType)
	}

	return r
}
