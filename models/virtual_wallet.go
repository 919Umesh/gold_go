package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type VirtualWallet struct {
	ID              string          `json:"id,omitempty"`
	UserID          string          `json:"user_id"`
	Balance         decimal.Decimal `json:"balance"`
	TotalInvested   decimal.Decimal `json:"total_invested"`
	TotalProfitLoss decimal.Decimal `json:"total_profit_loss"`
	CreatedAt       time.Time       `json:"created_at,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at,omitempty"`

	User *User `json:"user,omitempty"`
}
