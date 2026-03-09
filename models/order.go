package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Order represents a buy or sell order in the order book
type Order struct {
	ID        string          `json:"id,omitempty"`
	UserID    string          `json:"user_id"`
	CompanyID string          `json:"company_id"`
	Side      string          `json:"side"`       // buy, sell
	OrderType string          `json:"order_type"` // limit, market
	Price     decimal.Decimal `json:"price"`
	Quantity  int64           `json:"quantity"`
	FilledQty int64           `json:"filled_qty"`
	Status    string          `json:"status"` // open, partially_filled, filled, cancelled
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
}

const (
	OrderSideBuy  = "buy"
	OrderSideSell = "sell"

	OrderTypeLimit  = "limit"
	OrderTypeMarket = "market"

	OrderStatusOpen            = "open"
	OrderStatusPartiallyFilled = "partially_filled"
	OrderStatusFilled          = "filled"
	OrderStatusCancelled       = "cancelled"
)

// RemainingQty returns the unfilled quantity
func (o *Order) RemainingQty() int64 {
	return o.Quantity - o.FilledQty
}
