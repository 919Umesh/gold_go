package models

import (
	"time"
)

// CompanyEvent represents a corporate event (AGM, dividend, etc.)
type CompanyEvent struct {
	ID          string    `json:"id,omitempty"`
	CompanyID   string    `json:"company_id"`
	EventType   string    `json:"event_type"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	EventDate   time.Time `json:"event_date"`
	FiscalYear  string    `json:"fiscal_year,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

const (
	EventTypeAGM               = "agm"
	EventTypeDividend          = "dividend"
	EventTypeBonusShare        = "bonus_share"
	EventTypeRightShare        = "right_share"
	EventTypeQuarterlyReport   = "quarterly_report"
	EventTypeBoardMeeting      = "board_meeting"
	EventTypeFinancialResults  = "financial_results"
	EventTypeStockSplit        = "stock_split"
	EventTypeMergerAcquisition = "merger_acquisition"
	EventTypeIPOAnnouncement   = "ipo_announcement"

	EventStatusUpcoming  = "upcoming"
	EventStatusOngoing   = "ongoing"
	EventStatusCompleted = "completed"
	EventStatusCancelled = "cancelled"
)
