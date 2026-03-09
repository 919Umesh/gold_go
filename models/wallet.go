package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// MainWallet holds user's deposited funds (real-money equivalent)
type MainWallet struct {
	ID        string          `json:"id,omitempty"`
	UserID    string          `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
}

// TradingWallet holds funds used for buying shares and receiving sale proceeds
type TradingWallet struct {
	ID            string          `json:"id,omitempty"`
	UserID        string          `json:"user_id"`
	Balance       decimal.Decimal `json:"balance"`
	LockedBalance decimal.Decimal `json:"locked_balance"` // reserved for pending buy orders
	CreatedAt     time.Time       `json:"created_at,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
}

// AvailableBalance returns the balance minus locked funds
func (tw *TradingWallet) AvailableBalance() decimal.Decimal {
	return tw.Balance.Sub(tw.LockedBalance)
}

// WalletTransfer represents a fund transfer between main and trading wallets
type WalletTransfer struct {
	ID        string          `json:"id,omitempty"`
	UserID    string          `json:"user_id"`
	Amount    decimal.Decimal `json:"amount"`
	Direction string          `json:"direction"` // main_to_trading, trading_to_main
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
}

const (
	TransferMainToTrading = "main_to_trading"
	TransferTradingToMain = "trading_to_main"
)
