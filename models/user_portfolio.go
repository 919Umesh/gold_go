package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// UserPortfolio represents a user's stock holdings
type UserPortfolio struct {
	ID            string          `json:"$id,omitempty"`
	UserID        string          `json:"user_id"`
	CompanyID     string          `json:"company_id"`
	Quantity      int             `json:"quantity"`
	AvgBuyPrice   decimal.Decimal `json:"average_price"`
	TotalInvested decimal.Decimal `json:"total_invested"`
	CreatedAt     time.Time       `json:"$createdAt,omitempty"`
	UpdatedAt     time.Time       `json:"$updatedAt,omitempty"`

	// Relationships
	User    *User    `json:"user,omitempty"`
	Company *Company `json:"company,omitempty"`
}

func (u *UserPortfolio) CalculateCurrentValue(currentPrice decimal.Decimal) decimal.Decimal {
	return currentPrice.Mul(decimal.NewFromInt(int64(u.Quantity)))
}

func (u *UserPortfolio) CalculateProfitLoss(currentPrice decimal.Decimal) decimal.Decimal {
	currentValue := u.CalculateCurrentValue(currentPrice)
	return currentValue.Sub(u.TotalInvested)
}

func (u *UserPortfolio) CalculateProfitLossPercentage(currentPrice decimal.Decimal) decimal.Decimal {
	if u.TotalInvested.IsZero() {
		return decimal.Zero
	}
	profitLoss := u.CalculateProfitLoss(currentPrice)
	// (profit / invested) * 100
	return profitLoss.Div(u.TotalInvested).Mul(decimal.NewFromInt(100))
}
