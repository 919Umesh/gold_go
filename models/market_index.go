package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// MarketIndex represents the overall market index (like NEPSE index)
// Provides a single-number view of the entire market's performance
type MarketIndex struct {
	IndexValue     decimal.Decimal `json:"index_value"`
	Change         decimal.Decimal `json:"change"`
	ChangePercent  decimal.Decimal `json:"change_percent"`
	TotalTurnover  decimal.Decimal `json:"total_turnover"`
	TotalVolume    int64           `json:"total_volume"`
	TotalMarketCap decimal.Decimal `json:"total_market_cap"`
	Advances       int             `json:"advances"`
	Declines       int             `json:"declines"`
	Unchanged      int             `json:"unchanged"`
	TotalCompanies int             `json:"total_companies"`
	PreviousClose  decimal.Decimal `json:"previous_close"`
	Timestamp      time.Time       `json:"timestamp"`
}

// MarketSummary provides a comprehensive market overview combining
// the market index with top gainers, losers, and most active stocks
type MarketSummary struct {
	Index         MarketIndex       `json:"index"`
	TopGainers    []LiveTradingData `json:"top_gainers"`
	TopLosers     []LiveTradingData `json:"top_losers"`
	MostActive    []LiveTradingData `json:"most_active"`
	SectorSummary []SectorIndex     `json:"sector_summary"`
	AsOf          time.Time         `json:"as_of"`
}

// SectorIndex represents aggregated performance of a sector
type SectorIndex struct {
	Sector        string          `json:"sector"`
	Change        decimal.Decimal `json:"change"`
	ChangePercent decimal.Decimal `json:"change_percent"`
	Turnover      decimal.Decimal `json:"turnover"`
	Volume        int64           `json:"volume"`
	CompanyCount  int             `json:"company_count"`
	Advances      int             `json:"advances"`
	Declines      int             `json:"declines"`
}
