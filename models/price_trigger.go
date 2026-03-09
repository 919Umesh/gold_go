package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// PriceTrigger represents an auto-sell trigger
type PriceTrigger struct {
	ID           string          `json:"id,omitempty"`
	UserID       string          `json:"user_id"`
	CompanyID    string          `json:"company_id"`
	TriggerPrice decimal.Decimal `json:"trigger_price"`
	SharesQty    int64           `json:"shares_qty"`
	Direction    string          `json:"direction"` // above, below
	Status       string          `json:"status"`    // active, triggered, cancelled
	CreatedAt    time.Time       `json:"created_at,omitempty"`
	UpdatedAt    time.Time       `json:"updated_at,omitempty"`
}

const (
	TriggerDirectionAbove = "above"
	TriggerDirectionBelow = "below"

	TriggerStatusActive    = "active"
	TriggerStatusTriggered = "triggered"
	TriggerStatusCancelled = "cancelled"
)
