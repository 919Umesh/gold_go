package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/919Umesh/gold_go/pkg/redis"
	"github.com/gin-gonic/gin"
)

type CacheMiddleware struct {
	redisClient *redis.Client
}

func NewCacheMiddleware(redisClient *redis.Client) *CacheMiddleware {
	return &CacheMiddleware{
		redisClient: redisClient,
	}
}

func (cm *CacheMiddleware) Cache(duration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != "GET" {
			ctx.Next()
			return
		}

		cacheKey := cm.generateCacheKey(ctx)

		cached, err := cm.redisClient.Get(ctx.Request.Context(), cacheKey)
		if err == nil && cached != "" {
			ctx.Header("X-Cache", "HIT")
			ctx.Data(200, "application/json", []byte(cached))
			ctx.Abort()
			return
		}

		ctx.Header("X-Cache", "MISS")
		blw := &bodyLogWriter{body: []byte{}, ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		ctx.Next()

		if ctx.Writer.Status() == 200 {
			go func(key string, body []byte, dur time.Duration) {
				cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				cm.redisClient.Set(cacheCtx, key, string(body), dur)
			}(cacheKey, blw.body, duration)
		}
	}
}

func (cm *CacheMiddleware) generateCacheKey(c *gin.Context) string {
	key := c.Request.URL.String()

	if userID, exists := c.Get("user_id"); exists {
		key += ":user:" + strconv.FormatUint(uint64(userID.(uint)), 10)
	}

	if c.Request.URL.RawQuery != "" {
		key += "?" + c.Request.URL.RawQuery
	}

	hash := sha256.Sum256([]byte(key))
	return "cache:" + hex.EncodeToString(hash[:])
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}
