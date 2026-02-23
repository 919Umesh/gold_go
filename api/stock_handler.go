package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/pkg/apperr"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type StockHandler struct {
	service     stock.Service
	priceEngine *market.PriceEngine
	eventHub    *market.EventHub
}

func NewStockHandler(service stock.Service, priceEngine *market.PriceEngine, eventHub *market.EventHub) *StockHandler {
	return &StockHandler{
		service:     service,
		priceEngine: priceEngine,
		eventHub:    eventHub,
	}
}

func (h *StockHandler) GetCompany(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Symbol parameter required")
		return
	}

	company, err := h.service.GetCompany(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "Company not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"company": company})
}

func (h *StockHandler) ListCompanies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	companies, total, err := h.service.ListCompaniesWithTotal(limit, offset)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch companies")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"companies": companies,
		"limit":     limit,
		"offset":    offset,
		"count":     len(companies),
		"total":     total,
	})
}

func (h *StockHandler) SearchCompanies(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Search query required")
		return
	}

	companies, err := h.service.SearchCompanies(query)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Search failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"companies": companies})
}

func (h *StockHandler) GetCompaniesBySector(c *gin.Context) {
	sector := c.Param("sector")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	companies, total, err := h.service.GetCompaniesBySectorWithTotal(sector, limit, offset)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch companies")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sector":    sector,
		"companies": companies,
		"limit":     limit,
		"offset":    offset,
		"count":     len(companies),
		"total":     total,
	})
}

func (h *StockHandler) GetCurrentPrice(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Symbol parameter required")
		return
	}

	price, err := h.service.GetCurrentPrice(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": price})
}

func (h *StockHandler) GetPriceHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Symbol parameter required")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "1D")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	if days <= 0 {
		days = 30
	}
	if days > 365*5 {
		days = 365 * 5
	}

	// Validate timeframe: only allow known values
	// "all" returns all timeframes, "1D" is daily candles, etc.
	validTimeframes := map[string]bool{
		"1D": true, "1W": true, "1M": true, "all": true,
	}
	if !validTimeframes[timeframe] {
		timeframe = "1D"
	}

	prices, err := h.service.GetPriceHistory(symbol, timeframe, days)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"symbol":    symbol,
		"timeframe": timeframe,
		"days":      days,
		"count":     len(prices),
		"prices":    prices,
	})
}

func (h *StockHandler) GetMarketOverview(c *gin.Context) {
	overview, err := h.service.GetMarketOverview()
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch market overview")
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (h *StockHandler) GetTopGainers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	gainers, err := h.service.GetTopGainers(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch top gainers")
		return
	}

	c.JSON(http.StatusOK, gin.H{"gainers": gainers})
}

func (h *StockHandler) GetTopLosers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	losers, err := h.service.GetTopLosers(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch top losers")
		return
	}

	c.JSON(http.StatusOK, gin.H{"losers": losers})
}

func (h *StockHandler) GetMostActive(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	active, err := h.service.GetMostActive(limit)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch most active")
		return
	}

	c.JSON(http.StatusOK, gin.H{"active": active})
}

func (h *StockHandler) GetUpcomingEvents(c *gin.Context) {
	symbol := c.Param("symbol")

	events, err := h.service.GetUpcomingEvents(symbol)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, "Events not available")
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *StockHandler) GetAllSectors(c *gin.Context) {
	sectors, err := h.service.GetAllSectors()
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch sectors")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sectors retrieved successfully",
		"sectors": sectors,
		"count":   len(sectors),
	})
}

func (h *StockHandler) GetSectorStats(c *gin.Context) {
	sector := c.Param("sector")

	companies, err := h.service.GetCompaniesBySector(sector, 1000, 0)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch sector stats")
		return
	}

	if len(companies) == 0 {
		apperr.RespondWithMessage(c, http.StatusNotFound, "no companies found in this sector")
		return
	}

	totalMarketCap := decimal.Zero
	totalEmployees := 0
	avgYear := 0

	for _, comp := range companies {
		totalMarketCap = totalMarketCap.Add(comp.MarketCap)
		totalEmployees += comp.Employees
		avgYear += comp.FoundedYear
	}
	avgYear = avgYear / len(companies)

	topCompanies := companies
	if len(topCompanies) > 5 {
		topCompanies = topCompanies[:5]
	}

	avgMarketCap := decimal.Zero
	if len(companies) > 0 {
		avgMarketCap = totalMarketCap.Div(decimal.NewFromInt(int64(len(companies))))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sector statistics",
		"sector":  sector,
		"statistics": gin.H{
			"company_count":    len(companies),
			"total_market_cap": totalMarketCap,
			"avg_market_cap":   avgMarketCap,
			"total_employees":  totalEmployees,
			"avg_employees":    totalEmployees / len(companies),
			"avg_founded_year": avgYear,
		},
		"top_5_companies": topCompanies,
	})
}

// =====================================================
// LIVE TRADING DATA
// Returns real-time trading data for all stocks
// Format: Symbol | LTP | %Change | Open | High | Low | Qty | PClose | Diff
// =====================================================

func (h *StockHandler) GetLiveTradingData(c *gin.Context) {
	if h.priceEngine == nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Price engine not available")
		return
	}

	data, err := h.priceEngine.GetLiveTradingData()
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to fetch live trading data")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"live_trading": data,
		"count":        len(data),
	})
}

// =====================================================
// CANDLESTICK / OHLCV DATA
// Returns candle data for charting (like the price chart in the screenshot)
// =====================================================

func (h *StockHandler) GetCandlestickData(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		apperr.RespondWithMessage(c, http.StatusBadRequest, "Symbol parameter required")
		return
	}

	timeframe := c.DefaultQuery("timeframe", "1D")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	if days > 365 {
		days = 365
	}

	if h.priceEngine == nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Price engine not available")
		return
	}

	candles, err := h.priceEngine.GetCandlestickData(symbol, timeframe, days)
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"symbol":    symbol,
		"timeframe": timeframe,
		"days":      days,
		"candles":   candles,
		"count":     len(candles),
	})
}

// =====================================================
// MARKET INDEX
// Returns overall market index (like NEPSE index)
// Shows total market value, advances, declines, turnover
// =====================================================

func (h *StockHandler) GetMarketIndex(c *gin.Context) {
	if h.priceEngine == nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Price engine not available")
		return
	}

	index, err := h.priceEngine.GetMarketIndex()
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to calculate market index")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"market_index": index,
	})
}

// =====================================================
// MARKET SUMMARY
// Comprehensive market overview with index, gainers, losers,
// most active, and sector breakdown
// =====================================================

func (h *StockHandler) GetMarketSummary(c *gin.Context) {
	if h.priceEngine == nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Price engine not available")
		return
	}

	summary, err := h.priceEngine.GetMarketSummary()
	if err != nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Failed to get market summary")
		return
	}

	c.JSON(http.StatusOK, summary)
}

// =====================================================
// SSE STREAM
// Server-Sent Events endpoint for real-time price updates
// Clients connect and receive live trade/price events
// =====================================================

func (h *StockHandler) StreamPrices(c *gin.Context) {
	if h.eventHub == nil {
		apperr.RespondWithMessage(c, http.StatusInternalServerError, "Event hub not available")
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Subscribe to events
	ch := h.eventHub.Subscribe()
	defer h.eventHub.Unsubscribe(ch)

	// Stream events to client
	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-ch:
			if !ok {
				return false
			}
			c.SSEvent("message", msg)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}
