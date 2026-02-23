package ml

import (
	"fmt"
	"math"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/shopspring/decimal"
)

// PredictionResult holds prediction from a single algorithm
type PredictionResult struct {
	Symbol         string          `json:"symbol"`
	CurrentPrice   decimal.Decimal `json:"current_price"`
	PredictedPrice decimal.Decimal `json:"predicted_price"`
	Change         decimal.Decimal `json:"predicted_change"`
	ChangePct      decimal.Decimal `json:"predicted_change_pct"`
	RMSE           float64         `json:"rmse"`
	MAE            float64         `json:"mae"`
	Algorithm      string          `json:"algorithm"`
	Datapoints     int             `json:"datapoints"`
}

// ComparisonResult holds predictions from all algorithms for comparison
type ComparisonResult struct {
	Symbol       string               `json:"symbol"`
	CurrentPrice decimal.Decimal      `json:"current_price"`
	Datapoints   int                  `json:"datapoints"`
	Ensemble     *PredictionResult    `json:"ensemble"`
	Models       []PredictionResult   `json:"models"`
	BestModel    string               `json:"best_model"`
	ModelDetails []EnsemblePrediction `json:"model_details"`
}

// Service provides ML prediction capabilities
type Service interface {
	// PredictStockPrice predicts using a specific algorithm
	// algorithm: "linear_regression", "ema", "holt", "knn", "ar", "wma", "ensemble"
	PredictStockPrice(symbol string, days int) (*PredictionResult, error)

	// PredictWithAlgorithm predicts using a specific algorithm
	PredictWithAlgorithm(symbol string, algorithm string) (*PredictionResult, error)

	// CompareAlgorithms runs all algorithms and returns comparison
	CompareAlgorithms(symbol string) (*ComparisonResult, error)
}

type service struct {
	stockRepo stock.Repository
}

// NewService creates a new ML prediction service
func NewService(stockRepo stock.Repository) Service {
	return &service{
		stockRepo: stockRepo,
	}
}

// fetchPriceData loads historical price data and converts to float64 slice
func (s *service) fetchPriceData(symbol string) ([]float64, decimal.Decimal, int, error) {
	company, err := s.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return nil, decimal.Zero, 0, fmt.Errorf("company not found: %w", err)
	}

	to := time.Now()
	from := to.AddDate(0, 0, -90) // Use 90 days for better predictions

	prices, err := s.stockRepo.GetPriceHistory(company.ID, "1D", from, to, 90)
	if err != nil {
		return nil, decimal.Zero, 0, fmt.Errorf("failed to fetch price history: %w", err)
	}

	if len(prices) < 10 {
		return nil, decimal.Zero, 0, fmt.Errorf("insufficient data for prediction (need at least 10 data points, got %d)", len(prices))
	}

	// Convert to chronological order (oldest first)
	n := len(prices)
	priceSlice := make([]float64, n)
	for i := 0; i < n; i++ {
		val, _ := prices[n-1-i].ClosePrice.Float64()
		priceSlice[i] = val
	}

	currentPrice := prices[0].ClosePrice // Latest price (prices are DESC from DB)

	return priceSlice, currentPrice, n, nil
}

// buildResult creates a PredictionResult from predicted value
func buildResult(symbol string, currentPrice decimal.Decimal, predicted float64, rmse, mae float64, algorithm string, datapoints int) *PredictionResult {
	predictedDec := decimal.NewFromFloat(predicted).Round(2)
	change := predictedDec.Sub(currentPrice)
	changePct := decimal.Zero
	if !currentPrice.IsZero() {
		changePct = change.Div(currentPrice).Mul(decimal.NewFromInt(100)).Round(2)
	}

	return &PredictionResult{
		Symbol:         symbol,
		CurrentPrice:   currentPrice,
		PredictedPrice: predictedDec,
		Change:         change,
		ChangePct:      changePct,
		RMSE:           math.Round(rmse*10000) / 10000,
		MAE:            math.Round(mae*10000) / 10000,
		Algorithm:      algorithm,
		Datapoints:     datapoints,
	}
}

