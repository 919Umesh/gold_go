package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Company struct {
	ID           string          `json:"id,omitempty"`
	Symbol       string          `json:"symbol"`
	Name         string          `json:"name"`
	Sector       string          `json:"sector"`
	TotalSupply  int64           `json:"total_supply"`
	CurrentPrice decimal.Decimal `json:"current_price"`
	MarketCap    decimal.Decimal `json:"market_cap"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at,omitempty"`
	UpdatedAt    time.Time       `json:"updated_at,omitempty"`
}
