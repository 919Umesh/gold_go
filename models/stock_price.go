package models

import (
	"time"

	"gorm.io/gorm"
)

// StockPrice represents OHLCV (Open, High, Low, Close, Volume) data for a stock
type StockPrice struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CompanyID  uint      `gorm:"index:idx_company_timestamp;not null" json:"company_id"`
	OpenPrice  float64   `gorm:"type:numeric(10,2);not null" json:"open_price"`
	HighPrice  float64   `gorm:"type:numeric(10,2);not null" json:"high_price"`
	LowPrice   float64   `gorm:"type:numeric(10,2);not null" json:"low_price"`
	ClosePrice float64   `gorm:"type:numeric(10,2);not null" json:"close_price"`
	Volume     int64     `gorm:"default:0" json:"volume"`
	Timestamp  time.Time `gorm:"index:idx_company_timestamp;index:idx_timeframe;not null" json:"timestamp"`
	Timeframe  string    `gorm:"size:20;default:'1d';index:idx_timeframe" json:"timeframe"` // 1m, 5m, 15m, 1h, 1d, 1w, 1M
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (s *StockPrice) BeforeCreate(tx *gorm.DB) error {
	s.CreatedAt = time.Now()
	return nil
}

// CalculateChange returns the price change percentage
func (s *StockPrice) CalculateChange() float64 {
	if s.OpenPrice == 0 {
		return 0
	}
	return ((s.ClosePrice - s.OpenPrice) / s.OpenPrice) * 100
}
