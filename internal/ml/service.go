package ml

import (
	"fmt"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/shopspring/decimal"
)

type PredictionResult struct {
	Symbol         string          `json:"symbol"`
	PredictedPrice decimal.Decimal `json:"predicted_price"`
	RMSE           float64         `json:"rmse"`
	MAE            float64         `json:"mae"`
	Algorithm      string          `json:"algorithm"`
	Datapoints     int             `json:"datapoints"`
}

type Service interface {
	PredictStockPrice(symbol string, days int) (*PredictionResult, error)
}

type service struct {
	stockRepo stock.Repository
}

func NewService(stockRepo stock.Repository) Service {
	return &service{
		stockRepo: stockRepo,
	}
}

func (s *service) PredictStockPrice(symbol string, days int) (*PredictionResult, error) {
	company, err := s.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// Fetch 30 days of data for training
	to := time.Now()
	from := to.AddDate(0, 0, -30)

	// Fetch available data points (up to 30 days)
	prices, err := s.stockRepo.GetPriceHistory(company.ID, "1D", from, to, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch price history: %w", err)
	}

	if len(prices) < 10 {
		return nil, fmt.Errorf("insufficient data for prediction (need at least 10 points)")
	}

	// Prepare data for Linear Regression
	// X = Day index (0, 1, 2...)
	// Y = Close Price
	// We need to reverse prices because GetPriceHistory likely returns DESC order,
	// but for regression time should increase.
	// Actually typical time series regression uses time as X.

	var x []float64
	var y []float64

	// Prices are usually returned DESC (newest first). Let's verify or stick to sorting.
	// Repository says: query.OrderDesc("timestamp")
	// So index 0 is today, index N is oldest.
	// Let's create X such that oldest is 0.

	n := len(prices)
	for i := n - 1; i >= 0; i-- {
		// Oldest price at index n-1. We want that to be x=0
		// Index i (where i goes n-1 -> 0)
		// x value should be (n - 1) - i

		val := float64((n - 1) - i)
		price, _ := prices[i].ClosePrice.Float64()

		x = append(x, val)
		y = append(y, price)
	}

	lr := NewLinearRegression()
	lr.Train(x, y)

	// Predict for next day: x = n
	nextDayX := float64(n)
	predictedFloat := lr.Predict(nextDayX)

	// Calculate Error Metrics on training data (just for display)
	var predictedY []float64
	for _, val := range x {
		predictedY = append(predictedY, lr.Predict(val))
	}
	rmse := CalculateRMSE(y, predictedY)
	mae := CalculateMAE(y, predictedY)

	return &PredictionResult{
		Symbol:         symbol,
		PredictedPrice: decimal.NewFromFloat(predictedFloat).Round(2),
		RMSE:           rmse,
		MAE:            mae,
		Algorithm:      "Linear Regression (Simple)",
		Datapoints:     n,
	}, nil
}
