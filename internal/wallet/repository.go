package wallet

import (
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
)

const (
	TableWallets      = "virtual_wallets"
	TableTransactions = "transactions"
	TableUsers        = "users"
)

type Repository interface {
	FindUserByID(userID string) (*models.User, error)
	GetByUserID(userID string) (*models.Wallet, error)
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet) error
	WithLock(userID string, fn func(*models.Wallet) error) error

	CreateTransaction(transaction *models.Transaction) error
	UpdateTransaction(transaction *models.Transaction) error

	GetUserTransaction(userID string) ([]models.Transaction, error)
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

func (r *repository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"
	err := r.client.ExecuteQueryRow(query, &user, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (r *repository) GetByUserID(userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	query := "SELECT * FROM virtual_wallets WHERE user_id = $1"
	err := r.client.ExecuteQueryRow(query, &wallet, userID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found")
	}
	return &wallet, nil
}

func (r *repository) Create(wallet *models.Wallet) error {
	query := `INSERT INTO virtual_wallets (user_id, balance, total_invested, total_profit_loss, fiat_balance, locked, version)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
	return r.client.ExecuteInsert(query, wallet,
		wallet.UserID, 0, 0, 0, wallet.FiatBalance.InexactFloat64(), wallet.Locked, wallet.Version)
}


func (r *repository) Update(wallet *models.Wallet) error {
	query := `UPDATE virtual_wallets SET fiat_balance = $1, locked = $2, version = $3
			  WHERE id = $4 RETURNING *`
	return r.client.ExecuteUpdate(query, wallet,
		wallet.FiatBalance.InexactFloat64(), wallet.Locked, wallet.Version, wallet.ID)
}


func (r *repository) CreateTransaction(transaction *models.Transaction) error {
	query := `INSERT INTO transactions (user_id, type, amount, status, reference_id)
			  VALUES ($1, $2, $3, $4, $5) RETURNING *`
	return r.client.ExecuteInsert(query, transaction,
		transaction.UserID, string(transaction.Type), transaction.Amount.InexactFloat64(),
		string(transaction.Status), transaction.ReferenceID)
}


func (r *repository) UpdateTransaction(transaction *models.Transaction) error {
	query := "UPDATE transactions SET status = $1 WHERE id = $2 RETURNING *"
	return r.client.ExecuteUpdate(query, transaction,
		string(transaction.Status), transaction.ID)
}


func (r *repository) GetUserTransaction(userID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := "SELECT * FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT 100"
	err := r.client.ExecuteQuery(query, &transactions, userID)
	if err != nil {
		return nil, err
	}
	if transactions == nil {
		transactions = []models.Transaction{}
	}
	return transactions, nil
}

func (r *repository) WithLock(userID string, fn func(*models.Wallet) error) error {
	wallet, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	if wallet.Locked {
		return ErrWalletLocked
	}

	wallet.Locked = true
	if err := r.Update(wallet); err != nil {
		return err
	}

	callbackErr := fn(wallet)

	wallet.Locked = false
	if updateErr := r.Update(wallet); updateErr != nil {
		if callbackErr == nil {
			return updateErr
		}
	}

	return callbackErr
}
