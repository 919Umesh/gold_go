package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/919Umesh/stock_market_sim/internal/market"
	"github.com/919Umesh/stock_market_sim/internal/stock"
	"github.com/919Umesh/stock_market_sim/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type MarketHandler struct {
	priceEngine   *market.PriceEngine
	stockService  stock.Service
	eventHub      *market.EventHub
	triggerWorker *market.TriggerWorker
}

func NewMarketHandler(pe *market.PriceEngine, ss stock.Service, eh *market.EventHub, tw *market.TriggerWorker) *MarketHandler {
	return &MarketHandler{
		priceEngine:   pe,
		stockService:  ss,
		eventHub:      eh,
		triggerWorker: tw,
	}
}

// ──────────────────── Companies ────────────────────

// ListCompanies godoc
// @Summary List all companies
// @Tags Market
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Router /market/companies [get]
func (h *MarketHandler) ListCompanies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	companies, err := h.stockService.ListCompanies(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companies, "count": len(companies)})
}

// GetCompanyDetail godoc
// @Summary Get company detail
// @Tags Market
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} models.Company
// @Router /market/companies/{id} [get]
func (h *MarketHandler) GetCompanyDetail(c *gin.Context) {
	id := c.Param("id")
	company, err := h.stockService.GetCompanyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": company})
}

// ──────────────────── Market Data ────────────────────

// GetLiveTradingData godoc
// @Summary Get live trading data for all companies
// @Tags Market
// @Produce json
// @Success 200 {array} models.LiveTradingData
// @Router /market/live [get]
func (h *MarketHandler) GetLiveTradingData(c *gin.Context) {
	data, err := h.priceEngine.GetLiveTradingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "count": len(data)})
}

// GetMarketIndex godoc
// @Summary Get market index summary
// @Tags Market
// @Produce json
// @Success 200 {object} models.MarketIndex
// @Router /market/index [get]
func (h *MarketHandler) GetMarketIndex(c *gin.Context) {
	index, err := h.priceEngine.GetMarketIndex()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": index})
}

// GetCandlestickData godoc
// @Summary Get candlestick OHLCV data for a company
// @Tags Market
// @Produce json
// @Param symbol query string true "Stock symbol"
// @Param timeframe query string false "Timeframe (1m, 5m, 15m, 1h, 1D)" default(1D)
// @Param days query int false "Number of days" default(90)
// @Success 200 {array} models.CandlestickData
// @Router /market/candlestick [get]
func (h *MarketHandler) GetCandlestickData(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}
	timeframe := c.DefaultQuery("timeframe", "1D")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	if days <= 0 {
		days = 90
	}

	candles, err := h.priceEngine.GetCandlestickData(symbol, timeframe, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": candles, "symbol": symbol, "timeframe": timeframe, "count": len(candles)})
}

// ──────────────────── Top Gainers / Losers / Active ────────────────────

// GetTopGainers godoc
// @Summary Get top gaining stocks
// @Tags Market
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.LiveTradingData
// @Router /market/top-gainers [get]
func (h *MarketHandler) GetTopGainers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.priceEngine.GetTopGainers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "count": len(data)})
}

// GetTopLosers godoc
// @Summary Get top losing stocks
// @Tags Market
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.LiveTradingData
// @Router /market/top-losers [get]
func (h *MarketHandler) GetTopLosers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.priceEngine.GetTopLosers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "count": len(data)})
}

// GetMostActive godoc
// @Summary Get most actively traded stocks
// @Tags Market
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.LiveTradingData
// @Router /market/most-active [get]
func (h *MarketHandler) GetMostActive(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.priceEngine.GetMostActive(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "count": len(data)})
}

// GetTopTurnover godoc
// @Summary Get stocks with highest turnover
// @Tags Market
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.LiveTradingData
// @Router /market/top-turnover [get]
func (h *MarketHandler) GetTopTurnover(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.priceEngine.GetTopTurnover(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "count": len(data)})
}

// ──────────────────── Sectors ────────────────────

// GetTopSectors godoc
// @Summary Get sector performance summary
// @Tags Market
// @Produce json
// @Success 200 {array} models.SectorPerformance
// @Router /market/sectors [get]
func (h *MarketHandler) GetTopSectors(c *gin.Context) {
	sectors, err := h.priceEngine.GetTopSectors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sectors, "count": len(sectors)})
}

// GetCompaniesBySector godoc
// @Summary Get companies in a sector
// @Tags Market
// @Produce json
// @Param sector path string true "Sector name"
// @Success 200 {array} models.LiveTradingData
// @Router /market/sectors/{sector}/companies [get]
func (h *MarketHandler) GetCompaniesBySector(c *gin.Context) {
	sector := c.Param("sector")
	data, err := h.priceEngine.GetCompaniesBySector(sector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "sector": sector, "count": len(data)})
}

// ──────────────────── New / Old Companies ────────────────────

// GetNewCompanies godoc
// @Summary Get recently listed companies
// @Tags Market
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.LiveTradingData
// @Router /market/companies/new [get]
func (h *MarketHandler) GetNewCompanies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.priceEngine.GetNewCompanies(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "count": len(data)})
}

// GetOldCompanies godoc
// @Summary Get oldest listed companies
// @Tags Market
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.LiveTradingData
// @Router /market/companies/old [get]
func (h *MarketHandler) GetOldCompanies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.priceEngine.GetOldCompanies(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "count": len(data)})
}

// ──────────────────── SSE Streaming ────────────────────

// StreamPrices godoc
// @Summary Stream live price updates (SSE)
// @Tags Market
// @Produce text/event-stream
// @Router /market/stream [get]
func (h *MarketHandler) StreamPrices(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	clientChan := h.eventHub.Subscribe()
	defer h.eventHub.Unsubscribe(clientChan)

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-clientChan:
			if !ok {
				return false
			}
			c.SSEvent("price_update", msg)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

// ──────────────────── Price Triggers ────────────────────

type CreateTriggerRequest struct {
	CompanyID    string          `json:"company_id" binding:"required"`
	TriggerPrice decimal.Decimal `json:"trigger_price" binding:"required"`
	SharesQty    int64           `json:"shares_qty" binding:"required"`
	Direction    string          `json:"direction" binding:"required"` // above, below
}

// CreateTrigger godoc
// @Summary Create a price trigger
// @Tags Triggers
// @Accept json
// @Produce json
// @Param body body CreateTriggerRequest true "Trigger request"
// @Security BearerAuth
// @Success 201 {object} models.PriceTrigger
// @Router /market/triggers [post]
func (h *MarketHandler) CreateTrigger(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trigger := &models.PriceTrigger{
		UserID:       userID.(string),
		CompanyID:    req.CompanyID,
		TriggerPrice: req.TriggerPrice,
		SharesQty:    req.SharesQty,
		Direction:    req.Direction,
		Status:       models.TriggerStatusActive,
	}

	if err := h.triggerWorker.CreateTrigger(trigger); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": trigger})
}

// CancelTrigger godoc
// @Summary Cancel a price trigger
// @Tags Triggers
// @Produce json
// @Param id path string true "Trigger ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /market/triggers/{id}/cancel [put]
func (h *MarketHandler) CancelTrigger(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	triggerID := c.Param("id")
	if err := h.triggerWorker.CancelTrigger(triggerID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trigger cancelled"})
}

// GetUserTriggers godoc
// @Summary Get user's price triggers
// @Tags Triggers
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PriceTrigger
// @Router /market/triggers [get]
func (h *MarketHandler) GetUserTriggers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	triggers, err := h.triggerWorker.GetUserTriggers(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": triggers, "count": len(triggers)})
}
