package market

import (
	"log/slog"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/shopspring/decimal"
)

const (
	// MaxDailyChangePercent is the circuit breaker limit (like NEPSE ±10%)
	MaxDailyChangePercent = 10.0

	// BaseSensitivity controls how much a trade affects the price.
	// Higher = more volatile. 0.5 means a trade of 100% of reference volume
	// would move price by ~sqrt(1)*0.5 = 50%. Smaller trades have proportionally
	// smaller sqrt-dampened impact.
	BaseSensitivity = 0.5

	// MinPrice is the minimum allowed stock price
	MinPrice = 1.0

	// MicroNoiseRange adds tiny random noise to each trade's impact
	// to simulate market microstructure (±0.1%)
	MicroNoiseRange = 0.001
)

// TradeImpact contains the result of processing a trade through the price engine
type TradeImpact struct {
	NewPrice       decimal.Decimal `json:"new_price"`
	PriceChange    decimal.Decimal `json:"price_change"`
	PriceChangePct decimal.Decimal `json:"price_change_pct"`
	DayOpen        decimal.Decimal `json:"day_open"`
	DayHigh        decimal.Decimal `json:"day_high"`
	DayLow         decimal.Decimal `json:"day_low"`
	DayVolume      int64           `json:"day_volume"`
	PreviousClose  decimal.Decimal `json:"previous_close"`
}

// PriceEngine implements order-driven price discovery.
// Instead of random price changes, prices move based on actual buy/sell activity.
//
// Key principles:
//   - Buy orders push prices UP (demand pressure)
//   - Sell orders push prices DOWN (supply pressure)
//   - Impact is proportional to trade size relative to available liquidity
//   - Square-root dampening prevents large trades from causing unrealistic moves
//   - Circuit breaker limits daily movement to ±10% (like NEPSE)
//   - Micro noise simulates other market participants
type PriceEngine struct {
	stockRepo stock.Repository
	eventHub  *EventHub

	mu             sync.RWMutex
	dayOpenPrices  map[string]decimal.Decimal // companyID -> day's opening price
	previousCloses map[string]decimal.Decimal // companyID -> previous day's close
	dailyVolumes   map[string]int64           // companyID -> today's total volume
	dailyHighs     map[string]decimal.Decimal // companyID -> today's high
	dailyLows      map[string]decimal.Decimal // companyID -> today's low
	lastTradeDay   time.Time
}

// NewPriceEngine creates a new PriceEngine and initializes day tracking data
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

// initializeDayData loads current prices as starting points for the trading day
func (pe *PriceEngine) initializeDayData() {
	companies, err := pe.stockRepo.ListCompanies(200, 0)
	if err != nil {
		slog.Error("PriceEngine: failed to load companies", "error", err)
		return
	}

	for _, c := range companies {
		price, err := pe.stockRepo.GetLatestPrice(c.ID)
		if err != nil {
			continue
		}
		pe.previousCloses[c.ID] = price.ClosePrice
		pe.dayOpenPrices[c.ID] = price.ClosePrice
		pe.dailyHighs[c.ID] = price.ClosePrice
		pe.dailyLows[c.ID] = price.ClosePrice
	}

	slog.Info("PriceEngine initialized", "companies_loaded", len(pe.previousCloses))
}

// checkDayReset resets daily tracking when a new trading day starts
func (pe *PriceEngine) checkDayReset() {
	today := time.Now().Truncate(24 * time.Hour)
	if today.After(pe.lastTradeDay) {
		slog.Info("PriceEngine: new trading day detected, resetting daily data")

		// Previous close = last known price for each company
		for companyID := range pe.dayOpenPrices {
			if lastPrice, err := pe.stockRepo.GetLatestPrice(companyID); err == nil {
				pe.previousCloses[companyID] = lastPrice.ClosePrice
			}
		}

		// Reset intraday tracking
		pe.dayOpenPrices = make(map[string]decimal.Decimal)
		pe.dailyVolumes = make(map[string]int64)
		pe.dailyHighs = make(map[string]decimal.Decimal)
		pe.dailyLows = make(map[string]decimal.Decimal)
		pe.lastTradeDay = today
	}
}

