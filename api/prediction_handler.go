package api

import (
	"net/http"

	"github.com/919Umesh/stock_market_sim/internal/ml"
	"github.com/919Umesh/stock_market_sim/pkg/apperr"
	"github.com/gin-gonic/gin"
)

type PredictionHandler struct {
	service ml.Service
}

func NewPredictionHandler(service ml.Service) *PredictionHandler {
	return &PredictionHandler{service: service}
}

func (h *PredictionHandler) GetPrediction(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Symbol is required")
		return
	}

	result, err := h.service.PredictStockPrice(symbol, 1) // Next 1 day
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to predict price: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prediction": result,
	})
}
