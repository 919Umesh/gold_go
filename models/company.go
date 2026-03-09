package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Company struct {
	ID                string          `json:"id,omitempty"`
	Symbol            string          `json:"symbol"`
	Name              string          `json:"name"`
	Sector            string          `json:"sector"`
	Description       string          `json:"description,omitempty"`
	TotalSupply       int64           `json:"total_supply"`
	SharesOutstanding decimal.Decimal `json:"shares_outstanding"`
	CurrentPrice      decimal.Decimal `json:"current_price"`
	MarketCap         decimal.Decimal `json:"market_cap"`
	EPS               decimal.Decimal `json:"eps"`
	PERatio           decimal.Decimal `json:"pe_ratio"`
	BookValue         decimal.Decimal `json:"book_value"`
	PBV               decimal.Decimal `json:"pbv"`
	Week52High        decimal.Decimal `json:"week_52_high"`
	Week52Low         decimal.Decimal `json:"week_52_low"`
	Avg120Day         decimal.Decimal `json:"avg_120_day"`
	Yield1Year        decimal.Decimal `json:"yield_1_year"`
	ListedDate        string          `json:"listed_date,omitempty"`
	IsActive          bool            `json:"is_active"`
	CreatedAt         time.Time       `json:"created_at,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at,omitempty"`
}
