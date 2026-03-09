package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// IPO represents an Initial Public Offering event
type IPO struct {
	ID              string          `json:"id,omitempty"`
	CompanyID       string          `json:"company_id"`
	PricePerShare   decimal.Decimal `json:"price_per_share"`
	TotalShares     int64           `json:"total_shares"`
	AllocatedShares int64           `json:"allocated_shares"`
	MaxPerApplicant int64           `json:"max_per_applicant"`
	OpenAt          time.Time       `json:"open_at"`
	CloseAt         time.Time       `json:"close_at"`
	Status          string          `json:"status"` // pending, open, closed, allocated
	CreatedAt       time.Time       `json:"created_at,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at,omitempty"`
}

const (
	IPOStatusPending   = "pending"
	IPOStatusOpen      = "open"
	IPOStatusClosed    = "closed"
	IPOStatusAllocated = "allocated"
)

// IPOApplication represents a user's application for an IPO
type IPOApplication struct {
	ID              string          `json:"id,omitempty"`
	IPOID           string          `json:"ipo_id"`
	UserID          string          `json:"user_id"`
	SharesRequested int64           `json:"shares_requested"`
	SharesAllocated int64           `json:"shares_allocated"`
	AmountPaid      decimal.Decimal `json:"amount_paid"`
	AmountRefunded  decimal.Decimal `json:"amount_refunded"`
	Status          string          `json:"status"` // pending, allocated, not_allocated, refunded
	CreatedAt       time.Time       `json:"created_at,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at,omitempty"`
}

const (
	IPOAppStatusPending      = "pending"
	IPOAppStatusAllocated    = "allocated"
	IPOAppStatusNotAllocated = "not_allocated"
	IPOAppStatusRefunded     = "refunded"
)
