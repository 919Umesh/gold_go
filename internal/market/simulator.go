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
	latestPrice, err := s.stockRepo.GetLatestPrice(company.ID)

	var currentPrice float64
	if err != nil {
		currentPrice = 100.0 
	} else {
		currentPrice = latestPrice.ClosePrice
	}


	changePct := (rand.Float64() - 0.5) * 0.04
	newPrice := currentPrice * (1 + changePct)

	if newPrice < 0.01 {
		newPrice = 0.01
	}


	now := time.Now()

	newStockPrice := &models.StockPrice{
		CompanyID:  company.ID,
		OpenPrice:  currentPrice,
		HighPrice:  max(currentPrice, newPrice),
		LowPrice:   min(currentPrice, newPrice),
		ClosePrice: newPrice,
		Volume:     int64(rand.Intn(1000)), 
		Timestamp:  now,
		Timeframe:  "1m", 
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
