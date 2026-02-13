package models

import (
	"time"
)

type StockPrice struct {
	ID         string    `json:"$id,omitempty"`
	CompanyID  string    `json:"company_id"`
	OpenPrice  float64   `json:"open_price"`
	HighPrice  float64   `json:"high_price"`
	LowPrice   float64   `json:"low_price"`
	ClosePrice float64   `json:"close_price"`
	Volume     int64     `json:"volume"`
	Timestamp  time.Time `json:"timestamp"`
	Timeframe  string    `json:"timeframe"`
	CreatedAt  time.Time `json:"$createdAt,omitempty"`


	Company *Company `json:"company,omitempty"`
}


func (s *StockPrice) CalculateChange() float64 {
	if s.OpenPrice == 0 {
		return 0
	}
	return ((s.ClosePrice - s.OpenPrice) / s.OpenPrice) * 100
}
