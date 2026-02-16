package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeTopUp  TransactionType = "topup"
	TransactionTypeRefund TransactionType = "refund"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	ID          string            `json:"$id,omitempty"`
	UserID      string            `json:"user_id"`
	Type        TransactionType   `json:"type"`
	Amount      decimal.Decimal   `json:"amount"`
	Status      TransactionStatus `json:"status"`
	ReferenceID string            `json:"reference_id"`
	CreatedAt   time.Time         `json:"$createdAt,omitempty"`
	UpdatedAt   time.Time         `json:"$updatedAt,omitempty"`

	User *User `json:"user,omitempty"`
}
