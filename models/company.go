package models

import "time"

type Company struct {
	ID          string    `json:"$id,omitempty"` 
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Sector      string    `json:"sector"`
	MarketCap   float64   `json:"market_cap"`
	Description string    `json:"description"`
	FoundedYear int       `json:"founded_year"`
	Employees   int       `json:"employees"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"$createdAt,omitempty"`
	UpdatedAt   time.Time `json:"$updatedAt,omitempty"`
}
