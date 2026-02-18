package market

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
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
		slog.Error("Error fetching companies for simulation", "error", err)
		return
	}

	for _, company := range companies {
		if err := s.SimulateStockPrice(&company); err != nil {
			slog.Error("Error simulating price", "symbol", company.Symbol, "error", err)
		}
	}
}

func (s *Simulator) SimulateStockPrice(company *models.Company) error {
	latestPrice, err := s.stockRepo.GetLatestPrice(company.ID)

	var currentPrice decimal.Decimal
	if err != nil {
		currentPrice = decimal.NewFromFloat(100.0)
	} else {
		currentPrice = latestPrice.ClosePrice
	}


	changePct := decimal.NewFromFloat((rand.Float64() - 0.5) * 0.04)
	newPrice := currentPrice.Mul(decimal.NewFromInt(1).Add(changePct))

	if newPrice.LessThan(decimal.NewFromFloat(0.01)) {
		newPrice = decimal.NewFromFloat(0.01)
	}

	now := time.Now()

	newStockPrice := &models.StockPrice{
		CompanyID:  company.ID,
		OpenPrice:  currentPrice,
		HighPrice:  decimal.Max(currentPrice, newPrice),
		LowPrice:   decimal.Min(currentPrice, newPrice),
		ClosePrice: newPrice,
		Volume:     int64(rand.Intn(1000)),
		Timestamp:  now,
		Timeframe:  "1m",
	}

	return s.stockRepo.CreateStockPrice(newStockPrice)
}
