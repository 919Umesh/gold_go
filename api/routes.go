package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/umesh/gold_investment/config"
	"github.com/umesh/gold_investment/internal/auth"
	"github.com/umesh/gold_investment/internal/gold"
	"github.com/umesh/gold_investment/internal/wallet"
	"github.com/umesh/gold_investment/pkg/middleware"
)

type Router struct {
	db     *gorm.DB
	cfg    *config.Config
	engine *gin.Engine
}

func NewRouter(db *gorm.DB, cfg *config.Config) *Router {
	router := &Router{
		db:     db,
		cfg:    cfg,
		engine: gin.Default(),
	}

	router.setupRoutes()
	return router
}

func (r *Router) setupRoutes() {
	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Auth routes
		authRepo := auth.NewRepository(r.db)
		authService := auth.NewService(authRepo, r.cfg.JWTSecret)
		authHandler := auth.NewHandler(authService)

		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)

		// Gold price routes (public)
		goldService := gold.NewService(r.db, r.cfg)
		goldHandler := gold.NewHandler(goldService)

		v1.GET("/gold/price", goldHandler.GetCurrentPrice)
		v1.GET("/gold/history", goldHandler.GetPriceHistory)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(r.cfg))
		{
			protected.GET("/auth/profile", authHandler.GetProfile)

			// Wallet routes
			walletRepo := wallet.NewRepository(r.db)
			walletService := wallet.NewService(walletRepo)
			walletHandler := wallet.NewHandler(walletService)

			protected.GET("/wallet", walletHandler.GetWallet)
			protected.POST("/wallet/topup", walletHandler.TopUp)
			protected.POST("/wallet/buy", walletHandler.BuyGold)
			protected.POST("/wallet/sell", walletHandler.SellGold)
		}
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