// PredictStockPrice predicts using the ensemble method (default, backward compatible)
func (s *service) PredictStockPrice(symbol string, days int) (*PredictionResult, error) {
	return s.PredictWithAlgorithm(symbol, "ensemble")
}

// PredictWithAlgorithm predicts using a specific algorithm
func (s *service) PredictWithAlgorithm(symbol string, algorithm string) (*PredictionResult, error) {
	priceSlice, currentPrice, n, err := s.fetchPriceData(symbol)
	if err != nil {
		return nil, err
	}

	switch algorithm {
	case "linear_regression":
		return s.predictLinearRegression(symbol, priceSlice, currentPrice, n)
	case "ema":
		return s.predictEMA(symbol, priceSlice, currentPrice, n)
	case "holt":
		return s.predictHolt(symbol, priceSlice, currentPrice, n)
	case "knn":
		return s.predictKNN(symbol, priceSlice, currentPrice, n)
	case "ar":
		return s.predictAR(symbol, priceSlice, currentPrice, n)
	case "wma":
		return s.predictWMA(symbol, priceSlice, currentPrice, n)
	case "ensemble":
		return s.predictEnsemble(symbol, priceSlice, currentPrice, n)
	default:
		return s.predictEnsemble(symbol, priceSlice, currentPrice, n)
	}
}

// CompareAlgorithms runs all algorithms and returns comparison
func (s *service) CompareAlgorithms(symbol string) (*ComparisonResult, error) {
	priceSlice, currentPrice, n, err := s.fetchPriceData(symbol)
	if err != nil {
		return nil, err
	}

	algorithms := []string{"linear_regression", "ema", "holt", "knn", "ar", "wma"}
	var models []PredictionResult

	for _, algo := range algorithms {
		var result *PredictionResult
		switch algo {
		case "linear_regression":
			result, err = s.predictLinearRegression(symbol, priceSlice, currentPrice, n)
		case "ema":
			result, err = s.predictEMA(symbol, priceSlice, currentPrice, n)
		case "holt":
			result, err = s.predictHolt(symbol, priceSlice, currentPrice, n)
		case "knn":
			result, err = s.predictKNN(symbol, priceSlice, currentPrice, n)
		case "ar":
			result, err = s.predictAR(symbol, priceSlice, currentPrice, n)
		case "wma":
			result, err = s.predictWMA(symbol, priceSlice, currentPrice, n)
		}
		if err == nil && result != nil {
			models = append(models, *result)
		}
	}

	// Run ensemble
	ensembleResult, err := s.predictEnsemble(symbol, priceSlice, currentPrice, n)
	if err != nil {
		return nil, err
	}

	// Get ensemble details
	_, ensembleDetails := EnsemblePredict(priceSlice)

	// Find best model (lowest RMSE among those with RMSE > 0)
	bestModel := ""
	bestRMSE := math.MaxFloat64
	for _, m := range models {
		if m.RMSE > 0 && m.RMSE < bestRMSE {
			bestRMSE = m.RMSE
			bestModel = m.Algorithm
		}
	}

	return &ComparisonResult{
		Symbol:       symbol,
		CurrentPrice: currentPrice,
		Datapoints:   n,
		Ensemble:     ensembleResult,
		Models:       models,
		BestModel:    bestModel,
		ModelDetails: ensembleDetails,
	}, nil
}

// =====================================================
// Individual algorithm implementations
// =====================================================

func (s *service) predictLinearRegression(symbol string, prices []float64, currentPrice decimal.Decimal, n int) (*PredictionResult, error) {
	x := make([]float64, n)
	for i := range x {
		x[i] = float64(i)
	}

	lr := NewLinearRegression()
	lr.Train(x, prices)
	predicted := lr.Predict(float64(n))

	fitted := make([]float64, n)
	for i := range x {
		fitted[i] = lr.Predict(x[i])
	}

	return buildResult(symbol, currentPrice, predicted,
		CalculateRMSE(prices, fitted),
		CalculateMAE(prices, fitted),
		"Linear Regression", n), nil
}

