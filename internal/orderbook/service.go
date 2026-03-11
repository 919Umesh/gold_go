package orderbook

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/internal/wallet"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

var (
	ErrInsufficientShares = errors.New("insufficient shares in portfolio")
	ErrOrderNotFound      = errors.New("order not found")
	ErrNotOrderOwner      = errors.New("you don't own this order")
	ErrCannotCancel       = errors.New("order is already filled or cancelled")
)

type OrderBookLevel struct {
	Price    decimal.Decimal `json:"price"`
	Quantity int64           `json:"quantity"`
	Orders   int             `json:"orders"`
}

type OrderBookView struct {
	CompanyID string           `json:"company_id"`
	Bids      []OrderBookLevel `json:"bids"`
	Asks      []OrderBookLevel `json:"asks"`
}

type Service interface {
	PlaceSellOrder(userID, companyID string, qty int64, price decimal.Decimal) (*models.Order, []MatchResult, error)
	PlaceBuyOrder(userID, companyID string, qty int64, price decimal.Decimal, orderType string) (*models.Order, []MatchResult, error)
	CancelOrder(userID, orderID string) error
	GetOrderBook(companyID string) (*OrderBookView, error)
	GetUserOrders(userID string, limit int) ([]models.Order, error)
	GetPortfolio(userID string) ([]models.Portfolio, error)
	GetUserTrades(userID string, limit int) ([]models.Trade, error)
	GetCompanyTrades(companyID string, limit int) ([]models.Trade, error)
	GetCompanyTransactions(companyID string, limit int) ([]PublicTransaction, error)
}

type PublicTransaction struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"` // buy, sell
	Price     decimal.Decimal `json:"price"`
	Quantity  int64           `json:"quantity"`
	CreatedAt time.Time       `json:"created_at"`
}

type service struct {
	repo        Repository
	engine      *Engine
	walletSvc   wallet.Service
	stockRepo   stock.Repository
	priceEngine *market.PriceEngine
}

func NewService(repo Repository, walletSvc wallet.Service, stockRepo stock.Repository, priceEngine *market.PriceEngine) Service {
	return &service{
		repo:        repo,
		engine:      NewEngine(repo),
		walletSvc:   walletSvc,
		stockRepo:   stockRepo,
		priceEngine: priceEngine,
	}
}

// ──────────────────── Place Sell Order ────────────────────

