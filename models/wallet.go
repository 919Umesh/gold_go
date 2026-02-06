package models

type Wallet struct {
	ID          string  `json:"$id,omitempty"`
	UserID      string  `json:"user_id"`
	FiatBalance float64 `json:"fiat_balance"`
	Locked      bool    `json:"locked"`
	Version     int     `json:"version"` // For optimistic locking
}
