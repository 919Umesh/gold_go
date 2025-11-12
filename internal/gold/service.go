package gold

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/919Umesh/gold_go/config"
	"github.com/919Umesh/gold_go/models"
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
}

type PriceCache struct {
	price float64
	mu    sync.RWMutex
	time  time.Time
}

func NewService(db *gorm.DB, cfg *config.Config) *Service {
	service := &Service{
		db:         db,
		cfg:        cfg,
		priceCache: &PriceCache{},
	}

	service.fetcher = &MockPriceFetcher{}

	return service
}

func (s *Service) StartPriceUpdater(ctx context.Context) {
	ticker := time.NewTicker(600 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			price, err := s.fetcher.FetchPrice(ctx)
			if err != nil {
				log.Printf("Failed to fetch gold price: %v", err)
				continue
			}

			s.priceCache.mu.Lock()
			s.priceCache.price = price
			s.priceCache.time = time.Now()
			s.priceCache.mu.Unlock()

			goldPrice := &models.GoldPrice{
				PricePerGram: price,
				Source:       "provider",
			}
			if err := s.db.Create(goldPrice).Error; err != nil {
				log.Printf("Failed to save gold price: %v", err)
			}

			log.Printf("Gold price updated: %.2f", price)

		case <-ctx.Done():
			return
		}
	}
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

	err := s.db.Where("updated_at >= ?", since).
		Order("updated_at desc").
		Find(&prices).Error

	return prices, err
}

type MockPriceFetcher struct{}

func (m *MockPriceFetcher) FetchPrice(ctx context.Context) (float64, error) {
	basePrice := 6500.0
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
