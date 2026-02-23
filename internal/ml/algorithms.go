package ml

import (
	"math"
	"sort"
)

// =====================================================
// MOVING AVERAGE (SMA & EMA)
// One of the most widely used technical analysis methods.
// SMA: Simple average of last N prices
// EMA: Exponentially weighted — recent prices have more influence
// =====================================================

// MovingAverage implements Simple and Exponential Moving Average prediction
type MovingAverage struct {
	Window int
}

// NewMovingAverage creates a new MovingAverage predictor
func NewMovingAverage(window int) *MovingAverage {
	if window <= 0 {
		window = 5
	}
	return &MovingAverage{Window: window}
}

// PredictSMA predicts the next value using Simple Moving Average
// Formula: SMA = (P1 + P2 + ... + Pn) / n
func (ma *MovingAverage) PredictSMA(prices []float64) float64 {
	n := len(prices)
	if n == 0 {
		return 0
	}

	window := ma.Window
	if window > n {
		window = n
	}

	sum := 0.0
	for i := n - window; i < n; i++ {
		sum += prices[i]
	}
	return sum / float64(window)
}

// PredictEMA predicts the next value using Exponential Moving Average
// Formula: EMA_t = α × Price_t + (1 - α) × EMA_(t-1)
// where α = 2 / (window + 1) is the smoothing factor
func (ma *MovingAverage) PredictEMA(prices []float64) float64 {
	n := len(prices)
	if n == 0 {
		return 0
	}

	// Smoothing factor: higher α means more weight on recent prices
	alpha := 2.0 / (float64(ma.Window) + 1.0)

	// Initialize EMA with first price
	ema := prices[0]

	// Calculate EMA iteratively
	for i := 1; i < n; i++ {
		ema = alpha*prices[i] + (1-alpha)*ema
	}

	return ema
}

// GetEMASeries returns the full EMA series (useful for error calculation)
func (ma *MovingAverage) GetEMASeries(prices []float64) []float64 {
	n := len(prices)
	if n == 0 {
		return nil
	}

	alpha := 2.0 / (float64(ma.Window) + 1.0)
	emas := make([]float64, n)
	emas[0] = prices[0]

	for i := 1; i < n; i++ {
		emas[i] = alpha*prices[i] + (1-alpha)*emas[i-1]
	}

	return emas
}

// =====================================================
// EXPONENTIAL SMOOTHING (Holt's Linear Trend Method)
// Also called Double Exponential Smoothing.
// Captures both LEVEL and TREND in time series data.
// This is one of the most classic forecasting methods
// used in econometrics and business forecasting.
//
// Reference: Holt, C.E. (1957) "Forecasting Trends and
//            Seasonal by Exponentially Weighted Averages"
// =====================================================

// HoltLinear implements Holt's Linear Trend (Double Exponential Smoothing)
type HoltLinear struct {
	Alpha float64 // Level smoothing (0-1), higher = more reactive to recent values
	Beta  float64 // Trend smoothing (0-1), higher = more reactive to recent trends
}

// NewHoltLinear creates a new Holt's Linear Trend model
// alpha: level smoothing factor (recommended 0.3-0.8)
// beta: trend smoothing factor (recommended 0.1-0.5)
func NewHoltLinear(alpha, beta float64) *HoltLinear {
	if alpha <= 0 || alpha >= 1 {
		alpha = 0.5
	}
	if beta <= 0 || beta >= 1 {
		beta = 0.3
	}
	return &HoltLinear{Alpha: alpha, Beta: beta}
}

// Predict forecasts the next value using Holt's method
// Formulas:
//
//	Level:    L_t = α × Y_t + (1 - α) × (L_(t-1) + T_(t-1))
//	Trend:    T_t = β × (L_t - L_(t-1)) + (1 - β) × T_(t-1)
//	Forecast: F_(t+h) = L_t + h × T_t
func (h *HoltLinear) Predict(prices []float64, stepsAhead int) float64 {
	n := len(prices)
	if n < 2 {
		if n == 1 {
			return prices[0]
		}
		return 0
	}

	// Initialize level and trend
	level := prices[0]
	trend := prices[1] - prices[0]

	// Iterate through all observations
	for i := 1; i < n; i++ {
		prevLevel := level
		level = h.Alpha*prices[i] + (1-h.Alpha)*(prevLevel+trend)
		trend = h.Beta*(level-prevLevel) + (1-h.Beta)*trend
	}

	// Forecast h steps ahead
	return level + float64(stepsAhead)*trend
}

