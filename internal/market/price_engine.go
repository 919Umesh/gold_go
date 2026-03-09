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
	SensitivityFactor     = 2.0
	MaxDailyChangePercent = 10.0
	MinPrice              = 0.01
)

type TradeImpact struct {
	NewPrice       decimal.Decimal `json:"new_price"`
	PriceChange    decimal.Decimal `json:"price_change"`
	PriceChangePct decimal.Decimal `json:"price_change_pct"`
}

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

	totalSupply := company.TotalSupply
	if totalSupply <= 0 {
		totalSupply = 10000
	}

	volumeRatio := float64(tradeQty) / float64(totalSupply)
	impact := volumeRatio * SensitivityFactor

	direction := 1.0
	if tradePrice.LessThan(currentPrice) {
		direction = -1.0
	}

	noise := (rand.Float64() - 0.5) * 0.002
	impact = impact*direction + noise

	impactDecimal := decimal.NewFromFloat(impact)
	newPrice := currentPrice.Mul(decimal.NewFromInt(1).Add(impactDecimal))

	minP := decimal.NewFromFloat(MinPrice)
	if newPrice.LessThan(minP) {
		newPrice = minP
	}

	newPrice = pe.applyCircuitBreaker(companyID, newPrice)
	newPrice = newPrice.Round(2)

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
	turnover := newPrice.Mul(decimal.NewFromInt(tradeQty))

	priceChange := newPrice.Sub(currentPrice)
	priceChangePct := decimal.Zero
	if !currentPrice.IsZero() {
		priceChangePct = priceChange.Div(currentPrice).Mul(decimal.NewFromInt(100))
	}

	stockPrice := &models.StockPrice{
		CompanyID:     companyID,
		OpenPrice:     currentPrice,
		HighPrice:     decimal.Max(currentPrice, newPrice),
		LowPrice:      decimal.Min(currentPrice, newPrice),
		ClosePrice:    newPrice,
		Volume:        tradeQty,
		Turnover:      turnover,
		ChangePercent: priceChangePct.Round(4),
		Timestamp:     now,
		Timeframe:     "1m",
	}
	if err := pe.stockRepo.CreateStockPrice(stockPrice); err != nil {
		slog.Error("PriceEngine: failed to create stock price", "error", err)
	}

	pe.updateDailyCandle(companyID, newPrice, now)

	newMarketCap := newPrice.Mul(decimal.NewFromInt(company.TotalSupply))
	if err := pe.stockRepo.UpdateCompanyPrice(companyID, newPrice.String(), newMarketCap.String()); err != nil {
		slog.Warn("PriceEngine: failed to update company price", "error", err)
	}

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

	changePct := decimal.Zero
	prevClose := pe.previousCloses[companyID]
	if !prevClose.IsZero() {
		changePct = newPrice.Sub(prevClose).Div(prevClose).Mul(decimal.NewFromInt(100))
	}

	dailyTurnover := newPrice.Mul(decimal.NewFromInt(dayVolume))

	dailyPrice := &models.StockPrice{
		CompanyID:     companyID,
		OpenPrice:     dayOpen,
		HighPrice:     dayHigh,
		LowPrice:      dayLow,
		ClosePrice:    newPrice,
		Volume:        dayVolume,
		Turnover:      dailyTurnover,
		ChangePercent: changePct.Round(4),
		Timestamp:     dayStart,
		Timeframe:     "1D",
	}
	if err := pe.stockRepo.UpsertDailyPrice(dailyPrice); err != nil {
		slog.Warn("PriceEngine: failed to upsert daily price", "error", err)
	}
}

// ──────────────────── Market Data Methods ────────────────────

func (pe *PriceEngine) GetLiveTradingData() ([]models.LiveTradingData, error) {
	pe.mu.RLock()
	defer pe.mu.RUnlock()

	companies, err := pe.stockRepo.ListCompanies(200, 0)
	if err != nil {
		return nil, err
	}

	result := make([]models.LiveTradingData, 0, len(companies))
	for _, c := range companies {
		result = append(result, pe.buildLiveTradingData(c))
	}

	return result, nil
}

func (pe *PriceEngine) buildLiveTradingData(c models.Company) models.LiveTradingData {
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

	return models.LiveTradingData{
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
	}
}

func (pe *PriceEngine) GetTopGainers(limit int) ([]models.LiveTradingData, error) {
	data, err := pe.GetLiveTradingData()
	if err != nil {
		return nil, err
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].ChangePercent.GreaterThan(data[j].ChangePercent)
	})
	if len(data) > limit {
		data = data[:limit]
	}
	return data, nil
}

