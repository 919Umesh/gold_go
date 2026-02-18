package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Company struct {
	ID              string          `json:"id,omitempty"`
	Symbol          string          `json:"symbol"`
	Name            string          `json:"name"`
	Sector          string          `json:"sector"`
	MarketCap       decimal.Decimal `json:"market_cap"`
	Description     string          `json:"description"`
	FoundedYear     int             `json:"founded_year"`
	Employees       int             `json:"employees"`
	TotalShares     int64           `json:"total_shares"`     
	AvailableShares int64           `json:"available_shares"` 
	IsActive        bool            `json:"is_active"`
	CreatedAt       time.Time       `json:"created_at,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at,omitempty"`
}
