package models

import (
	"time"
)

// VirtualWallet represents a user's virtual trading wallet
type VirtualWallet struct {
	ID              string    `json:"$id,omitempty"`
	UserID          string    `json:"user_id"`
	Balance         float64   `json:"balance"` // NPR 10 lakh initial
	TotalInvested   float64   `json:"total_invested"`
	TotalProfitLoss float64   `json:"total_profit_loss"`
	CreatedAt       time.Time `json:"$createdAt,omitempty"`
	UpdatedAt       time.Time `json:"$updatedAt,omitempty"`

	// Relationships
	User *User `json:"user,omitempty"`
}
