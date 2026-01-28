package middleware

import (
	"net/http"

	"github.com/919Umesh/gold_go/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminAuth(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("user_id")

		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authencticated"})
			ctx.Abort()
			return
		}

		query := ` 
               SELECT role 
               FROM users 
               WHERE id = ? 
               AND kyc_status = 'verified'
             `
		var role string

		err := db.Raw(query, userId).Scan(&role).Error
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			ctx.Abort()
			return
		}
		var user models.User
		userQuery := `SELECT * FROM users WHERE id = ? LIMIT 1`
		if err := db.Raw(userQuery, userId).Scan(&user).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			ctx.Abort()
			return
		}
		if user.ID == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			ctx.Abort()
			return
		}

		if user.Role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
