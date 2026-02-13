package models

import (
	"time"
)


type StockTransactionType string

const (
	StockTransactionBuy  StockTransactionType = "buy"
	StockTransactionSell StockTransactionType = "sell"
)

type StockTransactionStatus string

const (
	StockTransactionPending   StockTransactionStatus = "pending"
	StockTransactionCompleted StockTransactionStatus = "completed"
	StockTransactionFailed    StockTransactionStatus = "failed"
	StockTransactionCancelled StockTransactionStatus = "cancelled"
)


type StockTransaction struct {
	ID            string                 `json:"$id,omitempty"`
	UserID        string                 `json:"user_id"`
	CompanyID     string                 `json:"company_id"`
	Type          StockTransactionType   `json:"type"`
	Quantity      int                    `json:"quantity"`
	PricePerShare float64                `json:"price_per_share"`
	TotalAmount   float64                `json:"total_amount"`
	Status        StockTransactionStatus `json:"status"`
	ReferenceID   string                 `json:"reference_id"`
	CreatedAt     time.Time              `json:"$createdAt,omitempty"`

	User    *User    `json:"user,omitempty"`
	Company *Company `json:"company,omitempty"`
}
