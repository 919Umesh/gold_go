package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Portfolio represents a user's share holdings for a company
type Portfolio struct {
	ID          string          `json:"id,omitempty"`
	UserID      string          `json:"user_id"`
	CompanyID   string          `json:"company_id"`
	Quantity    int64           `json:"quantity"`
	AvgBuyPrice decimal.Decimal `json:"avg_buy_price"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at,omitempty"`
}
