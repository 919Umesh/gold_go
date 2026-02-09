package models

import (
	"time"
)

// StockPrediction represents ML-based price predictions
type StockPrediction struct {
	ID              string    `json:"$id,omitempty"`
	CompanyID       string    `json:"company_id"`
	PredictedPrice  float64   `json:"predicted_price"`
	ConfidenceScore float64   `json:"confidence_score"` // 0-100
	PredictionDate  time.Time `json:"prediction_date"`
	TargetDate      time.Time `json:"target_date"`
	ModelUsed       string    `json:"model_used"` // "LSTM", "ARIMA", etc.
	CreatedAt       time.Time `json:"$createdAt,omitempty"`

	// Relationships
	Company *Company `json:"company,omitempty"`
}