// GetFittedValues returns the in-sample fitted values (for error calculation)
func (h *HoltLinear) GetFittedValues(prices []float64) []float64 {
	n := len(prices)
	if n < 2 {
		return prices
	}

	fitted := make([]float64, n)
	level := prices[0]
	trend := prices[1] - prices[0]
	fitted[0] = prices[0]

	for i := 1; i < n; i++ {
		// One-step-ahead forecast
		fitted[i] = level + trend

		// Update level and trend
		prevLevel := level
		level = h.Alpha*prices[i] + (1-h.Alpha)*(prevLevel+trend)
		trend = h.Beta*(level-prevLevel) + (1-h.Beta)*trend
	}

	return fitted
}

// =====================================================
// K-NEAREST NEIGHBORS (KNN) REGRESSION
// A non-parametric, instance-based learning algorithm.
// For stock prediction: finds the K most similar historical
// price patterns and averages their next-day outcomes.
//
// This is one of the most famous ML algorithms, simple
// yet effective. Uses feature vectors of recent price windows.
//
// Reference: Cover, T. & Hart, P. (1967) "Nearest Neighbor
//            Pattern Classification"
// =====================================================

// KNNRegressor implements K-Nearest Neighbors for regression
type KNNRegressor struct {
	K            int  // Number of neighbors
	WindowSize   int  // Size of the feature window (lookback period)
	WeightByDist bool // Whether to weight predictions by inverse distance
}

// NewKNNRegressor creates a new KNN regressor
func NewKNNRegressor(k, windowSize int, weightByDist bool) *KNNRegressor {
	if k <= 0 {
		k = 5
	}
	if windowSize <= 0 {
		windowSize = 5
	}
	return &KNNRegressor{K: k, WindowSize: windowSize, WeightByDist: weightByDist}
}

type neighbor struct {
	distance float64
	target   float64
}

// Predict finds the K nearest historical pattern matches and predicts the next value
func (knn *KNNRegressor) Predict(prices []float64) float64 {
	n := len(prices)
	if n <= knn.WindowSize {
		if n > 0 {
			return prices[n-1]
		}
		return 0
	}

	// The query pattern is the last WindowSize prices (normalized)
	queryStart := n - knn.WindowSize
	query := normalize(prices[queryStart:])

	// Build dataset of historical windows and their targets (next price)
	var neighbors []neighbor
	for i := 0; i <= n-knn.WindowSize-1; i++ {
		endIdx := i + knn.WindowSize
		if endIdx >= n {
			break
		}

		window := normalize(prices[i:endIdx])
		target := prices[endIdx] // The price that followed this pattern

		dist := euclideanDistance(query, window)
		neighbors = append(neighbors, neighbor{distance: dist, target: target})
	}

	if len(neighbors) == 0 {
		return prices[n-1]
	}

	// Sort by distance (closest first)
	sort.Slice(neighbors, func(i, j int) bool {
		return neighbors[i].distance < neighbors[j].distance
	})

	// Take K nearest
	k := knn.K
	if k > len(neighbors) {
		k = len(neighbors)
	}
	nearest := neighbors[:k]

	if knn.WeightByDist {
		// Inverse-distance weighted average
		return knn.weightedAverage(nearest)
	}

	// Simple average of K nearest targets
	sum := 0.0
	for _, nb := range nearest {
		sum += nb.target
	}
	return sum / float64(k)
}

// weightedAverage computes inverse-distance weighted prediction
// Closer neighbors have more influence on the prediction
func (knn *KNNRegressor) weightedAverage(neighbors []neighbor) float64 {
	totalWeight := 0.0
	weightedSum := 0.0

	for _, nb := range neighbors {
		weight := 1.0 / (nb.distance + 1e-10) // Add small epsilon to avoid division by zero
		totalWeight += weight
		weightedSum += weight * nb.target
	}

	if totalWeight == 0 {
		return 0
	}
	return weightedSum / totalWeight
}

// normalize scales a slice to zero mean and unit variance (z-score normalization)
func normalize(data []float64) []float64 {
	n := len(data)
	if n == 0 {
		return data
	}

	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(n)

	stdDev := 0.0
	for _, v := range data {
		stdDev += (v - mean) * (v - mean)
	}
	stdDev = math.Sqrt(stdDev / float64(n))

	if stdDev == 0 {
		result := make([]float64, n)
		return result
	}

	result := make([]float64, n)
	for i, v := range data {
		result[i] = (v - mean) / stdDev
	}
	return result
}

