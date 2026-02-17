package models

import (
	"time"
)

type StockPrediction struct {
	ID              string    `json:"id,omitempty"`
	CompanyID       string    `json:"company_id"`
	PredictedPrice  float64   `json:"predicted_price"`
	ConfidenceScore float64   `json:"confidence_score"`
	PredictionDate  time.Time `json:"prediction_date"`
	TargetDate      time.Time `json:"target_date"`
	ModelUsed       string    `json:"model_used"`
	CreatedAt       time.Time `json:"created_at,omitempty"`

	Company *Company `json:"company,omitempty"`
}
