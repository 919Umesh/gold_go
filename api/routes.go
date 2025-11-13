package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/919Umesh/gold_go/config"
	"github.com/919Umesh/gold_go/internal/auth"
	"github.com/919Umesh/gold_go/internal/gold"
	"github.com/919Umesh/gold_go/internal/wallet"
	"github.com/919Umesh/gold_go/pkg/middleware"
	"github.com/919Umesh/gold_go/pkg/redis"
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
		{
			authRepo := auth.NewRepository(r.db)
			authService := auth.NewService(authRepo, r.cfg.JWTSecret)
			authHandler := auth.NewHandler(authService)

			public.POST("/auth/register", rateLimiter.RateLimit(), authHandler.Register)
			public.POST("/auth/login", rateLimiter.RateLimit(), authHandler.Login)

			goldService := gold.NewService(r.db, r.cfg)
			goldHandler := gold.NewHandler(goldService)

			public.GET("/gold/price", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), goldHandler.GetCurrentPrice)
			public.GET("/gold/history", rateLimiter.RateLimit(), cacheMiddleware.Cache(1*time.Minute), goldHandler.GetPriceHistory)

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
			walletService := wallet.NewService(walletRepo)
			walletHandler := wallet.NewHandler(walletService)

			protected.GET("/wallet", rateLimiter.RateLimit(), walletHandler.GetWallet)
			protected.POST("/wallet/topup", rateLimiter.RateLimit(), walletHandler.TopUp)
			protected.POST("/wallet/buy", rateLimiter.RateLimit(), walletHandler.BuyGold)
			protected.POST("/wallet/sell", rateLimiter.RateLimit(), walletHandler.SellGold)
		}
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
