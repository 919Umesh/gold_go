package wallet

import (
	"net/http"

	"github.com/919Umesh/stock_market_sim/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ──────────────────── Get Wallets ────────────────────

// GetWallets godoc
func (h *Handler) GetWallets(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	main, trading, err := h.service.GetBothWallets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"main_wallet":    main,
		"trading_wallet": trading,
	})
}

// GetMainWallet godoc
func (h *Handler) GetMainWallet(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	w, err := h.service.GetMainWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": w})
}

// GetTradingWallet godoc
func (h *Handler) GetTradingWallet(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	w, err := h.service.GetTradingWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": w})
}

// ──────────────────── Top Up ────────────────────

type TopUpRequest struct {
	Amount string `json:"amount" binding:"required"`
}

// TopUp godoc
func (h *Handler) TopUp(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || !amount.IsPositive() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	wallet, err := h.service.TopUp(userID, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "top-up successful",
		"wallet":  wallet,
	})
}

// ──────────────────── Transfer ────────────────────

type TransferRequest struct {
	Amount    string `json:"amount" binding:"required"`
	Direction string `json:"direction" binding:"required"` // main_to_trading, trading_to_main
}

// Transfer godoc
func (h *Handler) Transfer(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || !amount.IsPositive() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	var mainW *models.MainWallet
	var tradingW *models.TradingWallet

	switch req.Direction {
	case "main_to_trading":
		mainW, tradingW, err = h.service.TransferToTrading(userID, amount)
	case "trading_to_main":
		mainW, tradingW, err = h.service.TransferToMain(userID, amount)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "direction must be 'main_to_trading' or 'trading_to_main'"})
		return
	}

	if err != nil {
		status := http.StatusBadRequest
		if err == ErrInsufficientBalance || err == ErrInsufficientAvailable {
			status = http.StatusUnprocessableEntity
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "transfer successful",
		"main_wallet":    mainW,
		"trading_wallet": tradingW,
	})
}

// ──────────────────── Transfer History ────────────────────

// GetTransferHistory godoc
func (h *Handler) GetTransferHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	transfers, err := h.service.GetTransferHistory(userID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transfers": transfers,
	})
}
