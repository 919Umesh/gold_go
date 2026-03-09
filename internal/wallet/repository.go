package wallet

import (
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
)

type Repository interface {
	// Main Wallet
	GetMainWallet(userID string) (*models.MainWallet, error)
	CreateMainWallet(wallet *models.MainWallet) error
	UpdateMainWalletBalance(walletID string, balance string) error

	// Trading Wallet
	GetTradingWallet(userID string) (*models.TradingWallet, error)
	CreateTradingWallet(wallet *models.TradingWallet) error
	UpdateTradingWallet(walletID string, balance string, lockedBalance string) error

	// Wallet Transfers
	CreateWalletTransfer(transfer *models.WalletTransfer) error
	GetWalletTransfers(userID string, limit int) ([]models.WalletTransfer, error)

	// User lookup
	FindUserByID(userID string) (*models.User, error)
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

// ──────────────────── Main Wallet ────────────────────

func (r *repository) GetMainWallet(userID string) (*models.MainWallet, error) {
	var w models.MainWallet
	err := r.client.ExecuteQueryRow("SELECT * FROM main_wallets WHERE user_id = $1", &w, userID)
	if err != nil {
		return nil, fmt.Errorf("main wallet not found: %w", err)
	}
	return &w, nil
}

func (r *repository) CreateMainWallet(wallet *models.MainWallet) error {
	query := `INSERT INTO main_wallets (user_id, balance) VALUES ($1, $2) RETURNING *`
	return r.client.ExecuteInsert(query, wallet, wallet.UserID, wallet.Balance.String())
}

func (r *repository) UpdateMainWalletBalance(walletID string, balance string) error {
	query := `UPDATE main_wallets SET balance = $1 WHERE id = $2 RETURNING *`
	var w models.MainWallet
	return r.client.ExecuteUpdate(query, &w, balance, walletID)
}

// ──────────────────── Trading Wallet ────────────────────

func (r *repository) GetTradingWallet(userID string) (*models.TradingWallet, error) {
	var w models.TradingWallet
	err := r.client.ExecuteQueryRow("SELECT * FROM trading_wallets WHERE user_id = $1", &w, userID)
	if err != nil {
		return nil, fmt.Errorf("trading wallet not found: %w", err)
	}
	return &w, nil
}

func (r *repository) CreateTradingWallet(wallet *models.TradingWallet) error {
	query := `INSERT INTO trading_wallets (user_id, balance, locked_balance) VALUES ($1, $2, $3) RETURNING *`
	return r.client.ExecuteInsert(query, wallet, wallet.UserID, wallet.Balance.String(), wallet.LockedBalance.String())
}

func (r *repository) UpdateTradingWallet(walletID string, balance string, lockedBalance string) error {
	query := `UPDATE trading_wallets SET balance = $1, locked_balance = $2 WHERE id = $3 RETURNING *`
	var w models.TradingWallet
	return r.client.ExecuteUpdate(query, &w, balance, lockedBalance, walletID)
}

// ──────────────────── Wallet Transfers ────────────────────

func (r *repository) CreateWalletTransfer(transfer *models.WalletTransfer) error {
	query := `INSERT INTO wallet_transfers (user_id, amount, direction, status) VALUES ($1, $2, $3, $4) RETURNING *`
	return r.client.ExecuteInsert(query, transfer,
		transfer.UserID, transfer.Amount.String(), transfer.Direction, transfer.Status)
}

func (r *repository) GetWalletTransfers(userID string, limit int) ([]models.WalletTransfer, error) {
	var transfers []models.WalletTransfer
	query := "SELECT * FROM wallet_transfers WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2"
	err := r.client.ExecuteQuery(query, &transfers, userID, limit)
	if err != nil {
		return nil, err
	}
	if transfers == nil {
		transfers = []models.WalletTransfer{}
	}
	return transfers, nil
}

// ──────────────────── User ────────────────────

func (r *repository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	err := r.client.ExecuteQueryRow("SELECT * FROM users WHERE id = $1", &user, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}
