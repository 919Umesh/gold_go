package models

import (
	"time"
)


type VirtualWallet struct {
	ID              string    `json:"$id,omitempty"`
	UserID          string    `json:"user_id"`
	Balance         float64   `json:"balance"` 
	TotalInvested   float64   `json:"total_invested"`
	TotalProfitLoss float64   `json:"total_profit_loss"`
	CreatedAt       time.Time `json:"$createdAt,omitempty"`
	UpdatedAt       time.Time `json:"$updatedAt,omitempty"`


	User *User `json:"user,omitempty"`
}
