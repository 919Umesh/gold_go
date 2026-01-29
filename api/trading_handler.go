package api

import (
	"net/http"
	"strconv"

	"github.com/919Umesh/gold_go/internal/trading"
	"github.com/919Umesh/gold_go/pkg/apperr"
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

// GetWallet godoc
// @Summary Get virtual wallet
// @Tags trading
// @Security BearerAuth
// @Success 200 {object} models.VirtualWallet
// @Router /api/v1/trading/wallet [get]
func (h *TradingHandler) GetWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	wallet, err := h.service.GetOrCreateWallet(userID.(uint))
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
	userID, exists := c.Get("user_id")
	if !exists {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	portfolio, err := h.service.GetPortfolio(userID.(uint))
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to get portfolio")
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// BuyStock godoc
// @Summary Buy stocks
// @Tags trading
// @Security BearerAuth
// @Param request body BuyRequest true "Buy request"
// @Success 200 {object} trading.TradeResult
// @Router /api/v1/trading/buy [post]
func (h *TradingHandler) BuyStock(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req BuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Respond(c, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.BuyStock(userID.(uint), req.Symbol, req.Quantity)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Transaction failed")
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// SellStock godoc
// @Summary Sell stocks
// @Tags trading
// @Security BearerAuth
// @Param request body SellRequest true "Sell request"
// @Success 200 {object} trading.TradeResult
// @Router /api/v1/trading/sell [post]
func (h *TradingHandler) SellStock(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req SellRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperr.Respond(c, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.SellStock(userID.(uint), req.Symbol, req.Quantity)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Transaction failed")
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTransactionHistory godoc
// @Summary Get transaction history
// @Tags trading
// @Security BearerAuth
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.StockTransaction
// @Router /api/v1/trading/transactions [get]
func (h *TradingHandler) GetTransactionHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apperr.RespondWithMessage(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.service.GetTransactionHistory(userID.(uint), limit, offset)
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
