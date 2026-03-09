package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Trade represents a matched trade between a buy and sell order
type Trade struct {
	ID          string          `json:"id,omitempty"`
	CompanyID   string          `json:"company_id"`
	BuyOrderID  string          `json:"buy_order_id"`
	SellOrderID string          `json:"sell_order_id"`
	BuyerID     string          `json:"buyer_id"`
	SellerID    string          `json:"seller_id"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int64           `json:"quantity"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
}
