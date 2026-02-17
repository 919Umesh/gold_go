package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type StockPrice struct {
	ID         string          `json:"id,omitempty"`
	CompanyID  string          `json:"company_id"`
	OpenPrice  decimal.Decimal `json:"open_price"`
	HighPrice  decimal.Decimal `json:"high_price"`
	LowPrice   decimal.Decimal `json:"low_price"`
	ClosePrice decimal.Decimal `json:"close_price"`
	Volume     int64           `json:"volume"`
	Timestamp  time.Time       `json:"timestamp"`
	Timeframe  string          `json:"timeframe"`
	CreatedAt  time.Time       `json:"created_at,omitempty"`

	Company *Company `json:"company,omitempty"`
}

func (s *StockPrice) CalculateChange() decimal.Decimal {
	if s.OpenPrice.IsZero() {
		return decimal.Zero
	}
	// ((Close - Open) / Open) * 100
	return s.ClosePrice.Sub(s.OpenPrice).Div(s.OpenPrice).Mul(decimal.NewFromInt(100))
}