func (s *service) PlaceSellOrder(userID, companyID string, qty int64, price decimal.Decimal) (*models.Order, []MatchResult, error) {
	// Validate: user must have sufficient shares in portfolio
	portfolio, err := s.repo.GetPortfolioItem(userID, companyID)
	if err != nil {
		return nil, nil, ErrInsufficientShares
	}

	if portfolio.Quantity < qty {
		return nil, nil, ErrInsufficientShares
	}

	// Deduct shares from portfolio (lock them for the sell order)
	portfolio.Quantity -= qty
	if err := s.repo.UpdatePortfolioItem(portfolio); err != nil {
		return nil, nil, fmt.Errorf("failed to update portfolio: %w", err)
	}

	// Create sell order
	order := &models.Order{
		UserID:    userID,
		CompanyID: companyID,
		Side:      models.OrderSideSell,
		OrderType: models.OrderTypeLimit,
		Price:     price,
		Quantity:  qty,
		FilledQty: 0,
		Status:    models.OrderStatusOpen,
	}

	if err := s.repo.CreateOrder(order); err != nil {
		// Rollback portfolio
		portfolio.Quantity += qty
		_ = s.repo.UpdatePortfolioItem(portfolio)
		return nil, nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Try to match the sell order against existing buy orders
	matches, err := s.engine.MatchSellOrder(order)
	if err != nil {
		slog.Error("matching failed", "error", err)
	}

	// Process matched trades
	for _, m := range matches {
		s.processTradeSettlement(m)
	}

	return order, matches, nil
}

// ──────────────────── Place Buy Order ────────────────────

func (s *service) PlaceBuyOrder(userID, companyID string, qty int64, price decimal.Decimal, orderType string) (*models.Order, []MatchResult, error) {
	if orderType == "" {
		orderType = models.OrderTypeLimit
	}

	// Calculate total cost and lock funds
	totalCost := price.Mul(decimal.NewFromInt(qty))
	if err := s.walletSvc.LockFunds(userID, totalCost); err != nil {
		return nil, nil, fmt.Errorf("insufficient trading wallet balance: %w", err)
	}

	// Create buy order
	order := &models.Order{
		UserID:    userID,
		CompanyID: companyID,
		Side:      models.OrderSideBuy,
		OrderType: orderType,
		Price:     price,
		Quantity:  qty,
		FilledQty: 0,
		Status:    models.OrderStatusOpen,
	}

	if err := s.repo.CreateOrder(order); err != nil {
		_ = s.walletSvc.ReleaseFunds(userID, totalCost)
		return nil, nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Try to match against existing sell orders
	matches, err := s.engine.MatchBuyOrder(order)
	if err != nil {
		slog.Error("matching failed", "error", err)
	}

	// Process matched trades
	for _, m := range matches {
		s.processTradeSettlement(m)
	}

	// Release excess locked funds if partially filled at a better price
	if order.FilledQty > 0 {
		// Calculate actual cost of fills
		var actualCost decimal.Decimal
		for _, m := range matches {
			actualCost = actualCost.Add(m.FillPrice.Mul(decimal.NewFromInt(m.FillQty)))
		}
		// Release difference
		unfilledQty := order.RemainingQty()
		unfilledCost := price.Mul(decimal.NewFromInt(unfilledQty))
		excess := totalCost.Sub(actualCost).Sub(unfilledCost)
		if excess.IsPositive() {
			_ = s.walletSvc.ReleaseFunds(userID, excess)
		}
	}

	return order, matches, nil
}

// ──────────────────── Trade Settlement ────────────────────

func (s *service) processTradeSettlement(m MatchResult) {
	// 1. Deduct locked funds from buyer
	if err := s.walletSvc.DeductLockedFunds(m.BuyOrder.UserID, m.Trade.TotalAmount); err != nil {
		slog.Error("failed to deduct buyer funds", "error", err)
	}

	// 2. Credit seller's trading wallet
	if err := s.walletSvc.CreditTradingWallet(m.SellOrder.UserID, m.Trade.TotalAmount); err != nil {
		slog.Error("failed to credit seller", "error", err)
	}

	// 3. Credit shares to buyer's portfolio
	s.creditBuyerPortfolio(m.BuyOrder.UserID, m.Trade.CompanyID, m.FillQty, m.FillPrice)

	// 4. Update company price through price engine
	if s.priceEngine != nil {
		s.priceEngine.ProcessMatchedTrade(m.Trade.CompanyID, m.FillPrice, m.FillQty)
	}
}

func (s *service) creditBuyerPortfolio(userID, companyID string, qty int64, price decimal.Decimal) {
	existing, err := s.repo.GetPortfolioItem(userID, companyID)
	if err != nil {
		// Create new
		item := &models.Portfolio{
			UserID:      userID,
			CompanyID:   companyID,
			Quantity:    qty,
			AvgBuyPrice: price,
		}
		if createErr := s.repo.CreatePortfolioItem(item); createErr != nil {
			slog.Error("failed to create portfolio", "error", createErr)
		}
		return
	}

	// Update weighted average
	totalInvested := existing.AvgBuyPrice.Mul(decimal.NewFromInt(existing.Quantity))
	newInvested := price.Mul(decimal.NewFromInt(qty))
	newQty := existing.Quantity + qty
	newAvg := totalInvested.Add(newInvested).Div(decimal.NewFromInt(newQty))

	existing.Quantity = newQty
	existing.AvgBuyPrice = newAvg
	if err := s.repo.UpdatePortfolioItem(existing); err != nil {
		slog.Error("failed to update portfolio", "error", err)
	}
}

// ──────────────────── Cancel Order ────────────────────

func (s *service) CancelOrder(userID, orderID string) error {
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return ErrOrderNotFound
	}

	if order.UserID != userID {
		return ErrNotOrderOwner
	}

	if order.Status == models.OrderStatusFilled || order.Status == models.OrderStatusCancelled {
		return ErrCannotCancel
	}

	remainingQty := order.RemainingQty()

	if order.Side == models.OrderSideBuy {
		// Release locked funds for unfilled portion
		refund := order.Price.Mul(decimal.NewFromInt(remainingQty))
		_ = s.walletSvc.ReleaseFunds(userID, refund)
	} else {
		// Return shares to portfolio preserving original avg_buy_price
		existing, pErr := s.repo.GetPortfolioItem(userID, order.CompanyID)
		if pErr != nil {
			// Portfolio row was deleted; recreate with order price as best estimate
			item := &models.Portfolio{
				UserID:      userID,
				CompanyID:   order.CompanyID,
				Quantity:    remainingQty,
				AvgBuyPrice: order.Price,
			}
			_ = s.repo.CreatePortfolioItem(item)
		} else {
			// Just add quantity back without changing avg_buy_price
			existing.Quantity += remainingQty
			_ = s.repo.UpdatePortfolioItem(existing)
		}
	}

	return s.repo.CancelOrder(orderID)
}

// ──────────────────── Order Book View ────────────────────

func (s *service) GetOrderBook(companyID string) (*OrderBookView, error) {
	sellOrders, err := s.repo.GetOpenSellOrders(companyID)
	if err != nil {
		return nil, err
	}

	buyOrders, err := s.repo.GetOpenBuyOrders(companyID)
	if err != nil {
		return nil, err
	}

	// Aggregate by price level
	asks := aggregateLevels(sellOrders)
	bids := aggregateLevels(buyOrders)

	return &OrderBookView{
		CompanyID: companyID,
		Bids:      bids,
		Asks:      asks,
	}, nil
}

func aggregateLevels(orders []models.Order) []OrderBookLevel {
	levelMap := make(map[string]*OrderBookLevel)
	var levelKeys []string

	for _, o := range orders {
		key := o.Price.String()
		if lvl, ok := levelMap[key]; ok {
			lvl.Quantity += o.RemainingQty()
			lvl.Orders++
		} else {
			levelMap[key] = &OrderBookLevel{
				Price:    o.Price,
				Quantity: o.RemainingQty(),
				Orders:   1,
			}
			levelKeys = append(levelKeys, key)
		}
	}

	levels := make([]OrderBookLevel, 0, len(levelKeys))
	for _, k := range levelKeys {
		levels = append(levels, *levelMap[k])
	}
	return levels
}

// ──────────────────── Queries ────────────────────

func (s *service) GetUserOrders(userID string, limit int) ([]models.Order, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetUserOrders(userID, limit)
}

func (s *service) GetPortfolio(userID string) ([]models.Portfolio, error) {
	return s.repo.GetPortfolio(userID)
}

func (s *service) GetUserTrades(userID string, limit int) ([]models.Trade, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetUserTrades(userID, limit)
}

func (s *service) GetCompanyTrades(companyID string, limit int) ([]models.Trade, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetTradesByCompany(companyID, limit)
}

func (s *service) GetCompanyTransactions(companyID string, limit int) ([]PublicTransaction, error) {
	if limit <= 0 {
		limit = 50
	}

	trades, err := s.repo.GetTradesByCompany(companyID, limit)
	if err != nil {
		return nil, err
	}

	transactions := make([]PublicTransaction, 0, len(trades)*2)
	for _, t := range trades {
		// Buyer's side (Buy transaction)
		transactions = append(transactions, PublicTransaction{
			ID:        t.ID + "-buy",
			Type:      "buy",
			Price:     t.Price,
			Quantity:  t.Quantity,
			CreatedAt: t.CreatedAt,
		})
		// Seller's side (Sell transaction)
		transactions = append(transactions, PublicTransaction{
			ID:        t.ID + "-sell",
			Type:      "sell",
			Price:     t.Price,
			Quantity:  t.Quantity,
			CreatedAt: t.CreatedAt,
		})
	}

	return transactions, nil
}
