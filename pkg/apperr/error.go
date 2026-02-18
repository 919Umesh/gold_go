package apperr

import "github.com/gin-gonic/gin"


type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Respond(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

func RespondWithMessage(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
