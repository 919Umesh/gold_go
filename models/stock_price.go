package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type StockPrice struct {
	ID            string          `json:"id,omitempty"`
	CompanyID     string          `json:"company_id"`
	OpenPrice     decimal.Decimal `json:"open_price"`
	HighPrice     decimal.Decimal `json:"high_price"`
	LowPrice      decimal.Decimal `json:"low_price"`
	ClosePrice    decimal.Decimal `json:"close_price"`
	Volume        int64           `json:"volume"`
	Turnover      decimal.Decimal `json:"turnover"`
	ChangePercent decimal.Decimal `json:"change_percent"`
	Timestamp     time.Time       `json:"timestamp"`
	Timeframe     string          `json:"timeframe"`
	CreatedAt     time.Time       `json:"created_at,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
}

// CandlestickData is the frontend-friendly OHLCV representation
type CandlestickData struct {
	Timestamp     time.Time       `json:"timestamp"`
	Open          decimal.Decimal `json:"open"`
	High          decimal.Decimal `json:"high"`
	Low           decimal.Decimal `json:"low"`
	Close         decimal.Decimal `json:"close"`
	Volume        int64           `json:"volume"`
	Turnover      decimal.Decimal `json:"turnover"`
	ChangePercent decimal.Decimal `json:"change_percent"`
}

// LiveTradingData represents real-time trading info for a company
type LiveTradingData struct {
	Symbol        string          `json:"symbol"`
	CompanyID     string          `json:"company_id"`
	CompanyName   string          `json:"company_name"`
	Sector        string          `json:"sector,omitempty"`
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

// TradeFeedItem represents a single trade in the live feed
type TradeFeedItem struct {
	Symbol      string          `json:"symbol"`
	CompanyName string          `json:"company_name"`
	TradeType   string          `json:"trade_type"`
	Quantity    int64           `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	PriceImpact decimal.Decimal `json:"price_impact"`
	NewPrice    decimal.Decimal `json:"new_price"`
	Timestamp   time.Time       `json:"timestamp"`
}

// MarketIndex represents the overall market index
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

// MarketSummary combines index, gainers, losers, most active
type MarketSummary struct {
	Index      MarketIndex       `json:"index"`
	TopGainers []LiveTradingData `json:"top_gainers"`
	TopLosers  []LiveTradingData `json:"top_losers"`
	MostActive []LiveTradingData `json:"most_active"`
	AsOf       time.Time         `json:"as_of"`
}

// SectorPerformance represents sector-level aggregation
type SectorPerformance struct {
	Sector         string          `json:"sector"`
	CompanyCount   int             `json:"company_count"`
	AvgChange      decimal.Decimal `json:"avg_change_percent"`
	TotalTurnover  decimal.Decimal `json:"total_turnover"`
	TotalVolume    int64           `json:"total_volume"`
	TotalMarketCap decimal.Decimal `json:"total_market_cap"`
}
