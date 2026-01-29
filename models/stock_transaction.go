package models

import (
	"time"

	"gorm.io/gorm"
)

// StockTransactionType represents the type of stock transaction
type StockTransactionType string

const (
	StockTransactionBuy  StockTransactionType = "buy"
	StockTransactionSell StockTransactionType = "sell"
)

// StockTransactionStatus represents the status of a stock transaction
type StockTransactionStatus string

const (
	StockTransactionPending   StockTransactionStatus = "pending"
	StockTransactionCompleted StockTransactionStatus = "completed"
	StockTransactionFailed    StockTransactionStatus = "failed"
	StockTransactionCancelled StockTransactionStatus = "cancelled"
)

// StockTransaction represents a buy or sell transaction
type StockTransaction struct {
	ID            uint                   `gorm:"primaryKey" json:"id"`
	UserID        uint                   `gorm:"index;not null" json:"user_id"`
	CompanyID     uint                   `gorm:"index;not null" json:"company_id"`
	Type          StockTransactionType   `gorm:"size:10;not null" json:"type"`
	Quantity      int                    `gorm:"not null" json:"quantity"`
	PricePerShare float64                `gorm:"type:numeric(10,2);not null" json:"price_per_share"`
	TotalAmount   float64                `gorm:"type:numeric(14,2);not null" json:"total_amount"`
	Status        StockTransactionStatus `gorm:"size:20;default:'completed'" json:"status"`
	ReferenceID   string                 `gorm:"size:100" json:"reference_id"`
	CreatedAt     time.Time              `json:"created_at"`

	// Relationships
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (s *StockTransaction) BeforeCreate(tx *gorm.DB) error {
	s.CreatedAt = time.Now()
	return nil
}
