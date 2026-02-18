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


	to := time.Now()
	from := to.AddDate(0, 0, -30)

	prices, err := s.stockRepo.GetPriceHistory(company.ID, "1D", from, to, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch price history: %w", err)
	}

	if len(prices) < 10 {
		return nil, fmt.Errorf("insufficient data for prediction (need at least 10 points)")
	}

	

	var x []float64
	var y []float64


	n := len(prices)
	for i := n - 1; i >= 0; i-- {

		val := float64((n - 1) - i)
		price, _ := prices[i].ClosePrice.Float64()

		x = append(x, val)
		y = append(y, price)
	}

	lr := NewLinearRegression()
	lr.Train(x, y)

	
	nextDayX := float64(n)
	predictedFloat := lr.Predict(nextDayX)

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
