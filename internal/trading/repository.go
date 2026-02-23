package trading

import (
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

const (
	TableVirtualWallets    = "virtual_wallets"
	TableUserPortfolios    = "user_portfolios"
	TableStockTransactions = "stock_transactions"
)

type Repository interface {
	CreateVirtualWallet(wallet *models.VirtualWallet) error
	GetVirtualWallet(userID string) (*models.VirtualWallet, error)
	UpdateVirtualWallet(wallet *models.VirtualWallet) error

	GetPortfolio(userID string) ([]models.UserPortfolio, error)
	GetPortfolioItem(userID, companyID string) (*models.UserPortfolio, error)
	CreatePortfolioItem(item *models.UserPortfolio) error
	UpdatePortfolioItem(item *models.UserPortfolio) error
	DeletePortfolioItem(userID, companyID string) error

	CreateTransaction(tx *models.StockTransaction) error
	GetUserTransactions(userID string, limit, offset int) ([]models.StockTransaction, error)
	GetTransactionsByCompany(userID, companyID string, limit int) ([]models.StockTransaction, error)

	ExecuteBuy(userID, companyID string, quantity int, pricePerShare decimal.Decimal) error
	ExecuteSell(userID, companyID string, quantity int, pricePerShare decimal.Decimal) error
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

func (r *repository) CreateVirtualWallet(wallet *models.VirtualWallet) error {
	query := `INSERT INTO virtual_wallets (user_id, balance, total_invested, total_profit_loss)
			  VALUES ($1, $2, $3, $4) RETURNING *`
	return r.client.ExecuteInsert(query, wallet,
		wallet.UserID, wallet.Balance.InexactFloat64(),
		wallet.TotalInvested.InexactFloat64(), wallet.TotalProfitLoss.InexactFloat64())
}

func (r *repository) GetVirtualWallet(userID string) (*models.VirtualWallet, error) {
	var wallet models.VirtualWallet
	query := "SELECT * FROM virtual_wallets WHERE user_id = $1"
	err := r.client.ExecuteQueryRow(query, &wallet, userID)
	if err != nil {
		return nil, fmt.Errorf("virtual wallet not found for user %s", userID)
	}
	return &wallet, nil
}

func (r *repository) UpdateVirtualWallet(wallet *models.VirtualWallet) error {
	query := `UPDATE virtual_wallets SET balance = $1, total_invested = $2, total_profit_loss = $3
			  WHERE id = $4 RETURNING *`
	return r.client.ExecuteUpdate(query, wallet,
		wallet.Balance.InexactFloat64(), wallet.TotalInvested.InexactFloat64(),
		wallet.TotalProfitLoss.InexactFloat64(), wallet.ID)
}

func (r *repository) GetPortfolio(userID string) ([]models.UserPortfolio, error) {
	var portfolio []models.UserPortfolio
	query := "SELECT * FROM user_portfolios WHERE user_id = $1 AND quantity > $2"
	err := r.client.ExecuteQuery(query, &portfolio, userID, 0)
	if err != nil {
		return nil, err
	}
	if portfolio == nil {
		portfolio = []models.UserPortfolio{}
	}
	return portfolio, nil
}

func (r *repository) GetPortfolioItem(userID, companyID string) (*models.UserPortfolio, error) {
	var item models.UserPortfolio
	query := "SELECT * FROM user_portfolios WHERE user_id = $1 AND company_id = $2"
	err := r.client.ExecuteQueryRow(query, &item, userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("portfolio item not found")
	}
	return &item, nil
}

func (r *repository) CreatePortfolioItem(item *models.UserPortfolio) error {
	query := `INSERT INTO user_portfolios (user_id, company_id, quantity, average_price, total_invested)
			  VALUES ($1, $2, $3, $4, $5) RETURNING *`
	return r.client.ExecuteInsert(query, item,
		item.UserID, item.CompanyID, item.Quantity,
		item.AvgBuyPrice.InexactFloat64(), item.TotalInvested.InexactFloat64())
}

func (r *repository) UpdatePortfolioItem(item *models.UserPortfolio) error {
	query := `UPDATE user_portfolios SET quantity = $1, average_price = $2, total_invested = $3
			  WHERE id = $4 RETURNING *`
	return r.client.ExecuteUpdate(query, item,
		item.Quantity, item.AvgBuyPrice.InexactFloat64(), item.TotalInvested.InexactFloat64(), item.ID)
}

func (r *repository) DeletePortfolioItem(userID, companyID string) error {
	query := "DELETE FROM user_portfolios WHERE user_id = $1 AND company_id = $2"
	return r.client.ExecuteDelete(query, userID, companyID)
}

func (r *repository) CreateTransaction(tx *models.StockTransaction) error {
	query := `INSERT INTO stock_transactions (user_id, company_id, type, quantity, price_per_share, total_amount, status, reference_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	return r.client.ExecuteInsert(query, tx,
		tx.UserID, tx.CompanyID, tx.Type, tx.Quantity,
		tx.PricePerShare.InexactFloat64(), tx.TotalAmount.InexactFloat64(),
		tx.Status, tx.ReferenceID)
}

func (r *repository) GetUserTransactions(userID string, limit, offset int) ([]models.StockTransaction, error) {
	var transactions []models.StockTransaction
	query := "SELECT * FROM stock_transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	err := r.client.ExecuteQuery(query, &transactions, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	if transactions == nil {
		transactions = []models.StockTransaction{}
	}
	return transactions, nil
}

func (r *repository) GetTransactionsByCompany(userID, companyID string, limit int) ([]models.StockTransaction, error) {
	var transactions []models.StockTransaction
	query := "SELECT * FROM stock_transactions WHERE user_id = $1 AND company_id = $2 ORDER BY created_at DESC LIMIT $3"
	err := r.client.ExecuteQuery(query, &transactions, userID, companyID, limit)
	if err != nil {
		return nil, err
	}
	if transactions == nil {
		transactions = []models.StockTransaction{}
	}
	return transactions, nil
}

func (r *repository) ExecuteBuy(userID, companyID string, quantity int, pricePerShare decimal.Decimal) error {
	totalAmount := pricePerShare.Mul(decimal.NewFromInt(int64(quantity)))

	wallet, err := r.GetVirtualWallet(userID)
	if err != nil {
		return err
	}

	if wallet.Balance.IsZero() && wallet.TotalInvested.IsZero() {
		type WalletBalance struct {
			FiatBalance decimal.Decimal `json:"fiat_balance"`
		}
		var wb WalletBalance
		rawQuery := "SELECT fiat_balance FROM virtual_wallets WHERE user_id = $1"
		if rawErr := r.client.ExecuteQueryRow(rawQuery, &wb, userID); rawErr == nil {
			wallet.Balance = wb.FiatBalance
		}
	}

	if wallet.Balance.LessThan(totalAmount) {
		return fmt.Errorf("insufficient balance")
	}

	wallet.Balance = wallet.Balance.Sub(totalAmount)
	wallet.TotalInvested = wallet.TotalInvested.Add(totalAmount)
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
		newTotalInvested := item.TotalInvested.Add(totalAmount)
		newAvgPrice := newTotalInvested.Div(decimal.NewFromInt(int64(newQuantity)))

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

func (r *repository) ExecuteSell(userID, companyID string, quantity int, pricePerShare decimal.Decimal) error {
	totalAmount := pricePerShare.Mul(decimal.NewFromInt(int64(quantity)))

	item, err := r.GetPortfolioItem(userID, companyID)
	if err != nil {
		return fmt.Errorf("portfolio item not found: %w", err)
	}
	if item.Quantity < quantity {
		return fmt.Errorf("insufficient shares")
	}

	newQuantity := item.Quantity - quantity

	// Calculate cost basis of sold shares (proportional to quantity sold)
	soldInvestment := item.TotalInvested.Div(decimal.NewFromInt(int64(item.Quantity))).Mul(decimal.NewFromInt(int64(quantity)))
	newTotalInvested := item.TotalInvested.Sub(soldInvestment)

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

	profitLoss := totalAmount.Sub(soldInvestment)

	wallet.Balance = wallet.Balance.Add(totalAmount)
	wallet.TotalInvested = wallet.TotalInvested.Sub(soldInvestment)
	wallet.TotalProfitLoss = wallet.TotalProfitLoss.Add(profitLoss)

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
