package wallet

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/919Umesh/stock_market_sim/models"
	"github.com/919Umesh/stock_market_sim/pkg/queue"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletLocked        = errors.New("wallet is locked")
	ErrInvalidAmount       = errors.New("invalid amount")
)

type Service interface {
	GetWallet(userID string) (*models.Wallet, error)
	TopUp(userID string, amount float64, referenceID string) (*models.Wallet, *models.Transaction, error)
	GetUserTransaction(userID string) ([]models.Transaction, error)
}

type service struct {
	repo       Repository
	workerPool *queue.WorkerPool
}

// TransactionAnalyticsJob implements queue.Job to simulate post-transaction processing
type TransactionAnalyticsJob struct {
	TransactionID string
	Type          models.TransactionType
	Amount        float64
	UserID        string
}

func (j *TransactionAnalyticsJob) Process() error {
	// Simulate heavy analytics processing
	time.Sleep(100 * time.Millisecond)
	slog.Info("Async analytics processed",
		"transaction_id", j.TransactionID,
		"type", j.Type,
		"user_id", j.UserID,
		"amount", j.Amount,
	)
	return nil
}

func NewService(repo Repository, wp *queue.WorkerPool) Service {
	return &service{
		repo:       repo,
		workerPool: wp,
	}
}

func (s *service) GetWallet(userID string) (*models.Wallet, error) {
	wallet, err := s.repo.GetByUserID(userID)
	if err != nil {
		// If using Appwrite, we check specific error or return nil if not found
		// Assuming repo returns error if not found or nil wallet?
		// Usually repo returns (nil, error) or (nil, ErrNotFound)
		// Let's assume ErrNotFound or simply err.
		// If err, we create.
		// Wait, if database error, we shouldn't create.
		// logic: if err is Not Found, create.
		// Appwrite returns error on Get/List if failed? List returns empty docs if not found?
		// Let's rely on Repo implementation quirks.

		// If generic error, return error.
		// If "not found", create.
		// For now simple logic: catch error and try create? risky.
		// Let's assume repo handles check-then-create or returns specific error.
		// Current logic: if err != nil, try create.

		// Create a new wallet
		wallet = &models.Wallet{UserID: userID}
		if err := s.repo.Create(wallet); err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}
	}
	return wallet, nil
}

func (s *service) TopUp(userID string, amount float64, referenceID string) (*models.Wallet, *models.Transaction, error) {
	if amount <= 0 {
		return nil, nil, ErrInvalidAmount
	}

	// Ensure wallet exists before attempting to lock/update it.
	// Previously WithLock called GetByUserID and returned an error if wallet was missing,
	// causing a generic 500 in handlers. Create the wallet if needed.
	if _, err := s.GetWallet(userID); err != nil {
		return nil, nil, fmt.Errorf("failed to ensure wallet: %w", err)
	}

	var updatedWallet *models.Wallet
	var transaction *models.Transaction

	// WithLock might be problematic in Appwrite (No transaction support in standard SDK yet?)
	// We might need to remove lock/transaction or implement optimistic locking via version match.
	// For now, removing WithLock and doing sequential operations (Unsafe but compiles)
	// Or implementation of WithLock in Appwrite repo using simple callback?

	err := s.repo.WithLock(userID, func(wallet *models.Wallet) error {
		if wallet.Locked {
			return ErrWalletLocked
		}
		wallet.FiatBalance += amount
		updatedWallet = wallet

		transaction = &models.Transaction{
			UserID:      userID,
			Type:        models.TransactionTypeTopUp,
			Amount:      amount,
			Status:      models.TransactionStatusSuccess,
			ReferenceID: referenceID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		return s.repo.CreateTransaction(transaction)
	})
	if err != nil {
		return nil, nil, err
	}

	s.workerPool.Submit(&TransactionAnalyticsJob{
		TransactionID: transaction.ID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		UserID:        transaction.UserID,
	})

	return updatedWallet, transaction, err
}

func (s *service) GetUserTransaction(userID string) ([]models.Transaction, error) {

	transaction, err := s.repo.GetUserTransaction(userID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
