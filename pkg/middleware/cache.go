package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/umesh/gold_investment/pkg/redis"
)

type CacheMiddleWare struct {
	redisClient *redis.Client
}

func NewCacheMiddleWare(redisClient *redis.Client) *CacheMiddleWare {
	return &CacheMiddleWare{
		redisClient: redisClient,
	}
}

func (cm *CacheMiddleWare) Cache(duration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if ctx.Request.Method != "GET" {
			ctx.Next()
			return
		}

		cacheKey := cm.generateCacheKey(ctx)

		cached, err := cm.redisClient.Get(ctx.Request.Context(), cacheKey)
		if err != nil && cached != "" {
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
			go func() {
				ctx := context.Background()
				cm.redisClient.Set(ctx, cacheKey, string(blw.body), duration)
			}()
		}

	}
}

func (cm *CacheMiddleWare) generateCacheKey(c *gin.Context) string {
	key := c.Request.URL.String()

	if userID, exists := c.Get("user_id"); exists {
		key += "user:" + strconv.FormatUint(uint64(userID.(uint)), 10)
	}

	hash := sha256.Sum256([]byte(key))
	return "Cache" + hex.EncodeToString(hash[:])
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}
