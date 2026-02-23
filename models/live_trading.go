package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// LiveTradingData represents real-time trading data for a single stock
// Matches the NEPSE live trading board format:
// Symbol | LTP | % Change | Open | High | Low | Qty | PClose | Diff
type LiveTradingData struct {
	Symbol        string          `json:"symbol"`
	CompanyID     string          `json:"company_id"`
	CompanyName   string          `json:"company_name"`
	Sector        string          `json:"sector"`
	LTP           decimal.Decimal `json:"ltp"`
	ChangePercent decimal.Decimal `json:"change_percent"`
	Open          decimal.Decimal `json:"open"`
	High          decimal.Decimal `json:"high"`
	Low           decimal.Decimal `json:"low"`
	Volume        int64           `json:"volume"`
	PreviousClose decimal.Decimal `json:"previous_close"`
	Difference    decimal.Decimal `json:"difference"`
	Turnover      decimal.Decimal `json:"turnover"`
	LastUpdated   time.Time       `json:"last_updated"`
}

// TradeFeedItem represents a single trade in the live trade feed
type TradeFeedItem struct {
	Symbol      string          `json:"symbol"`
	CompanyName string          `json:"company_name"`
	TradeType   string          `json:"trade_type"` // "buy" or "sell"
	Quantity    int             `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	PriceImpact decimal.Decimal `json:"price_impact"`
	NewPrice    decimal.Decimal `json:"new_price"`
	Timestamp   time.Time       `json:"timestamp"`
}

// CandlestickData represents OHLCV data for charting (candlestick format)
type CandlestickData struct {
	Timestamp time.Time       `json:"timestamp"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    int64           `json:"volume"`
}
