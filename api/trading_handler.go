package api

import (
	"net/http"
	"strconv"

	"github.com/919Umesh/stock_market_sim/internal/trading"
	"github.com/919Umesh/stock_market_sim/pkg/apperr"
	"github.com/gin-gonic/gin"
)

type TradingHandler struct {
	service trading.Service
}

func NewTradingHandler(service trading.Service) *TradingHandler {
	return &TradingHandler{service: service}
}

type BuyRequest struct {
	Symbol   string `json:"symbol" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type SellRequest struct {
	Symbol   string `json:"symbol" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}


func (h *TradingHandler) GetWallet(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	wallet, err := h.service.GetOrCreateWallet(userID)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to get wallet")
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

// GetPortfolio godoc
// @Summary Get user portfolio
// @Tags trading
// @Security BearerAuth
// @Success 200 {object} trading.PortfolioSummary
// @Router /api/v1/trading/portfolio [get]
func (h *TradingHandler) GetPortfolio(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	portfolio, err := h.service.GetPortfolio(userID)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to get portfolio")
		return
	}

	c.JSON(http.StatusOK, portfolio)
}


func (h *TradingHandler) BuyStock(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req BuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Respond(c, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.BuyStock(userID, req.Symbol, req.Quantity)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TradingHandler) SellStock(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req SellRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Respond(c, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.SellStock(userID, req.Symbol, req.Quantity)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, result)
		return
	}

	c.JSON(http.StatusOK, result)
}


func (h *TradingHandler) GetTransactionHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.service.GetTransactionHistory(userID, limit, offset)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to get transactions")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	})
}
