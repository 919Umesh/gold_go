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

	if _, err := s.GetWallet(userID); err != nil {
		return nil, nil, fmt.Errorf("failed to ensure wallet: %w", err)
	}

	var updatedWallet *models.Wallet
	var transaction *models.Transaction


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
