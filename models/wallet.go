package models

import (
	"gorm.io/gorm"
)

type Wallet struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UserID      uint    `gorm:"uniqueIndex;not null" json:"user_id"`
	FiatBalance float64 `gorm:"type:numeric(14,2);default:0" json:"fiat_balance"`
	GoldGrams   float64 `gorm:"type:numeric(14,4);default:0" json:"gold_grams"`
	Locked      bool    `gorm:"default:false" json:"locked"`

	// Optimistic locking
	Version int `gorm:"default:1" json:"-"`
}

// SafeUpdate ensures thread-safe wallet updates
func (w *Wallet) SafeUpdate(tx *gorm.DB, updates map[string]interface{}) error {
	return tx.Model(w).Where("version = ?", w.Version).Updates(updates).Error
}