// ProcessTrade is the core method called after every buy/sell trade.
// It calculates realistic price impact based on supply/demand and:
//  1. Computes new price using order-impact formula
//  2. Applies circuit breaker (±10% daily limit)
//  3. Updates intraday tracking (day OHLCV)
//  4. Creates new stock_price record
//  5. Updates daily candle
//  6. Updates company market cap
//  7. Broadcasts real-time event via SSE
func (pe *PriceEngine) ProcessTrade(company *models.Company, quantity int, currentPrice decimal.Decimal, isBuy bool) (*TradeImpact, error) {
	pe.mu.Lock()
	defer pe.mu.Unlock()

	pe.checkDayReset()

	// 1. Calculate new price based on supply/demand
	newPrice := pe.calculatePriceImpact(company, quantity, currentPrice, isBuy)

	// 2. Apply circuit breaker
	newPrice = pe.applyCircuitBreaker(company.ID, newPrice)

	// 3. Update intraday tracking
	pe.dailyVolumes[company.ID] += int64(quantity)

	if pe.dayOpenPrices[company.ID].IsZero() {
		pe.dayOpenPrices[company.ID] = currentPrice
	}

	if pe.dailyHighs[company.ID].IsZero() || newPrice.GreaterThan(pe.dailyHighs[company.ID]) {
		pe.dailyHighs[company.ID] = newPrice
	}
	if pe.dailyLows[company.ID].IsZero() || newPrice.LessThan(pe.dailyLows[company.ID]) {
		pe.dailyLows[company.ID] = newPrice
	}

	now := time.Now()

	// 4. Create 1-minute candle record
	stockPrice := &models.StockPrice{
		CompanyID:  company.ID,
		OpenPrice:  currentPrice,
		HighPrice:  decimal.Max(currentPrice, newPrice),
		LowPrice:   decimal.Min(currentPrice, newPrice),
		ClosePrice: newPrice,
		Volume:     int64(quantity),
		Timestamp:  now,
		Timeframe:  "1m",
	}

	if err := pe.stockRepo.CreateStockPrice(stockPrice); err != nil {
		slog.Error("PriceEngine: failed to create stock price", "error", err)
		return nil, err
	}

	// 5. Update daily candle (1D timeframe)
	pe.updateDailyCandle(company.ID, newPrice, now)

	// 6. Calculate impact metrics
	priceChange := newPrice.Sub(currentPrice)
	priceChangePct := decimal.Zero
	if !currentPrice.IsZero() {
		priceChangePct = priceChange.Div(currentPrice).Mul(decimal.NewFromInt(100))
	}

	// 7. Update company market cap
	newMarketCap := newPrice.Mul(decimal.NewFromInt(company.TotalShares))
	if err := pe.stockRepo.UpdateCompanyMarketCap(company.ID, newMarketCap); err != nil {
		slog.Warn("PriceEngine: failed to update market cap", "error", err)
	}

	impact := &TradeImpact{
		NewPrice:       newPrice,
		PriceChange:    priceChange,
		PriceChangePct: priceChangePct.Round(2),
		DayOpen:        pe.dayOpenPrices[company.ID],
		DayHigh:        pe.dailyHighs[company.ID],
		DayLow:         pe.dailyLows[company.ID],
		DayVolume:      pe.dailyVolumes[company.ID],
		PreviousClose:  pe.previousCloses[company.ID],
	}

	// 8. Broadcast real-time event
	if pe.eventHub != nil {
		tradeType := "buy"
		if !isBuy {
			tradeType = "sell"
		}

		prevClose := pe.previousCloses[company.ID]
		diffFromPrevClose := newPrice.Sub(prevClose)
		changePctFromPrevClose := decimal.Zero
		if !prevClose.IsZero() {
			changePctFromPrevClose = diffFromPrevClose.Div(prevClose).Mul(decimal.NewFromInt(100))
		}

		pe.eventHub.Broadcast(Event{
			Type: "trade",
			Data: models.TradeFeedItem{
				Symbol:      company.Symbol,
				CompanyName: company.Name,
				TradeType:   tradeType,
				Quantity:    quantity,
				Price:       currentPrice,
				TotalAmount: currentPrice.Mul(decimal.NewFromInt(int64(quantity))),
				PriceImpact: priceChange,
				NewPrice:    newPrice,
				Timestamp:   now,
			},
		})

		pe.eventHub.Broadcast(Event{
			Type: "price_update",
			Data: models.LiveTradingData{
				Symbol:        company.Symbol,
				CompanyID:     company.ID,
				CompanyName:   company.Name,
				LTP:           newPrice,
				ChangePercent: changePctFromPrevClose.Round(2),
				Open:          pe.dayOpenPrices[company.ID],
				High:          pe.dailyHighs[company.ID],
				Low:           pe.dailyLows[company.ID],
				Volume:        pe.dailyVolumes[company.ID],
				PreviousClose: prevClose,
				Difference:    diffFromPrevClose.Round(2),
				Turnover:      newPrice.Mul(decimal.NewFromInt(pe.dailyVolumes[company.ID])),
				LastUpdated:   now,
			},
		})
	}

	slog.Info("PriceEngine: trade processed",
		"symbol", company.Symbol,
		"type", map[bool]string{true: "buy", false: "sell"}[isBuy],
		"quantity", quantity,
		"old_price", currentPrice.StringFixed(2),
		"new_price", newPrice.StringFixed(2),
		"impact", priceChangePct.StringFixed(2)+"%",
	)

	return impact, nil
}

