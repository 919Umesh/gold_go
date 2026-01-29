package models

import (
	"time"

	"gorm.io/gorm"
)

// UserPortfolio represents a user's stock holdings
type UserPortfolio struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex:idx_user_company;not null" json:"user_id"`
	CompanyID     uint      `gorm:"uniqueIndex:idx_user_company;not null" json:"company_id"`
	Quantity      int       `gorm:"not null;default:0" json:"quantity"`
	AvgBuyPrice   float64   `gorm:"type:numeric(10,2);not null" json:"avg_buy_price"`
	TotalInvested float64   `gorm:"type:numeric(14,2);not null" json:"total_invested"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (u *UserPortfolio) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *UserPortfolio) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// CalculateCurrentValue calculates the current value of the holding
func (u *UserPortfolio) CalculateCurrentValue(currentPrice float64) float64 {
	return float64(u.Quantity) * currentPrice
}

// CalculateProfitLoss calculates the profit or loss
func (u *UserPortfolio) CalculateProfitLoss(currentPrice float64) float64 {
	currentValue := u.CalculateCurrentValue(currentPrice)
	return currentValue - u.TotalInvested
}

// CalculateProfitLossPercentage calculates the profit or loss percentage
func (u *UserPortfolio) CalculateProfitLossPercentage(currentPrice float64) float64 {
	if u.TotalInvested == 0 {
		return 0
	}
	profitLoss := u.CalculateProfitLoss(currentPrice)
	return (profitLoss / u.TotalInvested) * 100
}
