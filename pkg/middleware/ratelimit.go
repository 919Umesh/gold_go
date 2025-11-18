package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/919Umesh/gold_go/pkg/redis"
	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	redisClient *redis.Client
}

func NewRateLimiter(redisClient *redis.Client) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
	}
}

type RateLimitConfig struct {
	Requests int
	Window   int
}

var endpointLimits = map[string]RateLimitConfig{
	"/api/v1/auth/login":         {Requests: 5, Window: 300},
	"/api/v1/gold/history":       {Requests: 50, Window: 60},
	"/api/v1/auth/profile":       {Requests: 60, Window: 60},
	"/api/v1/auth/register":      {Requests: 3, Window: 3600},
	"/api/v1/wallet/topup":       {Requests: 100, Window: 3600},
	"/api/v1/wallet/buy":         {Requests: 30, Window: 3600},
	"/api/v1/wallet/sell":        {Requests: 30, Window: 3600},
	"/api/v1/wallet/transaction": {Requests: 30, Window: 3600},
	"/auth/profile/update":       {Requests: 30, Window: 3600},
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var identifier string

		if userID, exists := ctx.Get("user_id"); exists {
			identifier = "user:" + strconv.FormatUint(uint64(userID.(uint)), 10)
		} else {
			identifier = "ip:" + ctx.ClientIP()
		}

		path := ctx.FullPath()
		config, exists := endpointLimits[path]
		if !exists {
			config = RateLimitConfig{Requests: 100, Window: 3600}
		}

		key := "rate_limit:" + identifier + ":" + path

		count, err := rl.redisClient.IncrementWithExpiry(ctx.Request.Context(), key, time.Duration(config.Window)*time.Second)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			ctx.Abort()
			return
		}

		if count > int64(config.Requests) {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			ctx.Abort()
			return
		}

		remaining := int64(config.Requests) - count
		if remaining < 0 {
			remaining = 0
		}

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
		ctx.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		ctx.Header("X-RateLimit-Reset", strconv.Itoa(config.Window))
		ctx.Next()
	}
}
