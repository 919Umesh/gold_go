package models

import (
	"time"
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
	ID               string          `json:"$id,omitempty"`
	CompanyID        string          `json:"company_id"`
	EventType        MarketEventType `json:"event_type"`
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	ImpactPercentage float64         `json:"impact_percentage"`
	EventDate        time.Time       `json:"event_date"`
	CreatedAt        time.Time       `json:"$createdAt,omitempty"`

	Company *Company `json:"company,omitempty"`
}
