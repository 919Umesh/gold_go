package models

import (
	"time"
)

// UserPortfolio represents a user's stock holdings
type UserPortfolio struct {
	ID            string    `json:"$id,omitempty"`
	UserID        string    `json:"user_id"`
	CompanyID     string    `json:"company_id"`
	Quantity      int       `json:"quantity"`
	AvgBuyPrice   float64   `json:"average_price"`
	TotalInvested float64   `json:"total_invested"`
	CreatedAt     time.Time `json:"$createdAt,omitempty"`
	UpdatedAt     time.Time `json:"$updatedAt,omitempty"`

	// Relationships
	User    *User    `json:"user,omitempty"`
	Company *Company `json:"company,omitempty"`
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
