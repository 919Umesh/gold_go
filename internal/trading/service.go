package trading

import (
	"fmt"

	"github.com/919Umesh/gold_go/internal/stock"
	"github.com/919Umesh/gold_go/models"
)

type Service interface {
	// Wallet operations
	GetOrCreateWallet(userID uint) (*models.VirtualWallet, error)
	GetWalletBalance(userID uint) (float64, error)

	// Portfolio operations
	GetPortfolio(userID uint) (*PortfolioSummary, error)
	GetPortfolioValue(userID uint) (float64, error)

	// Trading operations
	BuyStock(userID uint, symbol string, quantity int) (*TradeResult, error)
	SellStock(userID uint, symbol string, quantity int) (*TradeResult, error)

	// Transaction history
	GetTransactionHistory(userID uint, limit, offset int) ([]models.StockTransaction, error)
}

type PortfolioSummary struct {
	TotalValue        float64            `json:"total_value"`
	TotalInvested     float64            `json:"total_invested"`
	TotalProfitLoss   float64            `json:"total_profit_loss"`
	ProfitLossPercent float64            `json:"profit_loss_percent"`
	Holdings          []PortfolioHolding `json:"holdings"`
}

type PortfolioHolding struct {
	models.UserPortfolio
	CompanySymbol     string  `json:"company_symbol"`
	CompanyName       string  `json:"company_name"`
	CurrentPrice      float64 `json:"current_price"`
	CurrentValue      float64 `json:"current_value"`
	ProfitLoss        float64 `json:"profit_loss"`
	ProfitLossPercent float64 `json:"profit_loss_percent"`
}