// euclideanDistance computes the Euclidean distance between two vectors
func euclideanDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.MaxFloat64
	}

	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// =====================================================
// AUTO-REGRESSIVE (AR) MODEL
// A simplified ARIMA-like model (the AR component).
// Uses past values to linearly predict the next value.
// AR(p): Y_t = c + φ1×Y_(t-1) + φ2×Y_(t-2) + ... + φp×Y_(t-p) + ε
//
// ARIMA is the gold standard for time series forecasting
// in statistics and econometrics. This implements the AR
// component with OLS estimation.
//
// Reference: Box, G.E.P. & Jenkins, G.M. (1970)
//            "Time Series Analysis: Forecasting and Control"
// =====================================================

// AutoRegressive implements an AR(p) model
type AutoRegressive struct {
	Order        int       // p — number of lagged terms
	Coefficients []float64 // φ1, φ2, ..., φp
	Intercept    float64   // c (constant term)
}

// NewAutoRegressive creates a new AR model of given order
func NewAutoRegressive(order int) *AutoRegressive {
	if order <= 0 {
		order = 3
	}
	return &AutoRegressive{Order: order}
}

// Train fits the AR model using Ordinary Least Squares (OLS)
// This creates a regression where X = [Y_(t-1), Y_(t-2), ..., Y_(t-p)]
// and target = Y_t
func (ar *AutoRegressive) Train(prices []float64) {
	n := len(prices)
	if n <= ar.Order {
		return
	}

	// Number of training samples
	numSamples := n - ar.Order
	p := ar.Order

	// Build design matrix X and target vector y
	// Each row of X = [1, Y_(t-1), Y_(t-2), ..., Y_(t-p)]
	// y[i] = Y_t

	// We solve this via normal equations: β = (X'X)^(-1) × X'y
	// For simplicity and numerical stability, use iterative approach

	// Use multiple simple linear regression via gradient descent
	// Initialize coefficients
	ar.Coefficients = make([]float64, p)
	ar.Intercept = 0

	// Learning rate and iterations for gradient descent
	lr := 0.0001
	iterations := 1000

	// Normalize prices for stable training
	mean, std := meanStd(prices)
	normalized := make([]float64, n)
	for i, v := range prices {
		if std > 0 {
			normalized[i] = (v - mean) / std
		} else {
			normalized[i] = 0
		}
	}

	for iter := 0; iter < iterations; iter++ {
		// Compute gradients
		gradIntercept := 0.0
		gradCoeffs := make([]float64, p)

		for i := p; i < n; i++ {
			// Predicted value
			pred := ar.Intercept
			for j := 0; j < p; j++ {
				pred += ar.Coefficients[j] * normalized[i-j-1]
			}

			// Error
			err := pred - normalized[i]

			// Accumulate gradients
			gradIntercept += err
			for j := 0; j < p; j++ {
				gradCoeffs[j] += err * normalized[i-j-1]
			}
		}

		// Update parameters
		ar.Intercept -= lr * gradIntercept / float64(numSamples)
		for j := 0; j < p; j++ {
			ar.Coefficients[j] -= lr * gradCoeffs[j] / float64(numSamples)
		}
	}

	// De-normalize: convert coefficients back to original scale
	// Original prediction: pred_original = pred_normalized * std + mean
	// We keep coefficients in normalized space and denormalize at prediction time
}

// Predict forecasts the next value
func (ar *AutoRegressive) Predict(prices []float64) float64 {
	n := len(prices)
	if n < ar.Order || len(ar.Coefficients) == 0 {
		if n > 0 {
			return prices[n-1]
		}
		return 0
	}

	mean, std := meanStd(prices)

	// Normalize the last p prices
	pred := ar.Intercept
	for j := 0; j < ar.Order; j++ {
		normalizedVal := 0.0
		if std > 0 {
			normalizedVal = (prices[n-j-1] - mean) / std
		}
		pred += ar.Coefficients[j] * normalizedVal
	}

	// De-normalize
	if std > 0 {
		return pred*std + mean
	}
	return mean
}

// GetFittedValues returns in-sample predictions for error calculation
func (ar *AutoRegressive) GetFittedValues(prices []float64) []float64 {
	n := len(prices)
	if n <= ar.Order || len(ar.Coefficients) == 0 {
		return prices
	}

	mean, std := meanStd(prices)

	fitted := make([]float64, n)
	// First p values can't be predicted
	for i := 0; i < ar.Order; i++ {
		fitted[i] = prices[i]
	}

	for i := ar.Order; i < n; i++ {
		pred := ar.Intercept
		for j := 0; j < ar.Order; j++ {
			normalizedVal := 0.0
			if std > 0 {
				normalizedVal = (prices[i-j-1] - mean) / std
			}
			pred += ar.Coefficients[j] * normalizedVal
		}
		if std > 0 {
			fitted[i] = pred*std + mean
		} else {
			fitted[i] = mean
		}
	}

	return fitted
}

