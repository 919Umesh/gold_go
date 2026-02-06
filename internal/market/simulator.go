package market

import (
	"log"
	"math/rand"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/models"
)

type Simulator struct {
	stockRepo stock.Repository
	ticker    *time.Ticker
	quit      chan struct{}
}

func NewSimulator(stockRepo stock.Repository) *Simulator {
	return &Simulator{
		stockRepo: stockRepo,
		quit:      make(chan struct{}),
	}
}

func (s *Simulator) Start(interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.UpdateMarket()
			case <-s.quit:
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *Simulator) Stop() {
	close(s.quit)
}

func (s *Simulator) UpdateMarket() {
	// 1. Fetch all active companies
	companies, err := s.stockRepo.ListCompanies(100, 0)
	if err != nil {
		log.Printf("Error fetching companies for simulation: %v", err)
		return
	}

	for _, company := range companies {
		if err := s.SimulateStockPrice(&company); err != nil {
			log.Printf("Error simulating price for %s: %v", company.Symbol, err)
		}
	}
}

func (s *Simulator) SimulateStockPrice(company *models.Company) error {
	// Get latest price
	latestPrice, err := s.stockRepo.GetLatestPrice(company.ID)

	var currentPrice float64
	if err != nil {
		// No price yet? Start with some base
		currentPrice = 100.0 // Default start price
	} else {
		currentPrice = latestPrice.ClosePrice
	}

	// Simple random walk simulation
	// Volatility around 2%
	changePct := (rand.Float64() - 0.5) * 0.04
	newPrice := currentPrice * (1 + changePct)

	if newPrice < 0.01 {
		newPrice = 0.01
	}

	// Create OHLCV bar for this interval (simplified: assuming 1 tick per interval)
	// In reality we might Aggregate ticks. Here just saving "Close" as new candle.

	now := time.Now()

	newStockPrice := &models.StockPrice{
		CompanyID:  company.ID,
		OpenPrice:  currentPrice,
		HighPrice:  max(currentPrice, newPrice),
		LowPrice:   min(currentPrice, newPrice),
		ClosePrice: newPrice,
		Volume:     int64(rand.Intn(1000)), // Random volume
		Timestamp:  now,
		Timeframe:  "1m", // simulation interval assumed
	}

	return s.stockRepo.CreateStockPrice(newStockPrice)
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
