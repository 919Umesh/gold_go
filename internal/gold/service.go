package gold

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/919Umesh/gold_go/config"
	"github.com/919Umesh/gold_go/models"
	"github.com/919Umesh/gold_go/pkg/queue"
	"gorm.io/gorm"
)

type PriceFetcher interface {
	FetchPrice(ctx context.Context) (float64, error)
}

type Service struct {
	db         *gorm.DB
	cfg        *config.Config
	priceCache *PriceCache
	fetcher    PriceFetcher
	workerPool *queue.WorkerPool
	updateChan chan struct{} // Channel to trigger immediate update
}

type PriceCache struct {
	price float64
	mu    sync.RWMutex
	time  time.Time
}

// PriceHistoryJob implements queue.Job to save price history asynchronously
type PriceHistoryJob struct {
	DB        *gorm.DB
	GoldPrice *models.GoldPrice
}

func (j *PriceHistoryJob) Process() error {
	// Using raw SQL for insertion
	query := `INSERT INTO gold_prices (price_per_gram, source, updated_at) VALUES (?, ?, ?)`
	if err := j.DB.Exec(query, j.GoldPrice.PricePerGram, j.GoldPrice.Source, j.GoldPrice.UpdatedAt).Error; err != nil {
		return fmt.Errorf("failed to save gold price: %w", err)
	}
	slog.Info("Async job: Gold price history saved", "price", j.GoldPrice.PricePerGram)
	return nil
}

func NewService(db *gorm.DB, cfg *config.Config, wp *queue.WorkerPool) *Service {
	service := &Service{
		db:         db,
		cfg:        cfg,
		priceCache: &PriceCache{},
		workerPool: wp,
		updateChan: make(chan struct{}, 1),
	}

	// Use real fetcher if URL is provided, otherwise mock
	if cfg.GoldProvider != "" {
		service.fetcher = &RealPriceFetcher{
			client: &http.Client{Timeout: 10 * time.Second},
			url:    cfg.GoldProvider,
		}
	} else {
		service.fetcher = &MockPriceFetcher{}
	}

	return service
}

func (s *Service) StartPriceUpdater(ctx context.Context) {
	ticker := time.NewTicker(600 * time.Second) // Update every 10 minutes
	defer ticker.Stop()

	// Initial fetch
	s.updatePrice(ctx)

	for {
		select {
		case <-ticker.C:
			s.updatePrice(ctx)
		case <-s.updateChan:
			s.updatePrice(ctx)
		case <-ctx.Done():
			slog.Info("Stopping price updater")
			return
		}
	}
}

func (s *Service) TriggerUpdate() {
	select {
	case s.updateChan <- struct{}{}:
	default:
		// Channel full, update already pending
	}
}

func (s *Service) updatePrice(ctx context.Context) {
	price, err := s.fetcher.FetchPrice(ctx)
	if err != nil {
		slog.Error("Failed to fetch gold price", "error", err)
		return
	}

	s.priceCache.mu.Lock()
	s.priceCache.price = price
	s.priceCache.time = time.Now()
	s.priceCache.mu.Unlock()

	goldPrice := &models.GoldPrice{
		PricePerGram: price,
		Source:       "provider",
		UpdatedAt:    time.Now(),
	}

	// Submit job to worker pool for async DB insertion
	s.workerPool.Submit(&PriceHistoryJob{
		DB:        s.db,
		GoldPrice: goldPrice,
	})

	slog.Info("Gold price updated in memory", "price", price)
}

func (s *Service) GetCurrentPrice() (float64, time.Time, error) {
	s.priceCache.mu.RLock()
	defer s.priceCache.mu.RUnlock()

	if s.priceCache.price == 0 {
		return 0, time.Time{}, fmt.Errorf("price not available")
	}

	return s.priceCache.price, s.priceCache.time, nil
}

func (s *Service) GetPriceHistory(days int) ([]models.GoldPrice, error) {
	var prices []models.GoldPrice
	since := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT * FROM gold_prices 
		WHERE updated_at >= ? 
		ORDER BY updated_at desc
	`
	err := s.db.Raw(query, since).Scan(&prices).Error

	return prices, err
}

type MockPriceFetcher struct{}

func (m *MockPriceFetcher) FetchPrice(ctx context.Context) (float64, error) {
	basePrice := 6500.0
	// Simple simulation of price fluctuation
	variation := (float64(time.Now().Unix()%100) - 50) / 100.0
	return basePrice + (basePrice * variation), nil
}

type RealPriceFetcher struct {
	client *http.Client
	url    string
}

func (r *RealPriceFetcher) FetchPrice(ctx context.Context) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Price float64 `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Price, nil
}
