package ml

import (
	"math"
)


type LinearRegression struct {
	Slope     float64
	Intercept float64
}


func NewLinearRegression() *LinearRegression {
	return &LinearRegression{}
}


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


func (lr *LinearRegression) Predict(x float64) float64 {
	return lr.Slope*x + lr.Intercept
}


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
