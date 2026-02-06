package trading

import (
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/appwrite"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/query"
)

const (
	CollectionVirtualWallets    = "virtual_wallets"
	CollectionUserPortfolios    = "user_portfolios"
	CollectionStockTransactions = "stock_transactions"
)

type Repository interface {
	// Virtual wallet operations
	CreateVirtualWallet(wallet *models.VirtualWallet) error
	GetVirtualWallet(userID string) (*models.VirtualWallet, error)
	UpdateVirtualWallet(wallet *models.VirtualWallet) error

	// Portfolio operations
	GetPortfolio(userID string) ([]models.UserPortfolio, error)
	GetPortfolioItem(userID, companyID string) (*models.UserPortfolio, error)
	CreatePortfolioItem(item *models.UserPortfolio) error
	UpdatePortfolioItem(item *models.UserPortfolio) error
	DeletePortfolioItem(userID, companyID string) error

	// Transaction operations
	CreateTransaction(tx *models.StockTransaction) error
	GetUserTransactions(userID string, limit, offset int) ([]models.StockTransaction, error)
	GetTransactionsByCompany(userID, companyID string, limit int) ([]models.StockTransaction, error)

	// Atomic trading operations
	ExecuteBuy(userID, companyID string, quantity int, pricePerShare float64) error
	ExecuteSell(userID, companyID string, quantity int, pricePerShare float64) error
}

type repository struct {
	client *appwrite.Client
}

func NewRepository(client *appwrite.Client) Repository {
	return &repository{client: client}
}

// Virtual wallet operations

func (r *repository) CreateVirtualWallet(wallet *models.VirtualWallet) error {
	data := map[string]interface{}{
		"user_id":           wallet.UserID,
		"balance":           wallet.Balance,
		"total_invested":    wallet.TotalInvested,
		"total_profit_loss": wallet.TotalProfitLoss,
	}

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionVirtualWallets,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, wallet)
}

func (r *repository) GetVirtualWallet(userID string) (*models.VirtualWallet, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionVirtualWallets,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("user_id", userID),
			query.Limit(1),
		}),
	)
	if err != nil {
		return nil, err
	}
	if len(resp.Documents) == 0 {
		return nil, fmt.Errorf("virtual wallet not found")
	}

	var wallet models.VirtualWallet
	if err := appwrite.DecodeListItem(resp, 0, &wallet); err != nil {
		return nil, fmt.Errorf("failed to decode wallet: %w", err)
	}
	return &wallet, nil
}

func (r *repository) UpdateVirtualWallet(wallet *models.VirtualWallet) error {
	data := map[string]interface{}{
		"balance":           wallet.Balance,
		"total_invested":    wallet.TotalInvested,
		"total_profit_loss": wallet.TotalProfitLoss,
	}

	resp, err := r.client.Databases.UpdateDocument(
		r.client.Config.DatabaseID,
		CollectionVirtualWallets,
		wallet.ID,
		appwrite.WithUpdateDocumentData(data),
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, wallet)
}

// Portfolio operations

func (r *repository) GetPortfolio(userID string) ([]models.UserPortfolio, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionUserPortfolios,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("user_id", userID),
			query.GreaterThan("quantity", 0),
			query.OrderDesc("total_invested"),
		}),
	)
	if err != nil {
		return nil, err
	}

	var portfolio []models.UserPortfolio
	for i := range resp.Documents {
		var item models.UserPortfolio
		if err := appwrite.DecodeListItem(resp, i, &item); err == nil {
			portfolio = append(portfolio, item)
		}
	}
	return portfolio, nil
}

func (r *repository) GetPortfolioItem(userID, companyID string) (*models.UserPortfolio, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionUserPortfolios,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("user_id", userID),
			query.Equal("company_id", companyID),
			query.Limit(1),
		}),
	)
	if err != nil {
		return nil, err
	}
	if len(resp.Documents) == 0 {
		return nil, fmt.Errorf("portfolio item not found")
	}

	var item models.UserPortfolio
	if err := appwrite.DecodeListItem(resp, 0, &item); err != nil {
		return nil, fmt.Errorf("failed to decode portfolio item: %w", err)
	}
	return &item, nil
}

func (r *repository) CreatePortfolioItem(item *models.UserPortfolio) error {
	data := map[string]interface{}{
		"user_id":        item.UserID,
		"company_id":     item.CompanyID,
		"quantity":       item.Quantity,
		"avg_buy_price":  item.AvgBuyPrice,
		"total_invested": item.TotalInvested,
	}

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionUserPortfolios,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, item)
}

func (r *repository) UpdatePortfolioItem(item *models.UserPortfolio) error {
	data := map[string]interface{}{
		"quantity":       item.Quantity,
		"avg_buy_price":  item.AvgBuyPrice,
		"total_invested": item.TotalInvested,
	}

	resp, err := r.client.Databases.UpdateDocument(
		r.client.Config.DatabaseID,
		CollectionUserPortfolios,
		item.ID,
		appwrite.WithUpdateDocumentData(data),
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, item)
}

func (r *repository) DeletePortfolioItem(userID, companyID string) error {
	item, err := r.GetPortfolioItem(userID, companyID)
	if err != nil {
		return err
	}
	_, err = r.client.Databases.DeleteDocument(
		r.client.Config.DatabaseID,
		CollectionUserPortfolios,
		item.ID,
	)
	return err
}

// Transaction operations

