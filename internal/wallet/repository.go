package wallet

import (
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/query"
)

const (
	CollectionWallets      = "virtual_wallets"
	CollectionTransactions = "transactions" 
)

type Repository interface {
	GetByUserID(userID string) (*models.Wallet, error)
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet) error
	WithLock(userID string, fn func(*models.Wallet) error) error

	CreateTransaction(transaction *models.Transaction) error
	UpdateTransaction(transaction *models.Transaction) error

	GetUserTransaction(userID string) ([]models.Transaction, error)
}

type repository struct {
	client *appwrite.Client
}

func NewRepository(client *appwrite.Client) Repository {
	return &repository{client: client}
}

func (r *repository) GetByUserID(userID string) (*models.Wallet, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionWallets,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("user_id", userID),
			query.Limit(1),
		}),
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Documents) == 0 {
		return nil, fmt.Errorf("wallet not found")
	}

	var wallet models.Wallet
	if err := appwrite.DecodeListItem(resp, 0, &wallet); err != nil {
		return nil, fmt.Errorf("failed to decode wallet: %w", err)
	}
	return &wallet, nil
}

func (r *repository) Create(wallet *models.Wallet) error {
	data := map[string]interface{}{
		"user_id":           wallet.UserID,
		"balance":           0.0,
		"total_invested":    0.0,
		"total_profit_loss": 0.0,
		"fiat_balance":      wallet.FiatBalance,
		"locked":            wallet.Locked,
		"version":           wallet.Version,
	}

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionWallets,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, wallet)
}

func (r *repository) Update(wallet *models.Wallet) error {
	data := map[string]interface{}{
		"fiat_balance": wallet.FiatBalance,
		"locked":       wallet.Locked,
		"version":      wallet.Version,
	}

	resp, err := r.client.Databases.UpdateDocument(
		r.client.Config.DatabaseID,
		CollectionWallets,
		wallet.ID,
		r.client.Databases.WithUpdateDocumentData(data),
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, wallet)
}

func (r *repository) CreateTransaction(transaction *models.Transaction) error {
	data := map[string]interface{}{
		"user_id":      transaction.UserID,
		"type":         string(transaction.Type),
		"amount":       transaction.Amount,
		"status":       string(transaction.Status),
		"reference_id": transaction.ReferenceID,
	}

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionTransactions,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, transaction)
}

func (r *repository) UpdateTransaction(transaction *models.Transaction) error {
	data := map[string]interface{}{
		"status": string(transaction.Status),
	}

	resp, err := r.client.Databases.UpdateDocument(
		r.client.Config.DatabaseID,
		CollectionTransactions,
		transaction.ID,
		r.client.Databases.WithUpdateDocumentData(data),
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, transaction)
}

func (r *repository) GetUserTransaction(userID string) ([]models.Transaction, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionTransactions,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("user_id", userID),
			query.OrderDesc("$createdAt"),
			query.Limit(100),
		}),
	)
	if err != nil {
		return nil, err
	}

	transactions := make([]models.Transaction, len(resp.Documents))
	for i := range resp.Documents {
		if err := appwrite.DecodeListItem(resp, i, &transactions[i]); err != nil {
			return nil, fmt.Errorf("failed to decode transaction: %w", err)
		}
	}
	return transactions, nil
}

func (r *repository) WithLock(userID string, fn func(*models.Wallet) error) error {
	wallet, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	if wallet.Locked {
		return fmt.Errorf("wallet is locked")
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