func (pe *PriceEngine) GetTopLosers(limit int) ([]models.LiveTradingData, error) {
	data, err := pe.GetLiveTradingData()
	if err != nil {
		return nil, err
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].ChangePercent.LessThan(data[j].ChangePercent)
	})
	if len(data) > limit {
		data = data[:limit]
	}
	return data, nil
}

func (pe *PriceEngine) GetMostActive(limit int) ([]models.LiveTradingData, error) {
	data, err := pe.GetLiveTradingData()
	if err != nil {
		return nil, err
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Volume > data[j].Volume
	})
	if len(data) > limit {
		data = data[:limit]
	}
	return data, nil
}

func (pe *PriceEngine) GetTopTurnover(limit int) ([]models.LiveTradingData, error) {
	data, err := pe.GetLiveTradingData()
	if err != nil {
		return nil, err
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Turnover.GreaterThan(data[j].Turnover)
	})
	if len(data) > limit {
		data = data[:limit]
	}
	return data, nil
}

func (pe *PriceEngine) GetTopSectors() ([]models.SectorPerformance, error) {
	data, err := pe.GetLiveTradingData()
	if err != nil {
		return nil, err
	}

	sectorMap := make(map[string]*models.SectorPerformance)

	for _, d := range data {
		sp, ok := sectorMap[d.Sector]
		if !ok {
			sp = &models.SectorPerformance{
				Sector: d.Sector,
			}
			sectorMap[d.Sector] = sp
		}
		sp.CompanyCount++
		sp.AvgChange = sp.AvgChange.Add(d.ChangePercent)
		sp.TotalTurnover = sp.TotalTurnover.Add(d.Turnover)
		sp.TotalVolume += d.Volume
		mc := d.LTP.Mul(decimal.NewFromInt(d.Volume)) // approximate
		sp.TotalMarketCap = sp.TotalMarketCap.Add(mc)
	}

	result := make([]models.SectorPerformance, 0, len(sectorMap))
	for _, sp := range sectorMap {
		if sp.CompanyCount > 0 {
			sp.AvgChange = sp.AvgChange.Div(decimal.NewFromInt(int64(sp.CompanyCount))).Round(2)
		}
		result = append(result, *sp)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].AvgChange.GreaterThan(result[j].AvgChange)
	})

	return result, nil
}

func (pe *PriceEngine) GetCompaniesBySector(sector string) ([]models.LiveTradingData, error) {
	pe.mu.RLock()
	defer pe.mu.RUnlock()

	companies, err := pe.stockRepo.ListCompaniesBySector(sector, 100, 0)
	if err != nil {
		return nil, err
	}

	result := make([]models.LiveTradingData, 0, len(companies))
	for _, c := range companies {
		result = append(result, pe.buildLiveTradingData(c))
	}

	return result, nil
}

func (pe *PriceEngine) GetNewCompanies(limit int) ([]models.LiveTradingData, error) {
	pe.mu.RLock()
	defer pe.mu.RUnlock()

	companies, err := pe.stockRepo.ListCompanies(200, 0)
	if err != nil {
		return nil, err
	}

	sort.Slice(companies, func(i, j int) bool {
		return companies[i].CreatedAt.After(companies[j].CreatedAt)
	})

	if len(companies) > limit {
		companies = companies[:limit]
	}

	result := make([]models.LiveTradingData, 0, len(companies))
	for _, c := range companies {
		result = append(result, pe.buildLiveTradingData(c))
	}

	return result, nil
}

func (pe *PriceEngine) GetOldCompanies(limit int) ([]models.LiveTradingData, error) {
	pe.mu.RLock()
	defer pe.mu.RUnlock()

	companies, err := pe.stockRepo.ListCompanies(200, 0)
	if err != nil {
		return nil, err
	}

	sort.Slice(companies, func(i, j int) bool {
		return companies[i].CreatedAt.Before(companies[j].CreatedAt)
	})

	if len(companies) > limit {
		companies = companies[:limit]
	}

	result := make([]models.LiveTradingData, 0, len(companies))
	for _, c := range companies {
		result = append(result, pe.buildLiveTradingData(c))
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
			Timestamp:     p.Timestamp,
			Open:          p.OpenPrice,
			High:          p.HighPrice,
			Low:           p.LowPrice,
			Close:         p.ClosePrice,
			Volume:        p.Volume,
			Turnover:      p.Turnover,
			ChangePercent: p.ChangePercent,
		})
	}

	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Timestamp.Before(candles[j].Timestamp)
	})

	return candles, nil
}
