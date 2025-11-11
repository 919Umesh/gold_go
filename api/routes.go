package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/umesh/gold_investment/config"
	"github.com/umesh/gold_investment/internal/auth"
	"github.com/umesh/gold_investment/internal/gold"
	"github.com/umesh/gold_investment/internal/wallet"
	"github.com/umesh/gold_investment/pkg/middleware"
	"github.com/umesh/gold_investment/pkg/redis"
)

type Router struct {
	db          *gorm.DB
	cfg         *config.Config
	engine      *gin.Engine
	redisClient *redis.Client
}

func NewRouter(db *gorm.DB, cfg *config.Config) *Router {
	router := &Router{
		db:     db,
		cfg:    cfg,
		engine: gin.Default(),
	}

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
		public.Use(rateLimiter.RateLimit())
		{
			authRepo := auth.NewRepository(r.db)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			v1.POST("/auth/register", authHandler.Register)
			v1.POST("/auth/login", authHandler.Login)

			goldService := gold.NewService(r.db, r.cfg)
			goldHandler := gold.NewHandler(goldService)

			v1.GET("/gold/price",cacheMiddleware.Cache(1*time.Minute), goldHandler.GetCurrentPrice)
			v1.GET("/gold/history",cacheMiddleware.Cache(5*time.Minute), goldHandler.GetPriceHistory)
		}

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(r.cfg))
		protected.Use(rateLimiter.RateLimit())
		{
            authRepo := auth.NewRepository(r.db)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			protected.GET("/auth/profile", authHandler.GetProfile)

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
