package trading

import (
	"time"

	"github.com/919Umesh/gold_go/models"
	"gorm.io/gorm"
)

type Repository interface {
	// Virtual wallet operations
	CreateVirtualWallet(wallet *models.VirtualWallet) error
	GetVirtualWallet(userID uint) (*models.VirtualWallet, error)
	UpdateVirtualWallet(wallet *models.VirtualWallet) error

	// Portfolio operations
	GetPortfolio(userID uint) ([]models.UserPortfolio, error)
	GetPortfolioItem(userID, companyID uint) (*models.UserPortfolio, error)
	CreatePortfolioItem(item *models.UserPortfolio) error
	UpdatePortfolioItem(item *models.UserPortfolio) error
	DeletePortfolioItem(userID, companyID uint) error

	// Transaction operations
	CreateTransaction(tx *models.StockTransaction) error
	GetUserTransactions(userID uint, limit, offset int) ([]models.StockTransaction, error)
	GetTransactionsByCompany(userID, companyID uint, limit int) ([]models.StockTransaction, error)

	// Atomic trading operations
	ExecuteBuy(userID, companyID uint, quantity int, pricePerShare float64) error
	ExecuteSell(userID, companyID uint, quantity int, pricePerShare float64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Virtual wallet operations

func (r *repository) CreateVirtualWallet(wallet *models.VirtualWallet) error {
	query := `
		INSERT INTO virtual_wallets (user_id, balance, total_invested, total_profit_loss, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	return r.db.Raw(query, wallet.UserID, wallet.Balance, wallet.TotalInvested,
		wallet.TotalProfitLoss, time.Now(), time.Now()).Scan(&wallet.ID).Error
}

func (r *repository) GetVirtualWallet(userID uint) (*models.VirtualWallet, error) {
	var wallet models.VirtualWallet
	query := `SELECT * FROM virtual_wallets WHERE user_id = ? LIMIT 1`
	err := r.db.Raw(query, userID).Scan(&wallet).Error
	if err != nil {
		return nil, err
	}
	if wallet.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &wallet, nil
}

func (r *repository) UpdateVirtualWallet(wallet *models.VirtualWallet) error {
	query := `
		UPDATE virtual_wallets 
		SET balance = ?, total_invested = ?, total_profit_loss = ?, updated_at = ?
		WHERE id = ?
	`
	return r.db.Exec(query, wallet.Balance, wallet.TotalInvested, wallet.TotalProfitLoss,
		time.Now(), wallet.ID).Error
}

// Portfolio operations

func (r *repository) GetPortfolio(userID uint) ([]models.UserPortfolio, error) {
	var portfolio []models.UserPortfolio
	query := `
		SELECT up.*, c.symbol, c.name, c.sector
		FROM user_portfolios up
		INNER JOIN companies c ON up.company_id = c.id
		WHERE up.user_id = ? AND up.quantity > 0
		ORDER BY up.total_invested DESC
	`
	err := r.db.Raw(query, userID).Scan(&portfolio).Error
	return portfolio, err
}

func (r *repository) GetPortfolioItem(userID, companyID uint) (*models.UserPortfolio, error) {
	var item models.UserPortfolio
	query := `SELECT * FROM user_portfolios WHERE user_id = ? AND company_id = ? LIMIT 1`
	err := r.db.Raw(query, userID, companyID).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, nil
}

func (r *repository) CreatePortfolioItem(item *models.UserPortfolio) error {
	query := `
		INSERT INTO user_portfolios (user_id, company_id, quantity, avg_buy_price, total_invested, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	return r.db.Raw(query, item.UserID, item.CompanyID, item.Quantity, item.AvgBuyPrice,
		item.TotalInvested, time.Now(), time.Now()).Scan(&item.ID).Error
}

func (r *repository) UpdatePortfolioItem(item *models.UserPortfolio) error {
	query := `
		UPDATE user_portfolios 
		SET quantity = ?, avg_buy_price = ?, total_invested = ?, updated_at = ?
		WHERE id = ?
	`
	return r.db.Exec(query, item.Quantity, item.AvgBuyPrice, item.TotalInvested,
		time.Now(), item.ID).Error
}

func (r *repository) DeletePortfolioItem(userID, companyID uint) error {
	query := `DELETE FROM user_portfolios WHERE user_id = ? AND company_id = ?`
	return r.db.Exec(query, userID, companyID).Error
}

// Transaction operations

func (r *repository) CreateTransaction(tx *models.StockTransaction) error {
	query := `
		INSERT INTO stock_transactions (user_id, company_id, type, quantity, price_per_share, total_amount, status, reference_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	return r.db.Raw(query, tx.UserID, tx.CompanyID, tx.Type, tx.Quantity, tx.PricePerShare,
		tx.TotalAmount, tx.Status, tx.ReferenceID, time.Now()).Scan(&tx.ID).Error
}

func (r *repository) GetUserTransactions(userID uint, limit, offset int) ([]models.StockTransaction, error) {
	var transactions []models.StockTransaction
	query := `
		SELECT st.*, c.symbol, c.name
		FROM stock_transactions st
		INNER JOIN companies c ON st.company_id = c.id
		WHERE st.user_id = ?
		ORDER BY st.created_at DESC
		LIMIT ? OFFSET ?
	`
	err := r.db.Raw(query, userID, limit, offset).Scan(&transactions).Error
	return transactions, err
}

func (r *repository) GetTransactionsByCompany(userID, companyID uint, limit int) ([]models.StockTransaction, error) {
	var transactions []models.StockTransaction
	query := `
		SELECT * FROM stock_transactions 
		WHERE user_id = ? AND company_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`
	err := r.db.Raw(query, userID, companyID, limit).Scan(&transactions).Error
	return transactions, err
}

// Atomic trading operations

func (r *repository) ExecuteBuy(userID, companyID uint, quantity int, pricePerShare float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		totalAmount := float64(quantity) * pricePerShare

		// 1. Lock and check virtual wallet balance
		var wallet models.VirtualWallet
		walletQuery := `SELECT * FROM virtual_wallets WHERE user_id = ? FOR UPDATE`
		if err := tx.Raw(walletQuery, userID).Scan(&wallet).Error; err != nil {
			return err
		}

		if wallet.ID == 0 {
			return gorm.ErrRecordNotFound
		}

		if wallet.Balance < totalAmount {
			return gorm.ErrInvalidData // Insufficient balance
		}

		// 2. Deduct from wallet
		updateWalletQuery := `
			UPDATE virtual_wallets 
			SET balance = balance - ?, total_invested = total_invested + ?, updated_at = ?
			WHERE id = ?
		`
		if err := tx.Exec(updateWalletQuery, totalAmount, totalAmount, time.Now(), wallet.ID).Error; err != nil {
			return err
		}

		// 3. Update or create portfolio item
		var portfolio models.UserPortfolio
		portfolioQuery := `SELECT * FROM user_portfolios WHERE user_id = ? AND company_id = ? FOR UPDATE`
		err := tx.Raw(portfolioQuery, userID, companyID).Scan(&portfolio).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if portfolio.ID == 0 {
			// Create new portfolio item
			createPortfolioQuery := `
				INSERT INTO user_portfolios (user_id, company_id, quantity, avg_buy_price, total_invested, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`
			if err := tx.Exec(createPortfolioQuery, userID, companyID, quantity, pricePerShare,
				totalAmount, time.Now(), time.Now()).Error; err != nil {
				return err
			}
		} else {
			// Update existing portfolio item
			newQuantity := portfolio.Quantity + quantity
			newTotalInvested := portfolio.TotalInvested + totalAmount
			newAvgPrice := newTotalInvested / float64(newQuantity)

			updatePortfolioQuery := `
				UPDATE user_portfolios 
				SET quantity = ?, avg_buy_price = ?, total_invested = ?, updated_at = ?
				WHERE id = ?
			`
			if err := tx.Exec(updatePortfolioQuery, newQuantity, newAvgPrice, newTotalInvested,
				time.Now(), portfolio.ID).Error; err != nil {
				return err
			}
		}

		// 4. Record transaction
		createTxQuery := `
			INSERT INTO stock_transactions (user_id, company_id, type, quantity, price_per_share, total_amount, status, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
		return tx.Exec(createTxQuery, userID, companyID, models.StockTransactionBuy, quantity,
			pricePerShare, totalAmount, models.StockTransactionCompleted, time.Now()).Error
	})
}

func (r *repository) ExecuteSell(userID, companyID uint, quantity int, pricePerShare float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		totalAmount := float64(quantity) * pricePerShare

		// 1. Lock and check portfolio
		var portfolio models.UserPortfolio
		portfolioQuery := `SELECT * FROM user_portfolios WHERE user_id = ? AND company_id = ? FOR UPDATE`
		if err := tx.Raw(portfolioQuery, userID, companyID).Scan(&portfolio).Error; err != nil {
			return err
		}

		if portfolio.ID == 0 {
			return gorm.ErrRecordNotFound
		}

		if portfolio.Quantity < quantity {
			return gorm.ErrInvalidData // Insufficient shares
		}

		// 2. Update portfolio
		newQuantity := portfolio.Quantity - quantity
		soldInvestment := (portfolio.TotalInvested / float64(portfolio.Quantity)) * float64(quantity)
		newTotalInvested := portfolio.TotalInvested - soldInvestment

		if newQuantity == 0 {
			// Delete portfolio item if all shares sold
			deleteQuery := `DELETE FROM user_portfolios WHERE id = ?`
			if err := tx.Exec(deleteQuery, portfolio.ID).Error; err != nil {
				return err
			}
		} else {
			// Update portfolio item
			updatePortfolioQuery := `
				UPDATE user_portfolios 
				SET quantity = ?, total_invested = ?, updated_at = ?
				WHERE id = ?
			`
			if err := tx.Exec(updatePortfolioQuery, newQuantity, newTotalInvested,
				time.Now(), portfolio.ID).Error; err != nil {
				return err
			}
		}

		// 3. Add to wallet
		profitLoss := totalAmount - soldInvestment
		updateWalletQuery := `
			UPDATE virtual_wallets 
			SET balance = balance + ?, total_invested = total_invested - ?, total_profit_loss = total_profit_loss + ?, updated_at = ?
			WHERE user_id = ?
		`
		if err := tx.Exec(updateWalletQuery, totalAmount, soldInvestment, profitLoss,
			time.Now(), userID).Error; err != nil {
			return err
		}

		// 4. Record transaction
		createTxQuery := `
			INSERT INTO stock_transactions (user_id, company_id, type, quantity, price_per_share, total_amount, status, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
		return tx.Exec(createTxQuery, userID, companyID, models.StockTransactionSell, quantity,
			pricePerShare, totalAmount, models.StockTransactionCompleted, time.Now()).Error
	})
}
