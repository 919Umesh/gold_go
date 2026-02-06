package trading

import (
	"errors"
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/models"
)

type Service interface {
	// Wallet
	CreateWallet(userID string) (*models.VirtualWallet, error)
	GetWallet(userID string) (*models.VirtualWallet, error)
	GetOrCreateWallet(userID string) (*models.VirtualWallet, error)

	// Portfolio
	GetPortfolio(userID string) (*PortfolioSummary, error)

	// Trading
	BuyStock(userID, symbol string, quantity int) (*TradeResult, error)
	SellStock(userID, symbol string, quantity int) (*TradeResult, error)
	GetTransactionHistory(userID string, limit, offset int) ([]models.StockTransaction, error)
}

type PortfolioSummary struct {
	TotalValue      float64                `json:"total_value"`
	TotalInvested   float64                `json:"total_invested"`
	TotalProfitLoss float64                `json:"total_profit_loss"`
	ProfitLossPct   float64                `json:"profit_loss_pct"`
	Items           []PortfolioItemSummary `json:"items"`
}

type PortfolioItemSummary struct {
	models.UserPortfolio
	CurrentPrice  float64 `json:"current_price"`
	CurrentValue  float64 `json:"current_value"`
	ProfitLoss    float64 `json:"profit_loss"`
	ProfitLossPct float64 `json:"profit_loss_pct"`
	CompanyName   string  `json:"company_name"`
	CompanySymbol string  `json:"company_symbol"`
	CompanySector string  `json:"company_sector"`
}

type TradeResult struct {
	Success       bool    `json:"success"`
	Message       string  `json:"message"`
	TransactionID string  `json:"transaction_id,omitempty"`
	Quantity      int     `json:"quantity"`
	PricePerShare float64 `json:"price_per_share"`
	TotalAmount   float64 `json:"total_amount"`
	NewBalance    float64 `json:"new_balance"`
}

type service struct {
	repo      Repository
	stockRepo stock.Repository
}

func NewService(repo Repository, stockRepo stock.Repository) Service {
	return &service{
		repo:      repo,
		stockRepo: stockRepo,
	}
}

func (s *service) CreateWallet(userID string) (*models.VirtualWallet, error) {
	// Check if wallet exists
	_, err := s.repo.GetVirtualWallet(userID)
	if err == nil {
		return nil, errors.New("wallet already exists")
	}

	wallet := &models.VirtualWallet{
		UserID:          userID,
		Balance:         1000000.00, // Initial balance
		TotalInvested:   0,
		TotalProfitLoss: 0,
	}

	if err := s.repo.CreateVirtualWallet(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *service) GetWallet(userID string) (*models.VirtualWallet, error) {
	return s.repo.GetVirtualWallet(userID)
}

func (s *service) GetOrCreateWallet(userID string) (*models.VirtualWallet, error) {
	wallet, err := s.repo.GetVirtualWallet(userID)
	if err == nil {
		return wallet, nil
	}
	return s.CreateWallet(userID)
}

func (s *service) GetPortfolio(userID string) (*PortfolioSummary, error) {
	portfolio, err := s.repo.GetPortfolio(userID)
	if err != nil {
		return nil, err
	}

	summary := &PortfolioSummary{
		Items: make([]PortfolioItemSummary, 0, len(portfolio)),
	}

	for _, item := range portfolio {
		// Fetch current price
		price, err := s.stockRepo.GetLatestPrice(item.CompanyID)
		currentPrice := 0.0
		if err == nil {
			currentPrice = price.ClosePrice
		}

		// Fetch company details for name/symbol
		company, err := s.stockRepo.GetCompanyByID(item.CompanyID)
		companyName := ""
		companySymbol := ""
		companySector := ""
		if err == nil {
			companyName = company.Name
			companySymbol = company.Symbol
			companySector = company.Sector
		}

		currentValue := float64(item.Quantity) * currentPrice
		profitLoss := currentValue - item.TotalInvested
		profitLossPct := 0.0
		if item.TotalInvested > 0 {
			profitLossPct = (profitLoss / item.TotalInvested) * 100
		}

		summary.TotalValue += currentValue
		summary.TotalInvested += item.TotalInvested

		summary.Items = append(summary.Items, PortfolioItemSummary{
			UserPortfolio: item,
			CurrentPrice:  currentPrice,
			CurrentValue:  currentValue,
			ProfitLoss:    profitLoss,
			ProfitLossPct: profitLossPct,
			CompanyName:   companyName,
			CompanySymbol: companySymbol,
			CompanySector: companySector,
		})
	}

	summary.TotalProfitLoss = summary.TotalValue - summary.TotalInvested
	if summary.TotalInvested > 0 {
		summary.ProfitLossPct = (summary.TotalProfitLoss / summary.TotalInvested) * 100
	}

	return summary, nil
}

func (s *service) BuyStock(userID, symbol string, quantity int) (*TradeResult, error) {
	if quantity <= 0 {
		return &TradeResult{Success: false, Message: "Quantity must be positive"}, nil
	}

	// Resolve Symbol to CompanyID
	company, err := s.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return &TradeResult{Success: false, Message: "Company not found"}, nil
	}

	// 1. Get latest price
	price, err := s.stockRepo.GetLatestPrice(company.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock price: %w", err)
	}

	totalAmount := float64(quantity) * price.ClosePrice

	// 2. Execute buy
	err = s.repo.ExecuteBuy(userID, company.ID, quantity, price.ClosePrice)
	if err != nil {
		return &TradeResult{Success: false, Message: err.Error()}, nil
	}

	// 3. Get updated balance
	wallet, _ := s.repo.GetVirtualWallet(userID)

	return &TradeResult{
		Success:       true,
		Message:       "Stock purchased successfully",
		Quantity:      quantity,
		PricePerShare: price.ClosePrice,
		TotalAmount:   totalAmount,
		NewBalance:    wallet.Balance,
		// TransactionID: ??? // ExecuteBuy creates a transaction but doesn't return ID.
		// For now, omitting ID or we could fetch the last transaction.
	}, nil
}

func (s *service) SellStock(userID, symbol string, quantity int) (*TradeResult, error) {
	if quantity <= 0 {
		return &TradeResult{Success: false, Message: "Quantity must be positive"}, nil
	}

	// Resolve Symbol to CompanyID
	company, err := s.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return &TradeResult{Success: false, Message: "Company not found"}, nil
	}

	// 1. Get latest price
	price, err := s.stockRepo.GetLatestPrice(company.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock price: %w", err)
	}

	totalAmount := float64(quantity) * price.ClosePrice

	// 2. Execute sell
	err = s.repo.ExecuteSell(userID, company.ID, quantity, price.ClosePrice)
	if err != nil {
		return &TradeResult{Success: false, Message: err.Error()}, nil
	}

	// 3. Get updated balance
	wallet, _ := s.repo.GetVirtualWallet(userID)

	return &TradeResult{
		Success:       true,
		Message:       "Stock sold successfully",
		Quantity:      quantity,
		PricePerShare: price.ClosePrice,
		TotalAmount:   totalAmount,
		NewBalance:    wallet.Balance,
	}, nil
}

func (s *service) GetTransactionHistory(userID string, limit, offset int) ([]models.StockTransaction, error) {
	if limit <= 0 {
		limit = 20
	}
	// Renaming in service implies calling repository method which is GetUserTransactions
	return s.repo.GetUserTransactions(userID, limit, offset)
}
