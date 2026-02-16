package models

import "github.com/shopspring/decimal"

type Wallet struct {
	ID          string          `json:"$id,omitempty"`
	UserID      string          `json:"user_id"`
	FiatBalance decimal.Decimal `json:"fiat_balance"`
	Locked      bool            `json:"locked"`
	Version     int             `json:"version"`
}
