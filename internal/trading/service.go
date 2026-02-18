package trading

import (
	"errors"
	"fmt"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/models"

	"github.com/shopspring/decimal"
)

type Service interface {
	CreateWallet(userID string) (*models.VirtualWallet, error)
	GetWallet(userID string) (*models.VirtualWallet, error)
	GetOrCreateWallet(userID string) (*models.VirtualWallet, error)

	GetPortfolio(userID string) (*PortfolioSummary, error)

	BuyStock(userID, symbol string, quantity int) (*TradeResult, error)
	SellStock(userID, symbol string, quantity int) (*TradeResult, error)
	GetTransactionHistory(userID string, limit, offset int) ([]models.StockTransaction, error)
}

type PortfolioSummary struct {
	TotalValue      decimal.Decimal        `json:"total_value"`
	TotalInvested   decimal.Decimal        `json:"total_invested"`
	TotalProfitLoss decimal.Decimal        `json:"total_profit_loss"`
	ProfitLossPct   decimal.Decimal        `json:"profit_loss_pct"`
	Items           []PortfolioItemSummary `json:"items"`
}

type PortfolioItemSummary struct {
	models.UserPortfolio
	CurrentPrice  decimal.Decimal `json:"current_price"`
	CurrentValue  decimal.Decimal `json:"current_value"`
	ProfitLoss    decimal.Decimal `json:"profit_loss"`
	ProfitLossPct decimal.Decimal `json:"profit_loss_pct"`
	CompanyName   string          `json:"company_name"`
	CompanySymbol string          `json:"company_symbol"`
	CompanySector string          `json:"company_sector"`
}

type TradeResult struct {
	Success       bool            `json:"success"`
	Message       string          `json:"message"`
	TransactionID string          `json:"transaction_id,omitempty"`
	Quantity      int             `json:"quantity"`
	PricePerShare decimal.Decimal `json:"price_per_share"`
	TotalAmount   decimal.Decimal `json:"total_amount"`
	NewBalance    decimal.Decimal `json:"new_balance"`
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
	_, err := s.repo.GetVirtualWallet(userID)
	if err == nil {
		return nil, errors.New("wallet already exists")
	}

	wallet := &models.VirtualWallet{
		UserID:          userID,
		Balance:         decimal.Zero, 
		TotalInvested:   decimal.Zero,
		TotalProfitLoss: decimal.Zero,
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
		Items:         make([]PortfolioItemSummary, 0, len(portfolio)),
		TotalValue:    decimal.Zero,
		TotalInvested: decimal.Zero,
	}

	for _, item := range portfolio {

		price, err := s.stockRepo.GetLatestPrice(item.CompanyID)
		currentPrice := decimal.Zero
		if err == nil {
			currentPrice = price.ClosePrice
		}

		company, err := s.stockRepo.GetCompanyByID(item.CompanyID)
		companyName := ""
		companySymbol := ""
		companySector := ""
		if err == nil {
			companyName = company.Name
			companySymbol = company.Symbol
			companySector = company.Sector
		}

		currentValue := item.CalculateCurrentValue(currentPrice)
		profitLoss := item.CalculateProfitLoss(currentPrice)
		profitLossPct := item.CalculateProfitLossPercentage(currentPrice)

		summary.TotalValue = summary.TotalValue.Add(currentValue)
		summary.TotalInvested = summary.TotalInvested.Add(item.TotalInvested)

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

	summary.TotalProfitLoss = summary.TotalValue.Sub(summary.TotalInvested)

	if summary.TotalInvested.IsPositive() {
		summary.ProfitLossPct = summary.TotalProfitLoss.Div(summary.TotalInvested).Mul(decimal.NewFromInt(100))
	} else {
		summary.ProfitLossPct = decimal.Zero
	}

	return summary, nil
}

func (s *service) BuyStock(userID, symbol string, quantity int) (*TradeResult, error) {
	if quantity <= 0 {
		return &TradeResult{Success: false, Message: "Quantity must be positive"}, nil
	}

	company, err := s.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return &TradeResult{Success: false, Message: "Company not found"}, nil
	}

	if company.AvailableShares < int64(quantity) {
		return &TradeResult{
			Success: false,
			Message: fmt.Sprintf("Insufficient shares available. Only %d shares in market", company.AvailableShares),
		}, nil
	}

	price, err := s.stockRepo.GetLatestPrice(company.ID)
	if err != nil {
		initialPrice := decimal.NewFromInt(100)
		initialStockPrice := &models.StockPrice{
			CompanyID:  company.ID,
			OpenPrice:  initialPrice,
			HighPrice:  initialPrice,
			LowPrice:   initialPrice,
			ClosePrice: initialPrice,
			Volume:     0,
			Timestamp:  time.Now(),
			Timeframe:  "1D",
		}
		if createErr := s.stockRepo.CreateStockPrice(initialStockPrice); createErr != nil {
			return nil, fmt.Errorf("failed to create initial stock price: %w", createErr)
		}
		price = initialStockPrice
	}

	totalAmount := price.ClosePrice.Mul(decimal.NewFromInt(int64(quantity)))

	err = s.repo.ExecuteBuy(userID, company.ID, quantity, price.ClosePrice)
	if err != nil {
		return &TradeResult{Success: false, Message: err.Error()}, nil
	}

	newAvailableShares := company.AvailableShares - int64(quantity)
	if err := s.stockRepo.UpdateCompanyShares(company.ID, newAvailableShares); err != nil {
		fmt.Printf("Warning: failed to update available shares: %v\n", err)
	}

	if err := s.stockRepo.UpdateStockPriceVolume(price.ID, int64(quantity)); err != nil {
		fmt.Printf("Warning: failed to update volume: %v\n", err)
	}

	wallet, _ := s.repo.GetVirtualWallet(userID)

	return &TradeResult{
		Success:       true,
		Message:       "Stock purchased successfully",
		Quantity:      quantity,
		PricePerShare: price.ClosePrice,
		TotalAmount:   totalAmount,
		NewBalance:    wallet.Balance,
	}, nil
}

