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

// =============================================================================
// CreateVirtualWallet — INSERT INTO virtual_wallets (user_id, balance,
//
//	total_invested, total_profit_loss)
//	VALUES ($1, $2, $3, $4) RETURNING *
//
// =============================================================================
func (r *repository) CreateVirtualWallet(wallet *models.VirtualWallet) error {
	query := `INSERT INTO virtual_wallets (user_id, balance, total_invested, total_profit_loss)
			  VALUES ($1, $2, $3, $4) RETURNING *`
	return r.client.ExecuteInsert(query, wallet,
		wallet.UserID, wallet.Balance.InexactFloat64(),
		wallet.TotalInvested.InexactFloat64(), wallet.TotalProfitLoss.InexactFloat64())
}

// =============================================================================
// GetVirtualWallet — SELECT * FROM virtual_wallets WHERE user_id = $1
// =============================================================================
func (r *repository) GetVirtualWallet(userID string) (*models.VirtualWallet, error) {
	var wallet models.VirtualWallet
	query := "SELECT * FROM virtual_wallets WHERE user_id = $1"
	err := r.client.ExecuteQueryRow(query, &wallet, userID)
	if err != nil {
		return nil, fmt.Errorf("virtual wallet not found for user %s", userID)
	}
	return &wallet, nil
}

// =============================================================================
// UpdateVirtualWallet — UPDATE virtual_wallets SET balance=$1, total_invested=$2,
//
//	total_profit_loss=$3 WHERE id = $4 RETURNING *
//
// =============================================================================
func (r *repository) UpdateVirtualWallet(wallet *models.VirtualWallet) error {
	query := `UPDATE virtual_wallets SET balance = $1, total_invested = $2, total_profit_loss = $3
			  WHERE id = $4 RETURNING *`
	return r.client.ExecuteUpdate(query, wallet,
		wallet.Balance.InexactFloat64(), wallet.TotalInvested.InexactFloat64(),
		wallet.TotalProfitLoss.InexactFloat64(), wallet.ID)
}

// =============================================================================
// GetPortfolio — SELECT * FROM user_portfolios WHERE user_id = $1 AND quantity > $2
// =============================================================================
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

// =============================================================================
// GetPortfolioItem — SELECT * FROM user_portfolios WHERE user_id = $1
//
//	AND company_id = $2
//
// =============================================================================
func (r *repository) GetPortfolioItem(userID, companyID string) (*models.UserPortfolio, error) {
	var item models.UserPortfolio
	query := "SELECT * FROM user_portfolios WHERE user_id = $1 AND company_id = $2"
	err := r.client.ExecuteQueryRow(query, &item, userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("portfolio item not found")
	}
	return &item, nil
}

// =============================================================================
// CreatePortfolioItem — INSERT INTO user_portfolios (user_id, company_id,
//
//	quantity, average_price, total_invested)
//	VALUES ($1, $2, $3, $4, $5) RETURNING *
//
// =============================================================================
func (r *repository) CreatePortfolioItem(item *models.UserPortfolio) error {
	query := `INSERT INTO user_portfolios (user_id, company_id, quantity, average_price, total_invested)
			  VALUES ($1, $2, $3, $4, $5) RETURNING *`
	return r.client.ExecuteInsert(query, item,
		item.UserID, item.CompanyID, item.Quantity,
		item.AvgBuyPrice.InexactFloat64(), item.TotalInvested.InexactFloat64())
}

// =============================================================================
// UpdatePortfolioItem — UPDATE user_portfolios SET quantity=$1, average_price=$2,
//
//	total_invested=$3 WHERE id = $4 RETURNING *
//
// =============================================================================
func (r *repository) UpdatePortfolioItem(item *models.UserPortfolio) error {
	query := `UPDATE user_portfolios SET quantity = $1, average_price = $2, total_invested = $3
			  WHERE id = $4 RETURNING *`
	return r.client.ExecuteUpdate(query, item,
		item.Quantity, item.AvgBuyPrice.InexactFloat64(), item.TotalInvested.InexactFloat64(), item.ID)
}

// =============================================================================
// DeletePortfolioItem — DELETE FROM user_portfolios WHERE user_id = $1
//
//	AND company_id = $2
//
// =============================================================================
func (r *repository) DeletePortfolioItem(userID, companyID string) error {
	query := "DELETE FROM user_portfolios WHERE user_id = $1 AND company_id = $2"
	return r.client.ExecuteDelete(query, userID, companyID)
}

