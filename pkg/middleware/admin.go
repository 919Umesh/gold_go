package middleware

import (
	"net/http"

	"github.com/919Umesh/stock_market_sim/internal/auth"
	"github.com/gin-gonic/gin"
)

func AdminAuth(authRepo auth.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("user_id")

		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			ctx.Abort()
			return
		}

		userIdStr, ok := userId.(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id type"})
			ctx.Abort()
			return
		}

		user, err := authRepo.FindByID(userIdStr)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			ctx.Abort()
			return
		}

		if user.Role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			ctx.Abort()
			return
		}

		if user.KYCStatus != "verified" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "verified kyc required for admin access"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
