package models

import (
	"time"

	"gorm.io/gorm"
)

type GoldPrice struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PricePerGram float64   `gorm:"type:numeric(10,4);not null" json:"price_per_gram"`
	Source       string    `gorm:"size:50" json:"source"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (g *GoldPrice) BeforeCreate(tx *gorm.DB) error {
	g.UpdatedAt = time.Now()
	return nil
}