// =============================================================================
// CreateTransaction — INSERT INTO stock_transactions (user_id, company_id, type,
//
//	quantity, price_per_share, total_amount, status, reference_id)
//	VALUES (...) RETURNING *
//
// =============================================================================
func (r *repository) CreateTransaction(tx *models.StockTransaction) error {
	query := `INSERT INTO stock_transactions (user_id, company_id, type, quantity, price_per_share, total_amount, status, reference_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	return r.client.ExecuteInsert(query, tx,
		tx.UserID, tx.CompanyID, tx.Type, tx.Quantity,
		tx.PricePerShare.InexactFloat64(), tx.TotalAmount.InexactFloat64(),
		tx.Status, tx.ReferenceID)
}

// =============================================================================
// GetUserTransactions — SELECT * FROM stock_transactions WHERE user_id = $1
//
//	ORDER BY created_at DESC LIMIT $2 OFFSET $3
//
// =============================================================================
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

// =============================================================================
// GetTransactionsByCompany — SELECT * FROM stock_transactions
//
//	WHERE user_id = $1 AND company_id = $2
//	ORDER BY created_at DESC LIMIT $3
//
// =============================================================================
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

// =============================================================================
// ExecuteBuy — Business logic for buying stocks
// Deducts from wallet, updates/creates portfolio, records transaction
// =============================================================================
func (r *repository) ExecuteBuy(userID, companyID string, quantity int, pricePerShare decimal.Decimal) error {
	totalAmount := pricePerShare.Mul(decimal.NewFromInt(int64(quantity)))

	// 1. SELECT * FROM virtual_wallets WHERE user_id = $1 — Check wallet balance
	wallet, err := r.GetVirtualWallet(userID)
	if err != nil {
		return err
	}

	// If balance is 0 but fiat_balance has funds (user topped up), use fiat_balance as balance
	if wallet.Balance.IsZero() && wallet.TotalInvested.IsZero() {
		// Try to read raw fiat_balance from database
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

	// 2. UPDATE virtual_wallets SET balance = balance - $1, total_invested = total_invested + $2
	wallet.Balance = wallet.Balance.Sub(totalAmount)
	wallet.TotalInvested = wallet.TotalInvested.Add(totalAmount)
	if err := r.UpdateVirtualWallet(wallet); err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	// 3. SELECT * FROM user_portfolios WHERE user_id = $1 AND company_id = $2
	item, err := r.GetPortfolioItem(userID, companyID)
	if err != nil {
		// New position — INSERT INTO user_portfolios
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
		// Existing position — UPDATE user_portfolios SET ... (weighted average)
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

	// 4. INSERT INTO stock_transactions — Record the transaction
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

// =============================================================================
// ExecuteSell — Business logic for selling stocks
// Checks portfolio, updates wallet with profit/loss, records transaction
// =============================================================================
func (r *repository) ExecuteSell(userID, companyID string, quantity int, pricePerShare decimal.Decimal) error {
	totalAmount := pricePerShare.Mul(decimal.NewFromInt(int64(quantity)))

	// 1. SELECT * FROM user_portfolios WHERE user_id = $1 AND company_id = $2
	item, err := r.GetPortfolioItem(userID, companyID)
	if err != nil {
		return fmt.Errorf("portfolio item not found: %w", err)
	}
	if item.Quantity < quantity {
		return fmt.Errorf("insufficient shares")
	}

	newQuantity := item.Quantity - quantity

	if item.Quantity == 0 {
		return fmt.Errorf("invalid portfolio state: quantity is 0")
	}

	// 2. Calculate proportional investment being sold
	soldInvestment := item.TotalInvested.Div(decimal.NewFromInt(int64(item.Quantity))).Mul(decimal.NewFromInt(int64(quantity)))
	newTotalInvested := item.TotalInvested.Sub(soldInvestment)

	item.Quantity = newQuantity
	item.TotalInvested = newTotalInvested

	// 3. DELETE or UPDATE user_portfolios
	if newQuantity == 0 {
		// DELETE FROM user_portfolios WHERE user_id = $1 AND company_id = $2
		if err := r.DeletePortfolioItem(userID, companyID); err != nil {
			return fmt.Errorf("failed to delete portfolio item: %w", err)
		}
	} else {
		// UPDATE user_portfolios SET quantity=$1, total_invested=$2 WHERE id = $3
		if err := r.UpdatePortfolioItem(item); err != nil {
			return fmt.Errorf("failed to update portfolio item: %w", err)
		}
	}

	// 4. UPDATE virtual_wallets — Add sell proceeds and calculate profit/loss
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

	// 5. INSERT INTO stock_transactions — Record the transaction
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