// calculatePriceImpact determines the new price based on trade parameters.
//
// Formula: newPrice = currentPrice × (1 + direction × sensitivity × √(tradeQty / referenceShares) + noise)
//
// This models real market behavior where:
//   - Price impact is proportional to trade size
//   - Square-root dampening prevents large orders from moving price linearly
//   - More liquid stocks (more available shares) are less affected
//   - Small random noise simulates market microstructure
func (pe *PriceEngine) calculatePriceImpact(company *models.Company, quantity int, currentPrice decimal.Decimal, isBuy bool) decimal.Decimal {
	// Reference volume = available shares (liquidity indicator)
	referenceShares := company.AvailableShares
	if referenceShares <= 0 {
		referenceShares = company.TotalShares
	}
	if referenceShares <= 0 {
		referenceShares = 1000000 // fallback
	}

	// Liquidity ratio: what fraction of available supply this trade represents
	liquidityRatio := float64(quantity) / float64(referenceShares)

	// Square-root dampening: prevents large trades from having linear impact
	// A trade of 1% of available shares moves price by ~0.05%
	// A trade of 100% would move price by ~50% (but circuit breaker limits this)
	impact := BaseSensitivity * math.Sqrt(liquidityRatio)

	// Apply direction
	if !isBuy {
		impact = -impact
	}

	// Add micro noise for realism
	noise := (rand.Float64() - 0.5) * 2 * MicroNoiseRange
	impact += noise

	// Calculate new price
	impactDecimal := decimal.NewFromFloat(impact)
	newPrice := currentPrice.Mul(decimal.NewFromInt(1).Add(impactDecimal))

	// Enforce minimum price
	minPriceDecimal := decimal.NewFromFloat(MinPrice)
	if newPrice.LessThan(minPriceDecimal) {
		newPrice = minPriceDecimal
	}

	// Round to 2 decimal places
	return newPrice.Round(2)
}

// applyCircuitBreaker limits daily price movement to ±MaxDailyChangePercent
// from the previous day's close. This prevents extreme volatility.
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

