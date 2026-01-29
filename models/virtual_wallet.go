package models

import (
	"time"

	"gorm.io/gorm"
)

// VirtualWallet represents a user's virtual trading wallet
type VirtualWallet struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance         float64   `gorm:"type:numeric(14,2);default:1000000.00" json:"balance"` // NPR 10 lakh initial
	TotalInvested   float64   `gorm:"type:numeric(14,2);default:0" json:"total_invested"`
	TotalProfitLoss float64   `gorm:"type:numeric(14,2);default:0" json:"total_profit_loss"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (v *VirtualWallet) BeforeCreate(tx *gorm.DB) error {
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	return nil
}

func (v *VirtualWallet) BeforeUpdate(tx *gorm.DB) error {
	v.UpdatedAt = time.Now()
	return nil
}
