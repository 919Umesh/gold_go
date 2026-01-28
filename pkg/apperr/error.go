package apperr

import "github.com/gin-gonic/gin"

// AppError represents a structured error response
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Respond sends a JSON error response
func Respond(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// RespondWithMessage sends a JSON error response with a custom message
func RespondWithMessage(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
