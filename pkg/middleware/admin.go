package middleware

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuth(db *sql.DB) gin.HandlerFunc {
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
		err := db.QueryRow(query, userId).Scan(&role)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found or KYC not verified"})
			} else {
				log.Printf("Database error: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			ctx.Abort()
			return
		}
		log.Print("----------Raw_Query-------------")
		log.Print(role)

		if role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
