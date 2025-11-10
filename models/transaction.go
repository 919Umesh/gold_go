package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeBuy    TransactionType = "buy"
	TransactionTypeSell   TransactionType = "sell"
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
	ID           uint              `gorm:"primaryKey" json:"id"`
	UserID       uint              `gorm:"index;not null" json:"user_id"`
	Type         TransactionType   `gorm:"size:20;not null" json:"type"`
	Amount       float64           `gorm:"type:numeric(14,2)" json:"amount"`
	GoldGrams    float64           `gorm:"type:numeric(14,4)" json:"gold_grams"`
	PricePerGram float64           `gorm:"type:numeric(10,4)" json:"price_per_gram"`
	Status       TransactionStatus `gorm:"size:20;default:pending" json:"status"`
	ReferenceID  string            `gorm:"size:100" json:"reference_id"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}
