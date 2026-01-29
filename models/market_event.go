package models

import (
	"time"

	"gorm.io/gorm"
)

// MarketEventType represents different types of market events
type MarketEventType string

const (
	MarketEventEarnings MarketEventType = "earnings"
	MarketEventNews     MarketEventType = "news"
	MarketEventDividend MarketEventType = "dividend"
	MarketEventMerger   MarketEventType = "merger"
	MarketEventIPO      MarketEventType = "ipo"
	MarketEventSplit    MarketEventType = "split"
)

// MarketEvent represents events that can affect stock prices
type MarketEvent struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	CompanyID        uint            `gorm:"index;not null" json:"company_id"`
	EventType        MarketEventType `gorm:"size:50;not null" json:"event_type"`
	Title            string          `gorm:"size:200;not null" json:"title"`
	Description      string          `gorm:"type:text" json:"description"`
	ImpactPercentage float64         `gorm:"type:numeric(5,2)" json:"impact_percentage"` // -10.00 to +10.00
	EventDate        time.Time       `gorm:"not null" json:"event_date"`
	CreatedAt        time.Time       `json:"created_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (m *MarketEvent) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	return nil
}
