package gold

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetCurrentPrice(c *gin.Context) {
	price, updatedAt, err := h.service.GetCurrentPrice()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "price not available"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"price_per_gram": price,
		"updated_at":     updatedAt,
		"currency":       "NPR",
	})
}

func (h *Handler) GetPriceHistory(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 30 {
		days = 7
	}

	prices, err := h.service.GetPriceHistory(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch price history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prices": prices,
		"days":   days,
	})
}
