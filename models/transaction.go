package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeTopUp  TransactionType = "topup"
	TransactionTypeRefund TransactionType = "refund"

	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

// Transaction represents wallet top-up / refund transactions
type Transaction struct {
	ID          string            `json:"id,omitempty"`
	UserID      string            `json:"user_id"`
	Type        TransactionType   `json:"type"`
	Amount      decimal.Decimal   `json:"amount"`
	Status      TransactionStatus `json:"status"`
	ReferenceID string            `json:"reference_id"`
	CreatedAt   time.Time         `json:"created_at,omitempty"`
	UpdatedAt   time.Time         `json:"updated_at,omitempty"`
}
