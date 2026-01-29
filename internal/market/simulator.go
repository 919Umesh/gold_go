package market

import (
	"context"
	"log/slog"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/919Umesh/gold_go/internal/stock"
	"github.com/919Umesh/gold_go/models"
	"gorm.io/gorm"
)

// Simulator handles the market price simulation
type Simulator struct {
	db        *gorm.DB
	stockRepo stock.Repository
	mu        sync.RWMutex
	isRunning bool
	stopChan  chan struct{}
}

// PriceSimulationParams holds parameters for price simulation
type PriceSimulationParams struct {
	Drift      float64 // Expected return (mu)
	Volatility float64 // Standard deviation (sigma)
	TimeStep   float64 // Time increment (dt)
}

func NewSimulator(db *gorm.DB, stockRepo stock.Repository) *Simulator {
	return &Simulator{
		db:        db,
		stockRepo: stockRepo,
		stopChan:  make(chan struct{}),
	}
}

// Start begins the market simulation
func (s *Simulator) Start(ctx context.Context) {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return
	}
	s.isRunning = true
	s.mu.Unlock()

	slog.Info("Market simulator started")

	// Update prices every 5 minutes during market hours
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.IsMarketOpen() {
				s.UpdateAllPrices()
			}
		case <-ctx.Done():
			slog.Info("Stopping market simulator")
			s.mu.Lock()
			s.isRunning = false
			s.mu.Unlock()
			return
		case <-s.stopChan:
			slog.Info("Market simulator stopped")
			s.mu.Lock()
			s.isRunning = false
			s.mu.Unlock()
			return
		}
	}
}

// Stop stops the market simulation
func (s *Simulator) Stop() {
	close(s.stopChan)
}

// IsMarketOpen checks if the Nepal Stock Exchange is currently open
// NEPSE hours: Sunday-Thursday, 11:00 AM - 3:00 PM Nepal Time
func (s *Simulator) IsMarketOpen() bool {
	now := time.Now()

	// Check if it's Friday or Saturday (market closed)
	weekday := now.Weekday()
	if weekday == time.Friday || weekday == time.Saturday {
		return false
	}

	// Check market hours (11:00 AM - 3:00 PM)
	hour := now.Hour()
	if hour < 11 || hour >= 15 {
		return false
	}

	return true
}

// UpdateAllPrices updates prices for all active companies
func (s *Simulator) UpdateAllPrices() {
	companies, err := s.stockRepo.ListCompanies(1000, 0)
	if err != nil {
		slog.Error("Failed to fetch companies for price update", "error", err)
		return
	}

	for _, company := range companies {
		if err := s.UpdateCompanyPrice(company.ID); err != nil {
			slog.Error("Failed to update price", "company", company.Symbol, "error", err)
		}
	}

	slog.Info("Updated prices for all companies", "count", len(companies))
}

// UpdateCompanyPrice updates the price for a single company
func (s *Simulator) UpdateCompanyPrice(companyID uint) error {
	// Get the latest price
	latestPrice, err := s.stockRepo.GetLatestPrice(companyID)
	if err != nil {
		// If no price exists, skip (will be seeded separately)
		return nil
	}

	// Simulate new price using Geometric Brownian Motion
	params := s.getSimulationParams(companyID)
	newPrice := s.simulateNextPrice(latestPrice.ClosePrice, params)

	// Create OHLCV data (simplified - using same price for all)
	// In a more sophisticated version, we'd simulate intraday movements
	stockPrice := &models.StockPrice{
		CompanyID:  companyID,
		OpenPrice:  latestPrice.ClosePrice,
		HighPrice:  math.Max(latestPrice.ClosePrice, newPrice),
		LowPrice:   math.Min(latestPrice.ClosePrice, newPrice),
		ClosePrice: newPrice,
		Volume:     s.simulateVolume(),
		Timestamp:  time.Now(),
		Timeframe:  "5m",
	}

	return s.stockRepo.CreateStockPrice(stockPrice)
}

// simulateNextPrice calculates the next price using Geometric Brownian Motion
// Formula: S(t+dt) = S(t) * exp((mu - 0.5*sigma^2)*dt + sigma*sqrt(dt)*Z)
// Where Z is a random variable from standard normal distribution
func (s *Simulator) simulateNextPrice(currentPrice float64, params PriceSimulationParams) float64 {
	// Generate random number from standard normal distribution
	z := rand.NormFloat64()

	// Calculate drift component
	drift := (params.Drift - 0.5*params.Volatility*params.Volatility) * params.TimeStep

	// Calculate diffusion component
	diffusion := params.Volatility * math.Sqrt(params.TimeStep) * z

	// Calculate new price
	newPrice := currentPrice * math.Exp(drift+diffusion)

	// Add some bounds to prevent extreme movements
	maxChange := 0.10 // 10% max change per update
	if newPrice > currentPrice*(1+maxChange) {
		newPrice = currentPrice * (1 + maxChange)
	} else if newPrice < currentPrice*(1-maxChange) {
		newPrice = currentPrice * (1 - maxChange)
	}

	// Round to 2 decimal places
	return math.Round(newPrice*100) / 100
}

// getSimulationParams returns simulation parameters for a company
// In a real system, these could be customized per company/sector
func (s *Simulator) getSimulationParams(companyID uint) PriceSimulationParams {
	// Default parameters (can be customized based on sector, market cap, etc.)
	return PriceSimulationParams{
		Drift:      0.0001,           // 0.01% expected return per 5-minute interval
		Volatility: 0.02,             // 2% volatility
		TimeStep:   1.0 / (252 * 78), // 5 minutes in trading year (252 days, 78 5-min intervals per day)
	}
}

// simulateVolume generates a random trading volume
func (s *Simulator) simulateVolume() int64 {
	// Generate volume between 10,000 and 1,000,000
	baseVolume := 10000
	randomFactor := rand.Intn(990000)
	return int64(baseVolume + randomFactor)
}

// GenerateMarketEvent creates a random market event for a company
func (s *Simulator) GenerateMarketEvent(companyID uint) error {
	eventTypes := []models.MarketEventType{
		models.MarketEventEarnings,
		models.MarketEventNews,
		models.MarketEventDividend,
	}

	eventType := eventTypes[rand.Intn(len(eventTypes))]

	event := &models.MarketEvent{
		CompanyID:        companyID,
		EventType:        eventType,
		Title:            s.generateEventTitle(eventType),
		Description:      "Simulated market event",
		ImpactPercentage: (rand.Float64() - 0.5) * 10, // -5% to +5%
		EventDate:        time.Now().Add(time.Duration(rand.Intn(30)) * 24 * time.Hour),
	}

	return s.stockRepo.CreateMarketEvent(event)
}

func (s *Simulator) generateEventTitle(eventType models.MarketEventType) string {
	titles := map[models.MarketEventType][]string{
		models.MarketEventEarnings: {
			"Q1 Earnings Report Released",
			"Strong Quarterly Performance",
			"Earnings Beat Expectations",
		},
		models.MarketEventNews: {
			"Major Partnership Announced",
			"New Product Launch",
			"Expansion Plans Revealed",
		},
		models.MarketEventDividend: {
			"Dividend Declared",
			"Increased Dividend Payout",
			"Special Dividend Announced",
		},
	}

	options := titles[eventType]
	return options[rand.Intn(len(options))]
}