func (r *repository) CreateTransaction(tx *models.StockTransaction) error {
	data := map[string]interface{}{
		"user_id":         tx.UserID,
		"company_id":      tx.CompanyID,
		"type":            tx.Type,
		"quantity":        tx.Quantity,
		"price_per_share": tx.PricePerShare,
		"total_amount":    tx.TotalAmount,
		"status":          tx.Status,
		"reference_id":    tx.ReferenceID,
	}

	resp, err := r.client.Databases.CreateDocument(
		r.client.Config.DatabaseID,
		CollectionStockTransactions,
		id.Unique(),
		data,
	)
	if err != nil {
		return err
	}
	return appwrite.Decode(resp, tx)
}

func (r *repository) GetUserTransactions(userID string, limit, offset int) ([]models.StockTransaction, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionStockTransactions,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("user_id", userID),
			query.OrderDesc("$createdAt"),
			query.Limit(limit),
			query.Offset(offset),
		}),
	)
	if err != nil {
		return nil, err
	}

	var transactions []models.StockTransaction
	for i := range resp.Documents {
		var tx models.StockTransaction
		if err := appwrite.DecodeListItem(resp, i, &tx); err == nil {
			transactions = append(transactions, tx)
		}
	}
	return transactions, nil
}

func (r *repository) GetTransactionsByCompany(userID, companyID string, limit int) ([]models.StockTransaction, error) {
	resp, err := r.client.Databases.ListDocuments(
		r.client.Config.DatabaseID,
		CollectionStockTransactions,
		appwrite.WithListDocumentsQueries([]string{
			query.Equal("user_id", userID),
			query.Equal("company_id", companyID),
			query.OrderDesc("$createdAt"),
			query.Limit(limit),
		}),
	)
	if err != nil {
		return nil, err
	}

	var transactions []models.StockTransaction
	for i := range resp.Documents {
		var tx models.StockTransaction
		if err := appwrite.DecodeListItem(resp, i, &tx); err == nil {
			transactions = append(transactions, tx)
		}
	}
	return transactions, nil
}

// Atomic trading operations - Simulated with sequential calls

func (r *repository) ExecuteBuy(userID, companyID string, quantity int, pricePerShare float64) error {
	totalAmount := float64(quantity) * pricePerShare

	wallet, err := r.GetVirtualWallet(userID)
	if err != nil {
		return err
	}
	if wallet.Balance < totalAmount {
		return fmt.Errorf("insufficient balance")
	}

	wallet.Balance -= totalAmount
	wallet.TotalInvested += totalAmount
	if err := r.UpdateVirtualWallet(wallet); err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	item, err := r.GetPortfolioItem(userID, companyID)
	if err != nil {
		newItem := &models.UserPortfolio{
			UserID:        userID,
			CompanyID:     companyID,
			Quantity:      quantity,
			AvgBuyPrice:   pricePerShare,
			TotalInvested: totalAmount,
		}
		if err := r.CreatePortfolioItem(newItem); err != nil {
			return fmt.Errorf("failed to create portfolio item: %w", err)
		}
	} else {
		newQuantity := item.Quantity + quantity
		newTotalInvested := item.TotalInvested + totalAmount
		newAvgPrice := newTotalInvested / float64(newQuantity)

		item.Quantity = newQuantity
		item.TotalInvested = newTotalInvested
		item.AvgBuyPrice = newAvgPrice

		if err := r.UpdatePortfolioItem(item); err != nil {
			return fmt.Errorf("failed to update portfolio item: %w", err)
		}
	}

	tx := &models.StockTransaction{
		UserID:        userID,
		CompanyID:     companyID,
		Type:          models.StockTransactionBuy,
		Quantity:      quantity,
		PricePerShare: pricePerShare,
		TotalAmount:   totalAmount,
		Status:        models.StockTransactionCompleted,
	}
	return r.CreateTransaction(tx)
}

func (r *repository) ExecuteSell(userID, companyID string, quantity int, pricePerShare float64) error {
	totalAmount := float64(quantity) * pricePerShare

	item, err := r.GetPortfolioItem(userID, companyID)
	if err != nil {
		return fmt.Errorf("portfolio item not found: %w", err)
	}
	if item.Quantity < quantity {
		return fmt.Errorf("insufficient shares")
	}

	newQuantity := item.Quantity - quantity
	soldInvestment := (item.TotalInvested / float64(item.Quantity)) * float64(quantity)
	newTotalInvested := item.TotalInvested - soldInvestment

	item.Quantity = newQuantity
	item.TotalInvested = newTotalInvested

	if newQuantity == 0 {
		if err := r.DeletePortfolioItem(userID, companyID); err != nil {
			return fmt.Errorf("failed to delete portfolio item: %w", err)
		}
	} else {
		if err := r.UpdatePortfolioItem(item); err != nil {
			return fmt.Errorf("failed to update portfolio item: %w", err)
		}
	}

	wallet, err := r.GetVirtualWallet(userID)
	if err != nil {
		return err
	}

	profitLoss := totalAmount - soldInvestment
	wallet.Balance += totalAmount
	wallet.TotalInvested -= soldInvestment
	wallet.TotalProfitLoss += profitLoss

	if err := r.UpdateVirtualWallet(wallet); err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	tx := &models.StockTransaction{
		UserID:        userID,
		CompanyID:     companyID,
		Type:          models.StockTransactionSell,
		Quantity:      quantity,
		PricePerShare: pricePerShare,
		TotalAmount:   totalAmount,
		Status:        models.StockTransactionCompleted,
	}
	return r.CreateTransaction(tx)
}
