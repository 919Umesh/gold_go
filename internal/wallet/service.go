package wallet

import (
	"errors"
	"fmt"

	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

var (
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrInvalidAmount         = errors.New("amount must be positive")
	ErrInsufficientAvailable = errors.New("insufficient available balance (funds locked in orders)")
)

type Service interface {
	// Wallet retrieval (auto-creates if not found)
	GetMainWallet(userID string) (*models.MainWallet, error)
	GetTradingWallet(userID string) (*models.TradingWallet, error)
	GetBothWallets(userID string) (*models.MainWallet, *models.TradingWallet, error)

	// Top up main wallet
	TopUp(userID string, amount decimal.Decimal) (*models.MainWallet, error)

	// Transfer between wallets (fixes the critical bug)
	TransferToTrading(userID string, amount decimal.Decimal) (*models.MainWallet, *models.TradingWallet, error)
	TransferToMain(userID string, amount decimal.Decimal) (*models.MainWallet, *models.TradingWallet, error)

	// Fund locking for orders
	LockFunds(userID string, amount decimal.Decimal) error
	ReleaseFunds(userID string, amount decimal.Decimal) error
	DeductLockedFunds(userID string, amount decimal.Decimal) error

	// Credit trading wallet (e.g., from sell proceeds)
	CreditTradingWallet(userID string, amount decimal.Decimal) error

	// Transfer history
	GetTransferHistory(userID string, limit int) ([]models.WalletTransfer, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ──────────────────── Wallet Retrieval ────────────────────

func (s *service) GetMainWallet(userID string) (*models.MainWallet, error) {
	w, err := s.repo.GetMainWallet(userID)
	if err != nil {
		// Auto-create
		w = &models.MainWallet{
			UserID:  userID,
			Balance: decimal.Zero,
		}
		if createErr := s.repo.CreateMainWallet(w); createErr != nil {
			return nil, fmt.Errorf("failed to create main wallet: %w", createErr)
		}
		// Re-fetch to get ID
		w, err = s.repo.GetMainWallet(userID)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (s *service) GetTradingWallet(userID string) (*models.TradingWallet, error) {
	w, err := s.repo.GetTradingWallet(userID)
	if err != nil {
		w = &models.TradingWallet{
			UserID:        userID,
			Balance:       decimal.Zero,
			LockedBalance: decimal.Zero,
		}
		if createErr := s.repo.CreateTradingWallet(w); createErr != nil {
			return nil, fmt.Errorf("failed to create trading wallet: %w", createErr)
		}
		w, err = s.repo.GetTradingWallet(userID)
		if err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (s *service) GetBothWallets(userID string) (*models.MainWallet, *models.TradingWallet, error) {
	main, err := s.GetMainWallet(userID)
	if err != nil {
		return nil, nil, err
	}
	trading, err := s.GetTradingWallet(userID)
	if err != nil {
		return nil, nil, err
	}
	return main, trading, nil
}

// ──────────────────── Top Up ────────────────────

func (s *service) TopUp(userID string, amount decimal.Decimal) (*models.MainWallet, error) {
	if !amount.IsPositive() {
		return nil, ErrInvalidAmount
	}

	main, err := s.GetMainWallet(userID)
	if err != nil {
		return nil, err
	}

	newBalance := main.Balance.Add(amount)
	if err := s.repo.UpdateMainWalletBalance(main.ID, newBalance.String()); err != nil {
		return nil, fmt.Errorf("failed to update main wallet: %w", err)
	}

	main.Balance = newBalance
	return main, nil
}

// ──────────────────── Transfers (Critical Bug Fix) ────────────────────

func (s *service) TransferToTrading(userID string, amount decimal.Decimal) (*models.MainWallet, *models.TradingWallet, error) {
	if !amount.IsPositive() {
		return nil, nil, ErrInvalidAmount
	}

	main, err := s.GetMainWallet(userID)
	if err != nil {
		return nil, nil, err
	}

	if main.Balance.LessThan(amount) {
		return nil, nil, ErrInsufficientBalance
	}

	trading, err := s.GetTradingWallet(userID)
	if err != nil {
		return nil, nil, err
	}

	// Debit main wallet
	newMainBalance := main.Balance.Sub(amount)
	if err := s.repo.UpdateMainWalletBalance(main.ID, newMainBalance.String()); err != nil {
		return nil, nil, fmt.Errorf("failed to debit main wallet: %w", err)
	}

	// Credit trading wallet
	newTradingBalance := trading.Balance.Add(amount)
	if err := s.repo.UpdateTradingWallet(trading.ID, newTradingBalance.String(), trading.LockedBalance.String()); err != nil {
		// Rollback main wallet
		_ = s.repo.UpdateMainWalletBalance(main.ID, main.Balance.String())
		return nil, nil, fmt.Errorf("failed to credit trading wallet: %w", err)
	}

	// Log the transfer
	transfer := &models.WalletTransfer{
		UserID:    userID,
		Amount:    amount,
		Direction: models.TransferMainToTrading,
		Status:    "completed",
	}
	_ = s.repo.CreateWalletTransfer(transfer)

	main.Balance = newMainBalance
	trading.Balance = newTradingBalance
	return main, trading, nil
}

func (s *service) TransferToMain(userID string, amount decimal.Decimal) (*models.MainWallet, *models.TradingWallet, error) {
	if !amount.IsPositive() {
		return nil, nil, ErrInvalidAmount
	}

	trading, err := s.GetTradingWallet(userID)
	if err != nil {
		return nil, nil, err
	}

	available := trading.AvailableBalance()
	if available.LessThan(amount) {
		return nil, nil, ErrInsufficientAvailable
	}

	main, err := s.GetMainWallet(userID)
	if err != nil {
		return nil, nil, err
	}

	// Debit trading wallet
	newTradingBalance := trading.Balance.Sub(amount)
	if err := s.repo.UpdateTradingWallet(trading.ID, newTradingBalance.String(), trading.LockedBalance.String()); err != nil {
		return nil, nil, fmt.Errorf("failed to debit trading wallet: %w", err)
	}

	// Credit main wallet
	newMainBalance := main.Balance.Add(amount)
	if err := s.repo.UpdateMainWalletBalance(main.ID, newMainBalance.String()); err != nil {
		_ = s.repo.UpdateTradingWallet(trading.ID, trading.Balance.String(), trading.LockedBalance.String())
		return nil, nil, fmt.Errorf("failed to credit main wallet: %w", err)
	}

	transfer := &models.WalletTransfer{
		UserID:    userID,
		Amount:    amount,
		Direction: models.TransferTradingToMain,
		Status:    "completed",
	}
	_ = s.repo.CreateWalletTransfer(transfer)

	main.Balance = newMainBalance
	trading.Balance = newTradingBalance
	return main, trading, nil
}

// ──────────────────── Fund Locking for Orders ────────────────────

func (s *service) LockFunds(userID string, amount decimal.Decimal) error {
	trading, err := s.GetTradingWallet(userID)
	if err != nil {
		return err
	}

	available := trading.AvailableBalance()
	if available.LessThan(amount) {
		return ErrInsufficientAvailable
	}

	newLocked := trading.LockedBalance.Add(amount)
	return s.repo.UpdateTradingWallet(trading.ID, trading.Balance.String(), newLocked.String())
}

func (s *service) ReleaseFunds(userID string, amount decimal.Decimal) error {
	trading, err := s.GetTradingWallet(userID)
	if err != nil {
		return err
	}

	newLocked := trading.LockedBalance.Sub(amount)
	if newLocked.IsNegative() {
		newLocked = decimal.Zero
	}
	return s.repo.UpdateTradingWallet(trading.ID, trading.Balance.String(), newLocked.String())
}

func (s *service) DeductLockedFunds(userID string, amount decimal.Decimal) error {
	trading, err := s.GetTradingWallet(userID)
	if err != nil {
		return err
	}

	newBalance := trading.Balance.Sub(amount)
	newLocked := trading.LockedBalance.Sub(amount)
	if newLocked.IsNegative() {
		newLocked = decimal.Zero
	}
	return s.repo.UpdateTradingWallet(trading.ID, newBalance.String(), newLocked.String())
}

func (s *service) CreditTradingWallet(userID string, amount decimal.Decimal) error {
	trading, err := s.GetTradingWallet(userID)
	if err != nil {
		return err
	}

	newBalance := trading.Balance.Add(amount)
	return s.repo.UpdateTradingWallet(trading.ID, newBalance.String(), trading.LockedBalance.String())
}

// ──────────────────── Transfer History ────────────────────

func (s *service) GetTransferHistory(userID string, limit int) ([]models.WalletTransfer, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetWalletTransfers(userID, limit)
}