type TradeResult struct {
	Success       bool    `json:"success"`
	Message       string  `json:"message"`
	TransactionID uint    `json:"transaction_id,omitempty"`
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

func (s *service) GetOrCreateWallet(userID uint) (*models.VirtualWallet, error) {
	wallet, err := s.repo.GetVirtualWallet(userID)
	if err == nil {
		return wallet, nil
	}

	// Create new wallet with initial balance
	newWallet := &models.VirtualWallet{
		UserID:          userID,
		Balance:         1000000.00, // NPR 10 lakh
		TotalInvested:   0,
		TotalProfitLoss: 0,
	}

	if err := s.repo.CreateVirtualWallet(newWallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return newWallet, nil
}

func (s *service) GetWalletBalance(userID uint) (float64, error) {
	wallet, err := s.GetOrCreateWallet(userID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (s *service) GetPortfolio(userID uint) (*PortfolioSummary, error) {
	portfolio, err := s.repo.GetPortfolio(userID)
	if err != nil {
		return nil, err
	}

	holdings := make([]PortfolioHolding, 0, len(portfolio))
	totalValue := 0.0
	totalInvested := 0.0

	for _, item := range portfolio {
		// Get current price
		currentPrice, err := s.stockRepo.GetLatestPrice(item.CompanyID)
		if err != nil {
			continue
		}

		currentValue := item.CalculateCurrentValue(currentPrice.ClosePrice)
		profitLoss := item.CalculateProfitLoss(currentPrice.ClosePrice)
		profitLossPercent := item.CalculateProfitLossPercentage(currentPrice.ClosePrice)

		// Get company info
		company, _ := s.stockRepo.GetCompanyByID(item.CompanyID)

		holdings = append(holdings, PortfolioHolding{
			UserPortfolio:     item,
			CompanySymbol:     company.Symbol,
			CompanyName:       company.Name,
			CurrentPrice:      currentPrice.ClosePrice,
			CurrentValue:      currentValue,
			ProfitLoss:        profitLoss,
			ProfitLossPercent: profitLossPercent,
		})

		totalValue += currentValue
		totalInvested += item.TotalInvested
	}

	totalProfitLoss := totalValue - totalInvested
	profitLossPercent := 0.0
	if totalInvested > 0 {
		profitLossPercent = (totalProfitLoss / totalInvested) * 100
	}

	return &PortfolioSummary{
		TotalValue:        totalValue,
		TotalInvested:     totalInvested,
		TotalProfitLoss:   totalProfitLoss,
		ProfitLossPercent: profitLossPercent,
		Holdings:          holdings,
	}, nil
}

func (s *service) GetPortfolioValue(userID uint) (float64, error) {
	summary, err := s.GetPortfolio(userID)
	if err != nil {
		return 0, err
	}
	return summary.TotalValue, nil
}

func (s *service) BuyStock(userID uint, symbol string, quantity int) (*TradeResult, error) {
	if quantity <= 0 {
		return &TradeResult{
			Success: false,
			Message: "Quantity must be greater than 0",
		}, nil
	}

	// Get company
	company, err := s.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return &TradeResult{
			Success: false,
			Message: "Company not found",
		}, nil
	}

	// Get current price
	currentPrice, err := s.stockRepo.GetLatestPrice(company.ID)
	if err != nil {
		return &TradeResult{
			Success: false,
			Message: "Price not available",
		}, nil
	}

	pricePerShare := currentPrice.ClosePrice
	totalAmount := float64(quantity) * pricePerShare

	// Check wallet balance
	wallet, err := s.GetOrCreateWallet(userID)
	if err != nil {
		return &TradeResult{
			Success: false,
			Message: "Failed to get wallet",
		}, nil
	}

	if wallet.Balance < totalAmount {
		return &TradeResult{
			Success: false,
			Message: fmt.Sprintf("Insufficient balance. Required: NPR %.2f, Available: NPR %.2f", totalAmount, wallet.Balance),
		}, nil
	}

	// Execute buy transaction
	if err := s.repo.ExecuteBuy(userID, company.ID, quantity, pricePerShare); err != nil {
		return &TradeResult{
			Success: false,
			Message: "Transaction failed",
		}, fmt.Errorf("buy transaction failed: %w", err)
	}

	// Get updated wallet balance
	updatedWallet, _ := s.repo.GetVirtualWallet(userID)

	return &TradeResult{
		Success:       true,
		Message:       "Stock purchased successfully",
		Quantity:      quantity,
		PricePerShare: pricePerShare,
		TotalAmount:   totalAmount,
		NewBalance:    updatedWallet.Balance,
	}, nil
}

func (s *service) SellStock(userID uint, symbol string, quantity int) (*TradeResult, error) {
	if quantity <= 0 {
		return &TradeResult{
			Success: false,
			Message: "Quantity must be greater than 0",
		}, nil
	}

	// Get company
	company, err := s.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return &TradeResult{
			Success: false,
			Message: "Company not found",
		}, nil
	}

	// Check if user owns the stock
	portfolio, err := s.repo.GetPortfolioItem(userID, company.ID)
	if err != nil {
		return &TradeResult{
			Success: false,
			Message: "You don't own this stock",
		}, nil
	}

	if portfolio.Quantity < quantity {
		return &TradeResult{
			Success: false,
			Message: fmt.Sprintf("Insufficient shares. You own %d shares", portfolio.Quantity),
		}, nil
	}

	// Get current price
	currentPrice, err := s.stockRepo.GetLatestPrice(company.ID)
	if err != nil {
		return &TradeResult{
			Success: false,
			Message: "Price not available",
		}, nil
	}

	pricePerShare := currentPrice.ClosePrice
	totalAmount := float64(quantity) * pricePerShare

	// Execute sell transaction
	if err := s.repo.ExecuteSell(userID, company.ID, quantity, pricePerShare); err != nil {
		return &TradeResult{
			Success: false,
			Message: "Transaction failed",
		}, fmt.Errorf("sell transaction failed: %w", err)
	}

	// Get updated wallet balance
	updatedWallet, _ := s.repo.GetVirtualWallet(userID)

	return &TradeResult{
		Success:       true,
		Message:       "Stock sold successfully",
		Quantity:      quantity,
		PricePerShare: pricePerShare,
		TotalAmount:   totalAmount,
		NewBalance:    updatedWallet.Balance,
	}, nil
}

func (s *service) GetTransactionHistory(userID uint, limit, offset int) ([]models.StockTransaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.repo.GetUserTransactions(userID, limit, offset)
}