// updateDailyCandle creates/updates the 1D timeframe candle for today
func (pe *PriceEngine) updateDailyCandle(companyID string, newPrice decimal.Decimal, now time.Time) {
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	dayOpen := pe.dayOpenPrices[companyID]
	if dayOpen.IsZero() {
		dayOpen = newPrice
	}
	dayHigh := pe.dailyHighs[companyID]
	if dayHigh.IsZero() {
		dayHigh = newPrice
	}
	dayLow := pe.dailyLows[companyID]
	if dayLow.IsZero() {
		dayLow = newPrice
	}
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
		slog.Warn("PriceEngine: failed to upsert daily price", "companyID", companyID, "error", err)
	}
}

// GetLiveTradingData returns the current live trading snapshot for all companies
// This is the data shown in the "Live Trading" table
func (pe *PriceEngine) GetLiveTradingData() ([]models.LiveTradingData, error) {
	pe.mu.RLock()
	defer pe.mu.RUnlock()

	companies, err := pe.stockRepo.ListCompanies(200, 0)
	if err != nil {
		return nil, err
	}

	result := make([]models.LiveTradingData, 0, len(companies))

	for _, company := range companies {
		price, err := pe.stockRepo.GetLatestPrice(company.ID)
		if err != nil {
			continue
		}

		ltp := price.ClosePrice

		prevClose := pe.previousCloses[company.ID]
		if prevClose.IsZero() {
			prevClose = price.OpenPrice
		}

		diff := ltp.Sub(prevClose)
		changePct := decimal.Zero
		if !prevClose.IsZero() {
			changePct = diff.Div(prevClose).Mul(decimal.NewFromInt(100))
		}

		dayOpen := pe.dayOpenPrices[company.ID]
		if dayOpen.IsZero() {
			dayOpen = price.OpenPrice
		}
		dayHigh := pe.dailyHighs[company.ID]
		if dayHigh.IsZero() {
			dayHigh = ltp
		}
		dayLow := pe.dailyLows[company.ID]
		if dayLow.IsZero() {
			dayLow = ltp
		}
		dayVolume := pe.dailyVolumes[company.ID]

		turnover := ltp.Mul(decimal.NewFromInt(dayVolume))

		result = append(result, models.LiveTradingData{
			Symbol:        company.Symbol,
			CompanyID:     company.ID,
			CompanyName:   company.Name,
			Sector:        company.Sector,
			LTP:           ltp,
			ChangePercent: changePct.Round(2),
			Open:          dayOpen,
			High:          dayHigh,
			Low:           dayLow,
			Volume:        dayVolume,
			PreviousClose: prevClose,
			Difference:    diff.Round(2),
			Turnover:      turnover,
			LastUpdated:   price.Timestamp,
		})
	}

	return result, nil
}

// GetMarketIndex calculates the overall market index value from all companies.
// Like NEPSE index, it represents the total market performance in a single number.
//
// Index = TotalMarketCap / BaseDivisor
// TotalMarketCap = Σ(stock_price_i × total_shares_i)
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
	totalVolume := int64(0)
	advances := 0
	declines := 0
	unchanged := 0

	for _, company := range companies {
		price, err := pe.stockRepo.GetLatestPrice(company.ID)
		if err != nil {
			unchanged++
			continue
		}

		// Current market cap
		marketCap := price.ClosePrice.Mul(decimal.NewFromInt(company.TotalShares))
		totalMarketCap = totalMarketCap.Add(marketCap)

		// Previous day market cap
		prevClose := pe.previousCloses[company.ID]
		if !prevClose.IsZero() {
			prevMarketCap := prevClose.Mul(decimal.NewFromInt(company.TotalShares))
			prevTotalMarketCap = prevTotalMarketCap.Add(prevMarketCap)
		} else {
			prevTotalMarketCap = prevTotalMarketCap.Add(marketCap)
		}

		// Intraday data
		dayVol := pe.dailyVolumes[company.ID]
		totalVolume += dayVol
		turnover := price.ClosePrice.Mul(decimal.NewFromInt(dayVol))
		totalTurnover = totalTurnover.Add(turnover)

		// Advance/decline tracking
		if !prevClose.IsZero() {
			cmp := price.ClosePrice.Cmp(prevClose)
			if cmp > 0 {
				advances++
			} else if cmp < 0 {
				declines++
			} else {
				unchanged++
			}
		} else {
			unchanged++
		}
	}

	// Index = total market cap / base divisor (1 billion)
	baseDivisor := decimal.NewFromFloat(1000000000)
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