// meanStd calculates mean and standard deviation of a slice
func meanStd(data []float64) (float64, float64) {
	n := len(data)
	if n == 0 {
		return 0, 0
	}

	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(n)

	variance := 0.0
	for _, v := range data {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(n)

	return mean, math.Sqrt(variance)
}

// =====================================================
// WEIGHTED MOVING AVERAGE (WMA)
// More recent prices are given linearly increasing weights.
// Weight_i = i / (1 + 2 + 3 + ... + n) = i / (n*(n+1)/2)
//
// Used alongside SMA and EMA as a technical indicator.
// =====================================================

// WeightedMovingAverage implements WMA prediction
type WeightedMovingAverage struct {
	Window int
}

// NewWeightedMovingAverage creates a WMA predictor
func NewWeightedMovingAverage(window int) *WeightedMovingAverage {
	if window <= 0 {
		window = 5
	}
	return &WeightedMovingAverage{Window: window}
}

// Predict computes the weighted moving average prediction
func (wma *WeightedMovingAverage) Predict(prices []float64) float64 {
	n := len(prices)
	if n == 0 {
		return 0
	}

	window := wma.Window
	if window > n {
		window = n
	}

	// Denominator: 1 + 2 + 3 + ... + window = window*(window+1)/2
	denom := float64(window * (window + 1) / 2)

	sum := 0.0
	for i := 0; i < window; i++ {
		weight := float64(i + 1) // 1, 2, 3, ..., window
		sum += weight * prices[n-window+i]
	}

	return sum / denom
}

// =====================================================
// ENSEMBLE MODEL
// Combines all algorithms using weighted averaging.
// Each model's weight is inversely proportional to its
// error (lower error → higher weight).
//
// Ensemble methods are one of the most powerful ideas in
// ML — they reduce variance and improve generalization.
//
// Reference: Dietterich, T.G. (2000) "Ensemble Methods
//            in Machine Learning"
// =====================================================

// EnsemblePrediction holds a single model's contribution
type EnsemblePrediction struct {
	Algorithm string  `json:"algorithm"`
	Predicted float64 `json:"predicted_price"`
	RMSE      float64 `json:"rmse"`
	MAE       float64 `json:"mae"`
	Weight    float64 `json:"weight"`
}

// EnsemblePredict runs all models and returns the weighted ensemble prediction
func EnsemblePredict(prices []float64) (float64, []EnsemblePrediction) {
	n := len(prices)
	if n < 5 {
		if n > 0 {
			return prices[n-1], nil
		}
		return 0, nil
	}

	var results []EnsemblePrediction

	// 1. Linear Regression
	lr := NewLinearRegression()
	xLR := make([]float64, n)
	for i := range xLR {
		xLR[i] = float64(i)
	}
	lr.Train(xLR, prices)
	lrPred := lr.Predict(float64(n))
	lrFitted := make([]float64, n)
	for i := range xLR {
		lrFitted[i] = lr.Predict(xLR[i])
	}
	results = append(results, EnsemblePrediction{
		Algorithm: "Linear Regression",
		Predicted: lrPred,
		RMSE:      CalculateRMSE(prices, lrFitted),
		MAE:       CalculateMAE(prices, lrFitted),
	})

	// 2. EMA
	ema := NewMovingAverage(5)
	emaPred := ema.PredictEMA(prices)
	emaFitted := ema.GetEMASeries(prices)
	results = append(results, EnsemblePrediction{
		Algorithm: "Exponential Moving Average (EMA-5)",
		Predicted: emaPred,
		RMSE:      CalculateRMSE(prices, emaFitted),
		MAE:       CalculateMAE(prices, emaFitted),
	})

	// 3. Holt's Linear Trend
	holt := NewHoltLinear(0.5, 0.3)
	holtPred := holt.Predict(prices, 1)
	holtFitted := holt.GetFittedValues(prices)
	results = append(results, EnsemblePrediction{
		Algorithm: "Holt's Linear Trend (Double Exponential Smoothing)",
		Predicted: holtPred,
		RMSE:      CalculateRMSE(prices, holtFitted),
		MAE:       CalculateMAE(prices, holtFitted),
	})

	// 4. KNN Regression
	knn := NewKNNRegressor(5, 5, true)
	knnPred := knn.Predict(prices)
	results = append(results, EnsemblePrediction{
		Algorithm: "K-Nearest Neighbors (KNN-5, window=5)",
		Predicted: knnPred,
		RMSE:      0, // KNN doesn't produce fitted values easily
		MAE:       0,
	})

	// 5. Auto-Regressive AR(3)
	ar := NewAutoRegressive(3)
	ar.Train(prices)
	arPred := ar.Predict(prices)
	arFitted := ar.GetFittedValues(prices)
	arRMSE := CalculateRMSE(prices[3:], arFitted[3:]) // Skip first p values
	arMAE := CalculateMAE(prices[3:], arFitted[3:])
	results = append(results, EnsemblePrediction{
		Algorithm: "Auto-Regressive AR(3)",
		Predicted: arPred,
		RMSE:      arRMSE,
		MAE:       arMAE,
	})

	// 6. Weighted Moving Average
	wma := NewWeightedMovingAverage(5)
	wmaPred := wma.Predict(prices)
	results = append(results, EnsemblePrediction{
		Algorithm: "Weighted Moving Average (WMA-5)",
		Predicted: wmaPred,
		RMSE:      0,
		MAE:       0,
	})

	// Calculate ensemble weights (inverse RMSE)
	totalInvRMSE := 0.0
	for i := range results {
		if results[i].RMSE > 0 {
			results[i].Weight = 1.0 / results[i].RMSE
		} else {
			results[i].Weight = 1.0 // Default weight for models without RMSE
		}
		totalInvRMSE += results[i].Weight
	}

	// Normalize weights to sum to 1
	if totalInvRMSE > 0 {
		for i := range results {
			results[i].Weight /= totalInvRMSE
		}
	}

	// Weighted ensemble prediction
	ensemblePred := 0.0
	for _, r := range results {
		ensemblePred += r.Weight * r.Predicted
	}

	// Round weights for display
	for i := range results {
		results[i].Weight = math.Round(results[i].Weight*10000) / 10000
		results[i].Predicted = math.Round(results[i].Predicted*100) / 100
		results[i].RMSE = math.Round(results[i].RMSE*10000) / 10000
		results[i].MAE = math.Round(results[i].MAE*10000) / 10000
	}

	return math.Round(ensemblePred*100) / 100, results
}

// =====================================================
// LINEAR REGRESSION
// Classic statistical method — fits a straight line (y = mx + b)
// using the Ordinary Least Squares (OLS) method.
// =====================================================

// LinearRegression implements simple linear regression
type LinearRegression struct {
	Slope     float64
	Intercept float64
}

// NewLinearRegression creates a new LinearRegression model
func NewLinearRegression() *LinearRegression {
	return &LinearRegression{}
}

// Train fits the linear regression model to training data
// using Ordinary Least Squares: slope = (n·∑xy − ∑x·∑y) / (n·∑x² − (∑x)²)
func (lr *LinearRegression) Train(x []float64, y []float64) {
	n := float64(len(x))
	if n == 0 || n != float64(len(y)) {
		return
	}

	var sumX, sumY, sumXY, sumXX float64

	for i := 0; i < int(n); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumXX += x[i] * x[i]
	}

	lr.Slope = (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
	lr.Intercept = (sumY - lr.Slope*sumX) / n
}

// Predict returns the predicted y value for a given x
func (lr *LinearRegression) Predict(x float64) float64 {
	return lr.Slope*x + lr.Intercept
}

// =====================================================
// ERROR METRICS: RMSE & MAE
// =====================================================

// CalculateRMSE computes Root Mean Squared Error between actual and predicted values
func CalculateRMSE(actual []float64, predicted []float64) float64 {
	n := len(actual)
	if n == 0 || n != len(predicted) {
		return 0.0
	}

	var sumSquaredError float64
	for i := 0; i < n; i++ {
		err := actual[i] - predicted[i]
		sumSquaredError += err * err
	}

	return math.Sqrt(sumSquaredError / float64(n))
}

// CalculateMAE computes Mean Absolute Error between actual and predicted values
func CalculateMAE(actual []float64, predicted []float64) float64 {
	n := len(actual)
	if n == 0 || n != len(predicted) {
		return 0.0
	}

	var sumAbsError float64
	for i := 0; i < n; i++ {
		sumAbsError += math.Abs(actual[i] - predicted[i])
	}

	return sumAbsError / float64(n)
}