func (s *service) predictEMA(symbol string, prices []float64, currentPrice decimal.Decimal, n int) (*PredictionResult, error) {
	ema := NewMovingAverage(5)
	predicted := ema.PredictEMA(prices)
	fitted := ema.GetEMASeries(prices)

	return buildResult(symbol, currentPrice, predicted,
		CalculateRMSE(prices, fitted),
		CalculateMAE(prices, fitted),
		"Exponential Moving Average (EMA-5)", n), nil
}

func (s *service) predictHolt(symbol string, prices []float64, currentPrice decimal.Decimal, n int) (*PredictionResult, error) {
	holt := NewHoltLinear(0.5, 0.3)
	predicted := holt.Predict(prices, 1)
	fitted := holt.GetFittedValues(prices)

	return buildResult(symbol, currentPrice, predicted,
		CalculateRMSE(prices, fitted),
		CalculateMAE(prices, fitted),
		"Holt's Linear Trend (Double Exponential Smoothing)", n), nil
}

func (s *service) predictKNN(symbol string, prices []float64, currentPrice decimal.Decimal, n int) (*PredictionResult, error) {
	knn := NewKNNRegressor(5, 5, true)
	predicted := knn.Predict(prices)

	// KNN doesn't produce a full fitted series easily, so approximate error
	// by leave-one-out on recent data
	if n > 10 {
		trainPrices := prices[:n-1]
		actual := prices[n-1]
		knnTest := NewKNNRegressor(5, 5, true)
		testPred := knnTest.Predict(trainPrices)
		rmse := math.Abs(actual - testPred)
		return buildResult(symbol, currentPrice, predicted, rmse, rmse,
			"K-Nearest Neighbors (KNN-5, window=5)", n), nil
	}

	return buildResult(symbol, currentPrice, predicted, 0, 0,
		"K-Nearest Neighbors (KNN-5, window=5)", n), nil
}

func (s *service) predictAR(symbol string, prices []float64, currentPrice decimal.Decimal, n int) (*PredictionResult, error) {
	ar := NewAutoRegressive(3)
	ar.Train(prices)
	predicted := ar.Predict(prices)
	fitted := ar.GetFittedValues(prices)

	// Calculate error only on fitted portion (skip first p values)
	p := ar.Order
	if p >= n {
		p = 0
	}

	return buildResult(symbol, currentPrice, predicted,
		CalculateRMSE(prices[p:], fitted[p:]),
		CalculateMAE(prices[p:], fitted[p:]),
		"Auto-Regressive AR(3)", n), nil
}

func (s *service) predictWMA(symbol string, prices []float64, currentPrice decimal.Decimal, n int) (*PredictionResult, error) {
	wma := NewWeightedMovingAverage(5)
	predicted := wma.Predict(prices)

	// Compute rolling WMA for error estimation
	if n > 5 {
		fitted := make([]float64, n)
		for i := 0; i < 5; i++ {
			fitted[i] = prices[i]
		}
		for i := 5; i < n; i++ {
			subWma := NewWeightedMovingAverage(5)
			fitted[i] = subWma.Predict(prices[:i])
		}
		return buildResult(symbol, currentPrice, predicted,
			CalculateRMSE(prices[5:], fitted[5:]),
			CalculateMAE(prices[5:], fitted[5:]),
			"Weighted Moving Average (WMA-5)", n), nil
	}

	return buildResult(symbol, currentPrice, predicted, 0, 0,
		"Weighted Moving Average (WMA-5)", n), nil
}

func (s *service) predictEnsemble(symbol string, prices []float64, currentPrice decimal.Decimal, n int) (*PredictionResult, error) {
	ensemblePred, _ := EnsemblePredict(prices)

	return buildResult(symbol, currentPrice, ensemblePred, 0, 0,
		"Ensemble (Weighted Average of All Models)", n), nil
}
