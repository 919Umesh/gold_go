package orderbook

import (
	"log/slog"

	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

// MatchResult contains the outcome of a single match between two orders
type MatchResult struct {
	Trade     *models.Trade
	BuyOrder  *models.Order
	SellOrder *models.Order
	FillQty   int64
	FillPrice decimal.Decimal
}

// Engine handles order matching with price-time priority
type Engine struct {
	repo Repository
}

func NewEngine(repo Repository) *Engine {
	return &Engine{repo: repo}
}

// MatchBuyOrder tries to match a buy order against existing sell orders
// Returns list of matched trades
func (e *Engine) MatchBuyOrder(buyOrder *models.Order) ([]MatchResult, error) {
	sellOrders, err := e.repo.GetOpenSellOrders(buyOrder.CompanyID)
	if err != nil {
		return nil, err
	}

	var results []MatchResult
	remainingQty := buyOrder.RemainingQty()

	for _, sell := range sellOrders {
		if remainingQty <= 0 {
			break
		}

		// For limit orders: buy price must be >= sell price
		if buyOrder.OrderType == models.OrderTypeLimit {
			if buyOrder.Price.LessThan(sell.Price) {
				// No more matches possible (sells are sorted by price ASC)
				break
			}
		}

		// Determine fill price (passive order's price = sell order's price)
		fillPrice := sell.Price

		// Determine fill quantity (minimum of remaining on both sides)
		sellRemaining := sell.RemainingQty()
		fillQty := min64(remainingQty, sellRemaining)

		if fillQty <= 0 {
			continue
		}

		// Calculate total amount
		totalAmount := fillPrice.Mul(decimal.NewFromInt(fillQty))

		// Create trade record
		trade := &models.Trade{
			CompanyID:   buyOrder.CompanyID,
			BuyOrderID:  buyOrder.ID,
			SellOrderID: sell.ID,
			BuyerID:     buyOrder.UserID,
			SellerID:    sell.UserID,
			Price:       fillPrice,
			Quantity:    fillQty,
			TotalAmount: totalAmount,
		}

		if err := e.repo.CreateTrade(trade); err != nil {
			slog.Error("failed to create trade", "error", err)
			continue
		}

		// Update sell order
		sell.FilledQty += fillQty
		if sell.FilledQty >= sell.Quantity {
			sell.Status = models.OrderStatusFilled
		} else {
			sell.Status = models.OrderStatusPartiallyFilled
		}
		if err := e.repo.UpdateOrder(&sell); err != nil {
			slog.Error("failed to update sell order", "error", err)
		}

		// Update buy order
		buyOrder.FilledQty += fillQty
		remainingQty -= fillQty

		results = append(results, MatchResult{
			Trade:     trade,
			BuyOrder:  buyOrder,
			SellOrder: &sell,
			FillQty:   fillQty,
			FillPrice: fillPrice,
		})
	}

	// Update buy order status
	if buyOrder.FilledQty >= buyOrder.Quantity {
		buyOrder.Status = models.OrderStatusFilled
	} else if buyOrder.FilledQty > 0 {
		buyOrder.Status = models.OrderStatusPartiallyFilled
	}
	if err := e.repo.UpdateOrder(buyOrder); err != nil {
		slog.Error("failed to update buy order", "error", err)
	}

	return results, nil
}

// MatchSellOrder tries to match a sell order against existing buy orders
func (e *Engine) MatchSellOrder(sellOrder *models.Order) ([]MatchResult, error) {
	buyOrders, err := e.repo.GetOpenBuyOrders(sellOrder.CompanyID)
	if err != nil {
		return nil, err
	}

	var results []MatchResult
	remainingQty := sellOrder.RemainingQty()

	for _, buy := range buyOrders {
		if remainingQty <= 0 {
			break
		}

		// For limit orders: buy price must be >= sell price
		if sellOrder.OrderType == models.OrderTypeLimit {
			if buy.Price.LessThan(sellOrder.Price) {
				break
			}
		}

		// Fill at the passive order's price (buy order's price)
		fillPrice := buy.Price

		buyRemaining := buy.RemainingQty()
		fillQty := min64(remainingQty, buyRemaining)

		if fillQty <= 0 {
			continue
		}

		totalAmount := fillPrice.Mul(decimal.NewFromInt(fillQty))

		trade := &models.Trade{
			CompanyID:   sellOrder.CompanyID,
			BuyOrderID:  buy.ID,
			SellOrderID: sellOrder.ID,
			BuyerID:     buy.UserID,
			SellerID:    sellOrder.UserID,
			Price:       fillPrice,
			Quantity:    fillQty,
			TotalAmount: totalAmount,
		}

		if err := e.repo.CreateTrade(trade); err != nil {
			slog.Error("failed to create trade", "error", err)
			continue
		}

		// Update buy order
		buy.FilledQty += fillQty
		if buy.FilledQty >= buy.Quantity {
			buy.Status = models.OrderStatusFilled
		} else {
			buy.Status = models.OrderStatusPartiallyFilled
		}
		if err := e.repo.UpdateOrder(&buy); err != nil {
			slog.Error("failed to update buy order", "error", err)
		}

		// Update sell order
		sellOrder.FilledQty += fillQty
		remainingQty -= fillQty

		results = append(results, MatchResult{
			Trade:     trade,
			BuyOrder:  &buy,
			SellOrder: sellOrder,
			FillQty:   fillQty,
			FillPrice: fillPrice,
		})
	}

	// Update sell order status
	if sellOrder.FilledQty >= sellOrder.Quantity {
		sellOrder.Status = models.OrderStatusFilled
	} else if sellOrder.FilledQty > 0 {
		sellOrder.Status = models.OrderStatusPartiallyFilled
	}
	if err := e.repo.UpdateOrder(sellOrder); err != nil {
		slog.Error("failed to update sell order", "error", err)
	}

	return results, nil
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
