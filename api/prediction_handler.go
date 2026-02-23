package api

import (
	"net/http"

	"github.com/919Umesh/stock_market_sim/internal/ml"
	"github.com/919Umesh/stock_market_sim/pkg/apperr"
	"github.com/gin-gonic/gin"
)

type PredictionHandler struct {
	service ml.Service
}

func NewPredictionHandler(service ml.Service) *PredictionHandler {
	return &PredictionHandler{service: service}
}

// GetPrediction predicts using default ensemble or a specific algorithm via ?algorithm= query param
// Supported algorithms: linear_regression, ema, holt, knn, ar, wma, ensemble (default)
func (h *PredictionHandler) GetPrediction(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Symbol is required")
		return
	}

	algorithm := c.DefaultQuery("algorithm", "ensemble")

	result, err := h.service.PredictWithAlgorithm(symbol, algorithm)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to predict price: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prediction": result,
	})
}

// CompareAlgorithms runs all ML algorithms on a stock and returns comparison
func (h *PredictionHandler) CompareAlgorithms(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Symbol is required")
		return
	}

	result, err := h.service.CompareAlgorithms(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to compare algorithms: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comparison": result,
	})
}

// ListAlgorithms returns all available prediction algorithms
func (h *PredictionHandler) ListAlgorithms(c *gin.Context) {
	algorithms := []gin.H{
		{
			"key":         "linear_regression",
			"name":        "Linear Regression",
			"description": "Simple linear regression fits a straight line to historical prices using least squares. Classic statistical method.",
			"type":        "Statistical",
		},
		{
			"key":         "ema",
			"name":        "Exponential Moving Average (EMA)",
			"description": "Applies exponentially decreasing weights to older observations. Widely used in technical analysis for trend identification.",
			"type":        "Technical Analysis",
		},
		{
			"key":         "holt",
			"name":        "Holt's Linear Trend (Double Exponential Smoothing)",
			"description": "Extension of exponential smoothing that captures both level and trend in time series data. A classic forecasting method.",
			"type":        "Time Series",
		},
		{
			"key":         "knn",
			"name":        "K-Nearest Neighbors (KNN) Regression",
			"description": "Non-parametric method that predicts based on the K most similar historical price patterns. Classic machine learning algorithm.",
			"type":        "Machine Learning",
		},
		{
			"key":         "ar",
			"name":        "Auto-Regressive AR(3) Model",
			"description": "Models the next value as a linear combination of previous values. Foundation of ARIMA, the gold standard for time series forecasting.",
			"type":        "Time Series",
		},
		{
			"key":         "wma",
			"name":        "Weighted Moving Average (WMA)",
			"description": "Moving average with linearly increasing weights, giving more importance to recent prices. Common technical indicator.",
			"type":        "Technical Analysis",
		},
		{
			"key":         "ensemble",
			"name":        "Ensemble (Weighted Average)",
			"description": "Combines predictions from all models using inverse-RMSE weighted averaging. Reduces variance and often outperforms individual models.",
			"type":        "Ensemble",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"algorithms": algorithms,
		"total":      len(algorithms),
	})
}
