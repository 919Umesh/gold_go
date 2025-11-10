package wallet

import (
	"errors"
	"fmt"

	"github.com/umesh/gold_investment/models"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletLocked        = errors.New("wallet is locked")
	ErrInvalidAmount       = errors.New("invalid amount")
)

type Service interface {
	GetWallet(userID uint) (*models.Wallet, error)
	TopUp(userID uint, amount float64) (*models.Wallet, error)
	BuyGold(userID uint, grams, pricePerGram float64) (*models.Wallet, error)
	SellGold(userID uint, grams, pricePerGram float64) (*models.Wallet, error)
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
		// Create wallet if not exists
		wallet = &models.Wallet{UserID: userID}
		if err := s.repo.Create(wallet); err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}
	}
	return wallet, nil
}

func (s *service) TopUp(userID uint, amount float64) (*models.Wallet, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	var updatedWallet *models.Wallet
	err := s.repo.WithLock(userID, func(wallet *models.Wallet) error {
		if wallet.Locked {
			return ErrWalletLocked
		}
		wallet.FiatBalance += amount
		updatedWallet = wallet
		return nil
	})

	return updatedWallet, err
}

func (s *service) BuyGold(userID uint, grams, pricePerGram float64) (*models.Wallet, error) {
	if grams <= 0 || pricePerGram <= 0 {
		return nil, ErrInvalidAmount
	}

	totalCost := grams * pricePerGram
	var updatedWallet *models.Wallet

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
		return nil
	})

	return updatedWallet, err
}

func (s *service) SellGold(userID uint, grams, pricePerGram float64) (*models.Wallet, error) {
	if grams <= 0 || pricePerGram <= 0 {
		return nil, ErrInvalidAmount
	}

	totalValue := grams * pricePerGram
	var updatedWallet *models.Wallet

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
		return nil
	})

	return updatedWallet, err
}
