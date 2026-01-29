package models

import (
	"time"

	"gorm.io/gorm"
)

// Company represents a demo company in the virtual stock market
type Company struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Symbol      string    `gorm:"size:10;uniqueIndex;not null" json:"symbol"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Sector      string    `gorm:"size:100;not null" json:"sector"`
	MarketCap   float64   `gorm:"type:numeric(20,2)" json:"market_cap"`
	Description string    `gorm:"type:text" json:"description"`
	FoundedYear int       `json:"founded_year"`
	Employees   int       `json:"employees"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	StockPrices []StockPrice `gorm:"foreignKey:CompanyID" json:"-"`
}

func (c *Company) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Company) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