func (s *service) SellStock(userID, symbol string, quantity int) (*TradeResult, error) {
	if quantity <= 0 {
		return &TradeResult{Success: false, Message: "Quantity must be positive"}, nil
	}

	company, err := s.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return &TradeResult{Success: false, Message: "Company not found"}, nil
	}

	price, err := s.stockRepo.GetLatestPrice(company.ID)
	if err != nil {
		initialPrice := decimal.NewFromInt(100)
		initialStockPrice := &models.StockPrice{
			CompanyID:  company.ID,
			OpenPrice:  initialPrice,
			HighPrice:  initialPrice,
			LowPrice:   initialPrice,
			ClosePrice: initialPrice,
			Volume:     0,
			Timestamp:  time.Now(),
			Timeframe:  "1D",
		}
		if createErr := s.stockRepo.CreateStockPrice(initialStockPrice); createErr != nil {
			return nil, fmt.Errorf("failed to create initial stock price: %w", createErr)
		}
		price = initialStockPrice
	}

	totalAmount := price.ClosePrice.Mul(decimal.NewFromInt(int64(quantity)))

	err = s.repo.ExecuteSell(userID, company.ID, quantity, price.ClosePrice)
	if err != nil {
		return &TradeResult{Success: false, Message: err.Error()}, nil
	}

	newAvailableShares := company.AvailableShares + int64(quantity)
	if newAvailableShares > company.TotalShares {
		newAvailableShares = company.TotalShares
	}
	if err := s.stockRepo.UpdateCompanyShares(company.ID, newAvailableShares); err != nil {
		fmt.Printf("Warning: failed to update available shares: %v\n", err)
	}

	if err := s.stockRepo.UpdateStockPriceVolume(price.ID, int64(quantity)); err != nil {
		fmt.Printf("Warning: failed to update volume: %v\n", err)
	}

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
	return s.repo.GetUserTransactions(userID, limit, offset)
}