// GetMarketSummary returns a comprehensive market overview with index,
// top gainers, top losers, most active stocks, and sector breakdown
func (pe *PriceEngine) GetMarketSummary() (*models.MarketSummary, error) {
	index, err := pe.GetMarketIndex()
	if err != nil {
		return nil, err
	}

	liveData, err := pe.GetLiveTradingData()
	if err != nil {
		return nil, err
	}

	// Sort for top gainers (highest % change)
	sortedByChange := make([]models.LiveTradingData, len(liveData))
	copy(sortedByChange, liveData)
	sort.Slice(sortedByChange, func(i, j int) bool {
		return sortedByChange[i].ChangePercent.GreaterThan(sortedByChange[j].ChangePercent)
	})

	topGainers := make([]models.LiveTradingData, 0, 5)
	for i, d := range sortedByChange {
		if i >= 5 {
			break
		}
		if d.ChangePercent.IsPositive() {
			topGainers = append(topGainers, d)
		}
	}

	// Top losers (lowest % change)
	topLosers := make([]models.LiveTradingData, 0, 5)
	for i := len(sortedByChange) - 1; i >= 0; i-- {
		if len(topLosers) >= 5 {
			break
		}
		if sortedByChange[i].ChangePercent.IsNegative() {
			topLosers = append(topLosers, sortedByChange[i])
		}
	}

	// Most active by volume
	sortedByVolume := make([]models.LiveTradingData, len(liveData))
	copy(sortedByVolume, liveData)
	sort.Slice(sortedByVolume, func(i, j int) bool {
		return sortedByVolume[i].Volume > sortedByVolume[j].Volume
	})
	mostActive := sortedByVolume
	if len(mostActive) > 5 {
		mostActive = mostActive[:5]
	}

	// Sector summary
	sectorMap := make(map[string]*models.SectorIndex)
	for _, d := range liveData {
		si, exists := sectorMap[d.Sector]
		if !exists {
			si = &models.SectorIndex{
				Sector:   d.Sector,
				Turnover: decimal.Zero,
				Change:   decimal.Zero,
			}
			sectorMap[d.Sector] = si
		}
		si.CompanyCount++
		si.Volume += d.Volume
		si.Turnover = si.Turnover.Add(d.Turnover)
		si.Change = si.Change.Add(d.ChangePercent)

		if d.ChangePercent.IsPositive() {
			si.Advances++
		} else if d.ChangePercent.IsNegative() {
			si.Declines++
		}
	}

	sectorSummary := make([]models.SectorIndex, 0, len(sectorMap))
	for _, si := range sectorMap {
		if si.CompanyCount > 0 {
			si.ChangePercent = si.Change.Div(decimal.NewFromInt(int64(si.CompanyCount))).Round(2)
		}
		sectorSummary = append(sectorSummary, *si)
	}
	sort.Slice(sectorSummary, func(i, j int) bool {
		return sectorSummary[i].Turnover.GreaterThan(sectorSummary[j].Turnover)
	})

	return &models.MarketSummary{
		Index:         *index,
		TopGainers:    topGainers,
		TopLosers:     topLosers,
		MostActive:    mostActive,
		SectorSummary: sectorSummary,
		AsOf:          time.Now(),
	}, nil
}

// GetCandlestickData returns OHLCV candle data for charting
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

	// Sort by timestamp ascending for charting
	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Timestamp.Before(candles[j].Timestamp)
	})

	return candles, nil
}
