package wallet

import (
	"errors"
	"fmt"
	"time"

	"github.com/919Umesh/gold_go/models"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletLocked        = errors.New("wallet is locked")
	ErrInvalidAmount       = errors.New("invalid amount")
)

type Service interface {
	GetWallet(userID uint) (*models.Wallet, error)
	TopUp(userID uint, amount float64, referenceID string) (*models.Wallet, *models.Transaction, error)
	BuyGold(userID uint, grams, pricePerGram float64, referenceID string) (*models.Wallet, *models.Transaction, error)
	SellGold(userID uint, grams, pricePerGram float64, referenceID string) (*models.Wallet, *models.Transaction, error)
	GetUserTransaction(userID uint) ([]models.Transaction, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetWallet(userID uint) (*models.Wallet, error) {
	wallet, err := s.repo.GetByUserID(userID)
	if err != nil {
		wallet = &models.Wallet{UserID: userID}
		if err := s.repo.Create(wallet); err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}
	}
	return wallet, nil
}

func (s *service) TopUp(userID uint, amount float64, referenceID string) (*models.Wallet, *models.Transaction, error) {
	if amount <= 0 {
		return nil, nil, ErrInvalidAmount
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
			UserID:       userID,
			Type:         models.TransactionTypeTopUp,
			Amount:       amount,
			GoldGrams:    0,
			PricePerGram: 0,
			Status:       models.TransactionStatusSuccess,
			ReferenceID:  referenceID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		return s.repo.CreateTransaction(transaction)
	})
	if err != nil {
		return nil, nil, err
	}
	return updatedWallet, transaction, err
}

func (s *service) BuyGold(userID uint, grams, pricePerGram float64, referenceID string) (*models.Wallet, *models.Transaction, error) {
	if grams <= 0 || pricePerGram <= 0 {
		return nil, nil, ErrInvalidAmount
	}

	totalCost := grams * pricePerGram
	var updatedWallet *models.Wallet
	var transaction *models.Transaction

	err := s.repo.WithLock(userID, func(wallet *models.Wallet) error {
		if wallet.Locked {
			return ErrWalletLocked
		}
		if wallet.FiatBalance < totalCost {
			return ErrInsufficientBalance
		}

		wallet.FiatBalance -= totalCost
		wallet.GoldGrams += grams
		updatedWallet = wallet

		transaction = &models.Transaction{
			UserID:       userID,
			Type:         models.TransactionTypeBuy,
			Amount:       totalCost,
			GoldGrams:    grams,
			PricePerGram: pricePerGram,
			Status:       models.TransactionStatusSuccess,
			ReferenceID:  referenceID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		return s.repo.CreateTransaction(transaction)
	})
	if err != nil {
		return nil, nil, err
	}

	return updatedWallet, transaction, err
}

func (s *service) SellGold(userID uint, grams, pricePerGram float64, referenceID string) (*models.Wallet, *models.Transaction, error) {
	if grams <= 0 || pricePerGram <= 0 {
		return nil, nil, ErrInvalidAmount
	}

	totalValue := grams * pricePerGram
	var updatedWallet *models.Wallet
	var transaction *models.Transaction

	err := s.repo.WithLock(userID, func(wallet *models.Wallet) error {
		if wallet.Locked {
			return ErrWalletLocked
		}
		if wallet.GoldGrams < grams {
			return ErrInsufficientBalance
		}

		wallet.GoldGrams -= grams
		wallet.FiatBalance += totalValue
		updatedWallet = wallet

		transaction = &models.Transaction{
			UserID:       userID,
			Type:         models.TransactionTypeSell,
			Amount:       totalValue,
			GoldGrams:    grams,
			PricePerGram: pricePerGram,
			Status:       models.TransactionStatusSuccess,
			ReferenceID:  referenceID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		return s.repo.CreateTransaction(transaction)
	})
	if err != nil {
		return nil, nil, err
	}

	return updatedWallet, transaction, err
}

func (s *service) GetUserTransaction(userID uint) ([]models.Transaction, error) {

	transaction, err := s.repo.GetUserTransaction(userID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
