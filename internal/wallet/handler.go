package wallet

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/919Umesh/stock_market_sim/models"
	"github.com/919Umesh/stock_market_sim/pkg/notification"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *Handler) GetWallet(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	wallet, err := h.service.GetWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

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

	referenceID := "topup_" + uuid.New().String()
	wallet, transaction, userEmail, err := h.service.TopUp(userID, req.Amount, "", referenceID)
	if err != nil {
		switch err {
		case ErrInvalidAmount:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		case ErrWalletLocked:
			c.JSON(http.StatusLocked, gin.H{"error": "wallet is locked"})
		default:
			slog.Error("TopUp failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "top-up failed"})
		}
		return
	}

	topupData := map[string]any{
		"message":        "Top-up successful",
		"amount":         transaction.Amount,
		"status":         transaction.Status,
		"reference_id":   transaction.ReferenceID,
		"created_at":     transaction.CreatedAt,
		"updated_at":     transaction.UpdatedAt,
		"wallet_balance": wallet.FiatBalance,
	}

	if userEmail != "" {
		notification.SendEmail(userEmail, "Top-Up", topupData)
	} else {
		slog.Warn("no user email available to send top-up notification", "user_id", userID)
	}

	push := models.PushInput{
		CustomerID:       userID,
		Title:            "Wallet Top-Up",
		Message:          fmt.Sprintf("₹%.2f has been added to your wallet", transaction.Amount.InexactFloat64()),
		NotificationType: "TOP-UP",
		LinkedID:         transaction.ID,
		Data: map[string]interface{}{
			"amount":         transaction.Amount.InexactFloat64(),
			"transaction_id": transaction.ID,
			"reference_id":   transaction.ReferenceID,
			"new_balance":    wallet.FiatBalance.InexactFloat64(),
		},
	}
	_, err = notification.SendOneSignalPush(push)
	if err != nil {
		slog.Error("failed to send push notification", "error", err)

	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "top-up successful",
		"wallet":      wallet,
		"transaction": transaction,
	})
}

func (h *Handler) GetUserTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	transaction, err := h.service.GetUserTransaction(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user transactions retrieved successfully",
		"data":    transaction,
	})

}
