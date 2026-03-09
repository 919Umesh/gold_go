package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/919Umesh/stock_market_sim/config"
	"github.com/919Umesh/stock_market_sim/internal/auth"
	"github.com/919Umesh/stock_market_sim/internal/ipo"
	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/internal/orderbook"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/internal/wallet"
	"github.com/919Umesh/stock_market_sim/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/shopspring/decimal"

	_ "github.com/919Umesh/stock_market_sim/docs"
)

type Router struct {
	client        *supabase.Client
	cfg           *config.Config
	engine        *gin.Engine
	eventHub      *market.EventHub
	priceEngine   *market.PriceEngine
	triggerWorker *market.TriggerWorker
}

func NewRouter(client *supabase.Client, cfg *config.Config) *Router {
	stockRepo := stock.NewRepository(client)
	eventHub := market.NewEventHub()
	priceEngine := market.NewPriceEngine(stockRepo, eventHub)
	triggerWorker := market.NewTriggerWorker(client)

	router := &Router{
		client:        client,
		cfg:           cfg,
		engine:        gin.Default(),
		eventHub:      eventHub,
		priceEngine:   priceEngine,
		triggerWorker: triggerWorker,
	}

	router.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.setupRoutes()
	return router
}

func (r *Router) setupRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": "2.0"})
	})

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ──────────────────── Init Repositories & Services ────────────────────

	authRepo := auth.NewRepository(r.client)
	stockRepo := stock.NewRepository(r.client)
	walletRepo := wallet.NewRepository(r.client)
	ipoRepo := ipo.NewRepository(r.client)
	obRepo := orderbook.NewRepository(r.client)

	authService := auth.NewService(authRepo, r.cfg.JWTSecret)
	walletService := wallet.NewService(walletRepo)
	ipoService := ipo.NewService(ipoRepo, stockRepo, walletService)
	obService := orderbook.NewService(obRepo, walletService, stockRepo, r.priceEngine)

	authHandler := auth.NewHandler(authService)
	walletHandler := wallet.NewHandler(walletService)
	ipoHandler := ipo.NewHandler(ipoService)
	obHandler := orderbook.NewHandler(obService)
	marketHandler := NewMarketHandler(r.priceEngine, r.eventHub, r.triggerWorker)

	// Wire up price trigger worker
	// When price updates → check triggers → auto-place sell orders
	r.priceEngine.SetOnPriceUpdate(func(companyID string, newPrice decimal.Decimal) {
		r.triggerWorker.CheckTriggers(companyID, newPrice)
	})

	r.triggerWorker.SetOrderPlacer(func(userID, companyID string, qty int64, price decimal.Decimal) error {
		_, _, err := obService.PlaceSellOrder(userID, companyID, qty, price)
		return err
	})

	// ──────────────────── Routes ────────────────────

	v1 := r.engine.Group("/api/v1")
	{
		// Auth (public)
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Market data (public)
		marketRoutes := v1.Group("/market")
		{
			marketRoutes.GET("/companies", marketHandler.ListCompanies)
			marketRoutes.GET("/companies/:symbol", marketHandler.GetCompanyDetail)
			marketRoutes.GET("/companies/:symbol/candles", marketHandler.GetCandlestickData)
			marketRoutes.GET("/live", marketHandler.GetLiveTradingData)
			marketRoutes.GET("/index", marketHandler.GetMarketIndex)
			marketRoutes.GET("/stream", marketHandler.StreamPrices)
		}

		// IPOs (public list/detail)
		ipoPublic := v1.Group("/ipos")
		{
			ipoPublic.GET("", ipoHandler.ListIPOs)
			ipoPublic.GET("/:id", ipoHandler.GetIPO)
		}

		// Order book (public view)
		v1.GET("/orderbook/:company_id", obHandler.GetOrderBook)

		// ──────────────────── Authenticated Routes ────────────────────
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(r.cfg))
		{
			// Profile
			profileRoutes := protected.Group("/auth")
			{
				profileRoutes.GET("/profile", authHandler.GetProfile)
				profileRoutes.PUT("/profile/update", authHandler.UpdateProfile)
				profileRoutes.POST("/profile/image", authHandler.UploadProfileImage)
			}

			// Wallet
			walletRoutes := protected.Group("/wallet")
			{
				walletRoutes.GET("", walletHandler.GetWallets)
				walletRoutes.GET("/main", walletHandler.GetMainWallet)
				walletRoutes.GET("/trading", walletHandler.GetTradingWallet)
				walletRoutes.POST("/topup", walletHandler.TopUp)
				walletRoutes.POST("/transfer", walletHandler.Transfer)
				walletRoutes.GET("/transfers", walletHandler.GetTransferHistory)
			}

			// IPO applications
			protected.POST("/ipos/:id/apply", ipoHandler.ApplyForIPO)

			// Orders
			orderRoutes := protected.Group("/orders")
			{
				orderRoutes.POST("/buy", obHandler.PlaceBuyOrder)
				orderRoutes.POST("/sell", obHandler.PlaceSellOrder)
				orderRoutes.DELETE("/:id", obHandler.CancelOrder)
				orderRoutes.GET("/my", obHandler.GetUserOrders)
			}

			// Portfolio & Trades
			protected.GET("/portfolio", obHandler.GetPortfolio)
			protected.GET("/trades", obHandler.GetUserTrades)

			// Price Triggers
			triggerRoutes := protected.Group("/triggers")
			{
				triggerRoutes.POST("", marketHandler.CreateTrigger)
				triggerRoutes.DELETE("/:id", marketHandler.CancelTrigger)
				triggerRoutes.GET("", marketHandler.GetUserTriggers)
			}
		}

		// ──────────────────── Admin Routes ────────────────────
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(r.cfg))
		admin.Use(middleware.AdminAuth(authRepo))
		{
			admin.POST("/companies", ipoHandler.CreateCompany)
			admin.POST("/ipos", ipoHandler.LaunchIPO)
			admin.POST("/ipos/:id/allocate", ipoHandler.AllocateIPO)
			admin.PUT("/users/:user_id/kyc", authHandler.UpdateKYC)
		}
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
