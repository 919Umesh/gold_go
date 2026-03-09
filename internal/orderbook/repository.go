package orderbook

import (
	"fmt"

	"github.com/919Umesh/stock_market_sim/internal/supabase"
	"github.com/919Umesh/stock_market_sim/models"
)

type Repository interface {
	// Orders
	CreateOrder(order *models.Order) error
	GetOrderByID(orderID string) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	GetOpenSellOrders(companyID string) ([]models.Order, error)
	GetOpenBuyOrders(companyID string) ([]models.Order, error)
	GetUserOrders(userID string, limit int) ([]models.Order, error)
	CancelOrder(orderID string) error

	// Trades
	CreateTrade(trade *models.Trade) error
	GetTradesByCompany(companyID string, limit int) ([]models.Trade, error)
	GetUserTrades(userID string, limit int) ([]models.Trade, error)

	// Portfolio
	GetPortfolioItem(userID, companyID string) (*models.Portfolio, error)
	CreatePortfolioItem(item *models.Portfolio) error
	UpdatePortfolioItem(item *models.Portfolio) error
	DeletePortfolioItem(portfolioID string) error
	GetPortfolio(userID string) ([]models.Portfolio, error)
}

type repository struct {
	client *supabase.Client
}

func NewRepository(client *supabase.Client) Repository {
	return &repository{client: client}
}

// ──────────────────── Orders ────────────────────

func (r *repository) CreateOrder(order *models.Order) error {
	query := `INSERT INTO orders (user_id, company_id, side, order_type, price, quantity, filled_qty, status)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	return r.client.ExecuteInsert(query, order,
		order.UserID, order.CompanyID, order.Side, order.OrderType,
		order.Price.String(), order.Quantity, order.FilledQty, order.Status)
}

func (r *repository) GetOrderByID(orderID string) (*models.Order, error) {
	var order models.Order
	err := r.client.ExecuteQueryRow("SELECT * FROM orders WHERE id = $1", &order, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	return &order, nil
}

func (r *repository) UpdateOrder(order *models.Order) error {
	query := `UPDATE orders SET filled_qty = $1, status = $2 WHERE id = $3 RETURNING *`
	return r.client.ExecuteUpdate(query, order, order.FilledQty, order.Status, order.ID)
}

// GetOpenSellOrders returns sell orders sorted by price ASC, then time ASC (best price first)
func (r *repository) GetOpenSellOrders(companyID string) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT * FROM orders 
			  WHERE company_id = $1 AND side = $2 AND (status = $3 OR status = $4) 
			  ORDER BY price LIMIT $5`
	err := r.client.ExecuteQuery(query, &orders, companyID, models.OrderSideSell,
		models.OrderStatusOpen, models.OrderStatusPartiallyFilled, 100)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		orders = []models.Order{}
	}
	return orders, nil
}

// GetOpenBuyOrders returns buy orders sorted by price DESC, then time ASC (best price first)
func (r *repository) GetOpenBuyOrders(companyID string) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT * FROM orders 
			  WHERE company_id = $1 AND side = $2 AND (status = $3 OR status = $4) 
			  ORDER BY price DESC LIMIT $5`
	err := r.client.ExecuteQuery(query, &orders, companyID, models.OrderSideBuy,
		models.OrderStatusOpen, models.OrderStatusPartiallyFilled, 100)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		orders = []models.Order{}
	}
	return orders, nil
}

func (r *repository) GetUserOrders(userID string, limit int) ([]models.Order, error) {
	var orders []models.Order
	query := "SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2"
	err := r.client.ExecuteQuery(query, &orders, userID, limit)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		orders = []models.Order{}
	}
	return orders, nil
}

func (r *repository) CancelOrder(orderID string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2 RETURNING *`
	var order models.Order
	return r.client.ExecuteUpdate(query, &order, models.OrderStatusCancelled, orderID)
}

// ──────────────────── Trades ────────────────────

func (r *repository) CreateTrade(trade *models.Trade) error {
	query := `INSERT INTO trades (company_id, buy_order_id, sell_order_id, buyer_id, seller_id, price, quantity, total_amount)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	return r.client.ExecuteInsert(query, trade,
		trade.CompanyID, trade.BuyOrderID, trade.SellOrderID,
		trade.BuyerID, trade.SellerID, trade.Price.String(),
		trade.Quantity, trade.TotalAmount.String())
}

func (r *repository) GetTradesByCompany(companyID string, limit int) ([]models.Trade, error) {
	var trades []models.Trade
	query := "SELECT * FROM trades WHERE company_id = $1 ORDER BY created_at DESC LIMIT $2"
	err := r.client.ExecuteQuery(query, &trades, companyID, limit)
	if err != nil {
		return nil, err
	}
	if trades == nil {
		trades = []models.Trade{}
	}
	return trades, nil
}

func (r *repository) GetUserTrades(userID string, limit int) ([]models.Trade, error) {
	var trades []models.Trade
	query := `SELECT * FROM trades WHERE (buyer_id = $1 OR seller_id = $2) ORDER BY created_at DESC LIMIT $3`
	err := r.client.ExecuteQuery(query, &trades, userID, userID, limit)
	if err != nil {
		return nil, err
	}
	if trades == nil {
		trades = []models.Trade{}
	}
	return trades, nil
}

// ──────────────────── Portfolio ────────────────────

func (r *repository) GetPortfolioItem(userID, companyID string) (*models.Portfolio, error) {
	var item models.Portfolio
	query := "SELECT * FROM portfolios WHERE user_id = $1 AND company_id = $2"
	err := r.client.ExecuteQueryRow(query, &item, userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("portfolio item not found")
	}
	return &item, nil
}

func (r *repository) CreatePortfolioItem(item *models.Portfolio) error {
	query := `INSERT INTO portfolios (user_id, company_id, quantity, avg_buy_price)
			  VALUES ($1, $2, $3, $4) RETURNING *`
	return r.client.ExecuteInsert(query, item,
		item.UserID, item.CompanyID, item.Quantity, item.AvgBuyPrice.String())
}

func (r *repository) UpdatePortfolioItem(item *models.Portfolio) error {
	query := `UPDATE portfolios SET quantity = $1, avg_buy_price = $2 WHERE id = $3 RETURNING *`
	return r.client.ExecuteUpdate(query, item, item.Quantity, item.AvgBuyPrice.String(), item.ID)
}

func (r *repository) DeletePortfolioItem(portfolioID string) error {
	query := "DELETE FROM portfolios WHERE id = $1"
	return r.client.ExecuteDelete(query, portfolioID)
}

func (r *repository) GetPortfolio(userID string) ([]models.Portfolio, error) {
	var items []models.Portfolio
	query := "SELECT * FROM portfolios WHERE user_id = $1 AND quantity > $2"
	err := r.client.ExecuteQuery(query, &items, userID, 0)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.Portfolio{}
	}
	return items, nil
}
