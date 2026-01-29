package models

import (
	"time"

	"gorm.io/gorm"
)

// StockPrediction represents ML-based price predictions
type StockPrediction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CompanyID       uint      `gorm:"index;not null" json:"company_id"`
	PredictedPrice  float64   `gorm:"type:numeric(10,2);not null" json:"predicted_price"`
	ConfidenceScore float64   `gorm:"type:numeric(5,2)" json:"confidence_score"` // 0-100
	PredictionDate  time.Time `gorm:"not null" json:"prediction_date"`
	TargetDate      time.Time `gorm:"not null" json:"target_date"`
	ModelUsed       string    `gorm:"size:50" json:"model_used"` // "LSTM", "ARIMA", etc.
	CreatedAt       time.Time `json:"created_at"`

	// Relationships
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (s *StockPrediction) BeforeCreate(tx *gorm.DB) error {
	s.CreatedAt = time.Now()
	return nil
}
