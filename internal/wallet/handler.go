package wallet

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type BuyGoldRequest struct {
	Grams        float64 `json:"grams" binding:"required,gt=0"`
	PricePerGram float64 `json:"price_per_gram" binding:"required,gt=0"`
}

func (h *Handler) GetWallet(c *gin.Context) {
	userID := c.GetUint("user_id")

	wallet, err := h.service.GetWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

func (h *Handler) TopUp(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.service.TopUp(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "top-up successful",
		"wallet":  wallet,
	})
}

func (h *Handler) BuyGold(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req BuyGoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.service.BuyGold(userID, req.Grams, req.PricePerGram)
	if err != nil {
		switch err {
		case ErrInsufficientBalance:
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient fiat balance"})
		case ErrWalletLocked:
			c.JSON(http.StatusLocked, gin.H{"error": "wallet is locked"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "gold purchase successful",
		"wallet":  wallet,
	})
}

func (h *Handler) SellGold(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req BuyGoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.service.SellGold(userID, req.Grams, req.PricePerGram)
	if err != nil {
		switch err {
		case ErrInsufficientBalance:
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient gold balance"})
		case ErrWalletLocked:
			c.JSON(http.StatusLocked, gin.H{"error": "wallet is locked"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "gold sale successful",
		"wallet":  wallet,
	})
}
