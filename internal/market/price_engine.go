package market

import (
	"log/slog"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

const (
	// SensitivityFactor controls price impact magnitude.
	// For small-volume markets: higher = more sensitive to each trade.
	SensitivityFactor = 2.0

	// MaxDailyChangePercent is the circuit breaker limit (±10%)
	MaxDailyChangePercent = 10.0

	// MinPrice floor
	MinPrice = 0.01
)

// TradeImpact is the result of processing a trade
type TradeImpact struct {
	NewPrice       decimal.Decimal `json:"new_price"`
	PriceChange    decimal.Decimal `json:"price_change"`
	PriceChangePct decimal.Decimal `json:"price_change_pct"`
}

// PriceEngine handles order-driven price discovery with high sensitivity
type PriceEngine struct {
	stockRepo stock.Repository
	eventHub  *EventHub

	mu             sync.RWMutex
	dayOpenPrices  map[string]decimal.Decimal
	previousCloses map[string]decimal.Decimal
	dailyVolumes   map[string]int64
	dailyHighs     map[string]decimal.Decimal
	dailyLows      map[string]decimal.Decimal
	lastTradeDay   time.Time

	// Trigger callback (set by trigger worker)
	onPriceUpdate func(companyID string, newPrice decimal.Decimal)
}

func NewPriceEngine(stockRepo stock.Repository, eventHub *EventHub) *PriceEngine {
	pe := &PriceEngine{
		stockRepo:      stockRepo,
		eventHub:       eventHub,
		dayOpenPrices:  make(map[string]decimal.Decimal),
		previousCloses: make(map[string]decimal.Decimal),
		dailyVolumes:   make(map[string]int64),
		dailyHighs:     make(map[string]decimal.Decimal),
		dailyLows:      make(map[string]decimal.Decimal),
		lastTradeDay:   time.Now().Truncate(24 * time.Hour),
	}
	pe.initializeDayData()
	return pe
}

// SetOnPriceUpdate registers a callback called after each price update
func (pe *PriceEngine) SetOnPriceUpdate(fn func(companyID string, newPrice decimal.Decimal)) {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	pe.onPriceUpdate = fn
}

func (pe *PriceEngine) initializeDayData() {
	companies, err := pe.stockRepo.ListCompanies(200, 0)
	if err != nil {
		slog.Error("PriceEngine: failed to load companies", "error", err)
		return
	}

	for _, c := range companies {
		pe.previousCloses[c.ID] = c.CurrentPrice
		pe.dayOpenPrices[c.ID] = c.CurrentPrice
		pe.dailyHighs[c.ID] = c.CurrentPrice
		pe.dailyLows[c.ID] = c.CurrentPrice
	}

	slog.Info("PriceEngine initialized", "companies_loaded", len(pe.previousCloses))
}

func (pe *PriceEngine) checkDayReset() {
	today := time.Now().Truncate(24 * time.Hour)
	if today.After(pe.lastTradeDay) {
		for companyID := range pe.dayOpenPrices {
			if p, err := pe.stockRepo.GetLatestPrice(companyID); err == nil {
				pe.previousCloses[companyID] = p.ClosePrice
			}
		}
		pe.dayOpenPrices = make(map[string]decimal.Decimal)
		pe.dailyVolumes = make(map[string]int64)
		pe.dailyHighs = make(map[string]decimal.Decimal)
		pe.dailyLows = make(map[string]decimal.Decimal)
		pe.lastTradeDay = today
	}
}

// ProcessMatchedTrade is called after each order match to update the price
// Uses high-sensitivity formula: impact = (tradeVolume / totalSupply) * SensitivityFactor
func (pe *PriceEngine) ProcessMatchedTrade(companyID string, tradePrice decimal.Decimal, tradeQty int64) {
	pe.mu.Lock()
	defer pe.mu.Unlock()

	pe.checkDayReset()

	company, err := pe.stockRepo.GetCompanyByID(companyID)
	if err != nil {
		slog.Error("PriceEngine: company not found", "company_id", companyID, "error", err)
		return
	}

	currentPrice := company.CurrentPrice
	if currentPrice.IsZero() {
		currentPrice = tradePrice
	}

	// High-sensitivity price impact formula
	// impact = (tradeVolume / totalSupply) * SensitivityFactor
	totalSupply := company.TotalSupply
	if totalSupply <= 0 {
		totalSupply = 10000
	}

	volumeRatio := float64(tradeQty) / float64(totalSupply)
	impact := volumeRatio * SensitivityFactor

	// Determine direction: if trade price > current, push up; if below, push down
	direction := 1.0
	if tradePrice.LessThan(currentPrice) {
		direction = -1.0
	}

	// Add small noise for realism
	noise := (rand.Float64() - 0.5) * 0.002
	impact = impact*direction + noise

	// Calculate new price
	impactDecimal := decimal.NewFromFloat(impact)
	newPrice := currentPrice.Mul(decimal.NewFromInt(1).Add(impactDecimal))

	// Enforce minimum price
	minP := decimal.NewFromFloat(MinPrice)
	if newPrice.LessThan(minP) {
		newPrice = minP
	}

	// Apply circuit breaker
	newPrice = pe.applyCircuitBreaker(companyID, newPrice)

	// Round to 2 decimal places
	newPrice = newPrice.Round(2)

	// Track intraday data
	pe.dailyVolumes[companyID] += tradeQty
	if pe.dayOpenPrices[companyID].IsZero() {
		pe.dayOpenPrices[companyID] = currentPrice
	}
	if pe.dailyHighs[companyID].IsZero() || newPrice.GreaterThan(pe.dailyHighs[companyID]) {
		pe.dailyHighs[companyID] = newPrice
	}
	if pe.dailyLows[companyID].IsZero() || newPrice.LessThan(pe.dailyLows[companyID]) {
		pe.dailyLows[companyID] = newPrice
	}

	now := time.Now()

	// Create 1-minute candle
	stockPrice := &models.StockPrice{
		CompanyID:  companyID,
		OpenPrice:  currentPrice,
		HighPrice:  decimal.Max(currentPrice, newPrice),
		LowPrice:   decimal.Min(currentPrice, newPrice),
		ClosePrice: newPrice,
		Volume:     tradeQty,
		Timestamp:  now,
		Timeframe:  "1m",
	}
	if err := pe.stockRepo.CreateStockPrice(stockPrice); err != nil {
		slog.Error("PriceEngine: failed to create stock price", "error", err)
	}

	// Update daily candle
	pe.updateDailyCandle(companyID, newPrice, now)

	// Update company price and market cap
	newMarketCap := newPrice.Mul(decimal.NewFromInt(company.TotalSupply))
	if err := pe.stockRepo.UpdateCompanyPrice(companyID, newPrice.String(), newMarketCap.String()); err != nil {
		slog.Warn("PriceEngine: failed to update company price", "error", err)
	}

	// Calculate change metrics
	priceChange := newPrice.Sub(currentPrice)
	priceChangePct := decimal.Zero
	if !currentPrice.IsZero() {
		priceChangePct = priceChange.Div(currentPrice).Mul(decimal.NewFromInt(100))
	}

	// Broadcast SSE events
	if pe.eventHub != nil {
		pe.eventHub.Broadcast(Event{
			Type: "price_update",
			Data: models.LiveTradingData{
				Symbol:        company.Symbol,
				CompanyID:     company.ID,
				CompanyName:   company.Name,
				LTP:           newPrice,
				ChangePercent: priceChangePct.Round(2),
				Open:          pe.dayOpenPrices[companyID],
				High:          pe.dailyHighs[companyID],
				Low:           pe.dailyLows[companyID],
				Volume:        pe.dailyVolumes[companyID],
				PreviousClose: pe.previousCloses[companyID],
				Difference:    priceChange.Round(2),
				Turnover:      newPrice.Mul(decimal.NewFromInt(pe.dailyVolumes[companyID])),
				LastUpdated:   now,
			},
		})
	}

	slog.Info("PriceEngine: trade processed",
		"symbol", company.Symbol,
		"quantity", tradeQty,
		"old_price", currentPrice.StringFixed(2),
		"new_price", newPrice.StringFixed(2),
		"impact", priceChangePct.StringFixed(2)+"%",
	)

	// Fire price update callback (for trigger worker)
	if pe.onPriceUpdate != nil {
		pe.onPriceUpdate(companyID, newPrice)
	}
}

func (pe *PriceEngine) applyCircuitBreaker(companyID string, newPrice decimal.Decimal) decimal.Decimal {
	prevClose, exists := pe.previousCloses[companyID]
	if !exists || prevClose.IsZero() {
		return newPrice
	}

	maxChange := prevClose.Mul(decimal.NewFromFloat(MaxDailyChangePercent / 100.0))
	upperLimit := prevClose.Add(maxChange)
	lowerLimit := prevClose.Sub(maxChange)

	if newPrice.GreaterThan(upperLimit) {
		return upperLimit
	}
	if newPrice.LessThan(lowerLimit) {
		return lowerLimit
	}
	return newPrice
}

func (pe *PriceEngine) updateDailyCandle(companyID string, newPrice decimal.Decimal, now time.Time) {
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dayOpen := pe.dayOpenPrices[companyID]
	if dayOpen.IsZero() {
		dayOpen = newPrice
	}
	dayHigh := pe.dailyHighs[companyID]
	dayLow := pe.dailyLows[companyID]
	dayVolume := pe.dailyVolumes[companyID]

	dailyPrice := &models.StockPrice{
		CompanyID:  companyID,
		OpenPrice:  dayOpen,
		HighPrice:  dayHigh,
		LowPrice:   dayLow,
		ClosePrice: newPrice,
		Volume:     dayVolume,
		Timestamp:  dayStart,
		Timeframe:  "1D",
	}
	if err := pe.stockRepo.UpsertDailyPrice(dailyPrice); err != nil {
		slog.Warn("PriceEngine: failed to upsert daily price", "error", err)
	}
}

// ──────────────────── Market Data Endpoints ────────────────────

func (pe *PriceEngine) GetLiveTradingData() ([]models.LiveTradingData, error) {
	pe.mu.RLock()
	defer pe.mu.RUnlock()

	companies, err := pe.stockRepo.ListCompanies(200, 0)
	if err != nil {
		return nil, err
	}

	result := make([]models.LiveTradingData, 0, len(companies))
	for _, c := range companies {
		ltp := c.CurrentPrice

		prevClose := pe.previousCloses[c.ID]
		if prevClose.IsZero() {
			prevClose = ltp
		}

		diff := ltp.Sub(prevClose)
		changePct := decimal.Zero
		if !prevClose.IsZero() {
			changePct = diff.Div(prevClose).Mul(decimal.NewFromInt(100))
		}

		dayOpen := pe.dayOpenPrices[c.ID]
		if dayOpen.IsZero() {
			dayOpen = ltp
		}
		dayHigh := pe.dailyHighs[c.ID]
		if dayHigh.IsZero() {
			dayHigh = ltp
		}
		dayLow := pe.dailyLows[c.ID]
		if dayLow.IsZero() {
			dayLow = ltp
		}

		result = append(result, models.LiveTradingData{
			Symbol:        c.Symbol,
			CompanyID:     c.ID,
			CompanyName:   c.Name,
			Sector:        c.Sector,
			LTP:           ltp,
			ChangePercent: changePct.Round(2),
			Open:          dayOpen,
			High:          dayHigh,
			Low:           dayLow,
			Volume:        pe.dailyVolumes[c.ID],
			PreviousClose: prevClose,
			Difference:    diff.Round(2),
			Turnover:      ltp.Mul(decimal.NewFromInt(pe.dailyVolumes[c.ID])),
			LastUpdated:   time.Now(),
		})
	}

	return result, nil
}

func (pe *PriceEngine) GetMarketIndex() (*models.MarketIndex, error) {
	pe.mu.RLock()
	defer pe.mu.RUnlock()

	companies, err := pe.stockRepo.ListCompanies(200, 0)
	if err != nil {
		return nil, err
	}

	totalMarketCap := decimal.Zero
	prevTotalMarketCap := decimal.Zero
	totalTurnover := decimal.Zero
	var totalVolume int64
	advances, declines, unchanged := 0, 0, 0

	for _, c := range companies {
		marketCap := c.CurrentPrice.Mul(decimal.NewFromInt(c.TotalSupply))
		totalMarketCap = totalMarketCap.Add(marketCap)

		prevClose := pe.previousCloses[c.ID]
		if !prevClose.IsZero() {
			prevMC := prevClose.Mul(decimal.NewFromInt(c.TotalSupply))
			prevTotalMarketCap = prevTotalMarketCap.Add(prevMC)

			cmp := c.CurrentPrice.Cmp(prevClose)
			if cmp > 0 {
				advances++
			} else if cmp < 0 {
				declines++
			} else {
				unchanged++
			}
		} else {
			prevTotalMarketCap = prevTotalMarketCap.Add(marketCap)
			unchanged++
		}

		vol := pe.dailyVolumes[c.ID]
		totalVolume += vol
		totalTurnover = totalTurnover.Add(c.CurrentPrice.Mul(decimal.NewFromInt(vol)))
	}

	baseDivisor := decimal.NewFromFloat(1e9)
	indexValue := decimal.Zero
	if !totalMarketCap.IsZero() {
		indexValue = totalMarketCap.Div(baseDivisor)
	}
	prevIndexValue := decimal.Zero
	if !prevTotalMarketCap.IsZero() {
		prevIndexValue = prevTotalMarketCap.Div(baseDivisor)
	}

	indexChange := indexValue.Sub(prevIndexValue)
	indexChangePct := decimal.Zero
	if !prevIndexValue.IsZero() {
		indexChangePct = indexChange.Div(prevIndexValue).Mul(decimal.NewFromInt(100))
	}

	return &models.MarketIndex{
		IndexValue:     indexValue.Round(2),
		Change:         indexChange.Round(2),
		ChangePercent:  indexChangePct.Round(2),
		TotalTurnover:  totalTurnover.Round(2),
		TotalVolume:    totalVolume,
		TotalMarketCap: totalMarketCap.Round(2),
		Advances:       advances,
		Declines:       declines,
		Unchanged:      unchanged,
		TotalCompanies: len(companies),
		PreviousClose:  prevIndexValue.Round(2),
		Timestamp:      time.Now(),
	}, nil
}

func (pe *PriceEngine) GetCandlestickData(symbol string, timeframe string, days int) ([]models.CandlestickData, error) {
	company, err := pe.stockRepo.GetCompanyBySymbol(symbol)
	if err != nil {
		return nil, err
	}

	to := time.Now()
	from := to.AddDate(0, 0, -days)
	if timeframe == "" {
		timeframe = "1D"
	}

	prices, err := pe.stockRepo.GetPriceHistory(company.ID, timeframe, from, to, 1000)
	if err != nil {
		return nil, err
	}

	candles := make([]models.CandlestickData, 0, len(prices))
	for _, p := range prices {
		candles = append(candles, models.CandlestickData{
			Timestamp: p.Timestamp,
			Open:      p.OpenPrice,
			High:      p.HighPrice,
			Low:       p.LowPrice,
			Close:     p.ClosePrice,
			Volume:    p.Volume,
		})
	}

	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Timestamp.Before(candles[j].Timestamp)
	})

	return candles, nil
}
